import { uploadDualImage } from '@/api/client'
import { processDualImage, revokeObjectUrl } from '@/utils/imageProcess'
import { createPreviewObjectUrl } from '@/utils/previewImage'
import { SLOT_STATUS } from '@/utils/dualSlot'

/** Desktop admin may compress/watermark several images at once. */
export const COMPRESS_CONCURRENCY = 2
/** Uploads stay serial: concurrency 3 previously caused gateway 504s. */
export const UPLOAD_CONCURRENCY = 1

/**
 * Wait for the browser to finish a paint before starting heavy canvas work, so
 * newly added local previews are visible first. The timer is a fallback for
 * background tabs where requestAnimationFrame is throttled.
 */
function yieldToPaint() {
  return new Promise((resolve) => {
    let settled = false
    const done = () => {
      if (settled) return
      settled = true
      resolve()
    }
    setTimeout(done, 100)
    if (typeof requestAnimationFrame === 'function') {
      requestAnimationFrame(() => setTimeout(done, 0))
    }
  })
}

/**
 * Shared concurrency-limited queue keyed by slotId (used for compress and upload).
 */
export function createUploadQueue(concurrency = UPLOAD_CONCURRENCY) {
  const limit = Math.max(1, concurrency)
  /** @type {{ slotId: string, task: (signal: AbortSignal) => Promise<void> }[]} */
  const queue = []
  /** @type {Map<string, AbortController>} */
  const controllers = new Map()
  const cancelled = new Set()
  /** Slots with a job currently running; avoids overlapping old/new work for one slot. */
  const activeSlots = new Set()
  let active = 0

  function pump() {
    while (active < limit && queue.length) {
      const index = queue.findIndex((job) => !activeSlots.has(job.slotId))
      if (index < 0) break
      const [job] = queue.splice(index, 1)
      if (!job || cancelled.has(job.slotId)) continue

      active += 1
      activeSlots.add(job.slotId)
      const ac = new AbortController()
      controllers.set(job.slotId, ac)

      Promise.resolve()
        .then(() => job.task(ac.signal))
        .catch(() => {})
        .finally(() => {
          if (controllers.get(job.slotId) === ac) controllers.delete(job.slotId)
          activeSlots.delete(job.slotId)
          active -= 1
          pump()
        })
    }
  }

  return {
    enqueue(slotId, task) {
      cancelled.delete(slotId)
      for (let i = queue.length - 1; i >= 0; i--) {
        if (queue[i].slotId === slotId) queue.splice(i, 1)
      }
      const prev = controllers.get(slotId)
      if (prev) prev.abort()
      queue.push({ slotId, task })
      pump()
    },
    cancel(slotId) {
      cancelled.add(slotId)
      for (let i = queue.length - 1; i >= 0; i--) {
        if (queue[i].slotId === slotId) queue.splice(i, 1)
      }
      controllers.get(slotId)?.abort()
    },
    isCancelled(slotId) {
      return cancelled.has(slotId)
    },
  }
}

let sharedUploadQueue = null
let sharedCompressQueue = null

export function getSharedUploadQueue() {
  if (!sharedUploadQueue) sharedUploadQueue = createUploadQueue(UPLOAD_CONCURRENCY)
  return sharedUploadQueue
}

export function getSharedCompressQueue() {
  if (!sharedCompressQueue) sharedCompressQueue = createUploadQueue(COMPRESS_CONCURRENCY)
  return sharedCompressQueue
}

let previewSeq = 0
/** @type {Map<string, number>} */
const previewTokenBySlot = new Map()

function beginPreview(slotId) {
  const token = ++previewSeq
  previewTokenBySlot.set(slotId, token)
  return token
}

function isPreviewCurrent(slotId, token) {
  return previewTokenBySlot.get(slotId) === token
}

function invalidatePreview(slotId) {
  previewTokenBySlot.delete(slotId)
}

/** Cancel compress + upload work for a slot. */
export function cancelSlotUpload(slotId) {
  if (!slotId) return
  invalidatePreview(slotId)
  getSharedCompressQueue().cancel(slotId)
  getSharedUploadQueue().cancel(slotId)
}

/**
 * HTTP-only upload of existing pending blobs.
 * @returns {Promise<object|void>} resolves with done patch, rejects on upload error
 */
export function enqueuePendingUpload({
  slotId,
  pending,
  category,
  patch,
  queue = getSharedUploadQueue(),
}) {
  if (!slotId || !pending?.original || !pending?.thumb) {
    return Promise.resolve()
  }

  patch({ status: SLOT_STATUS.QUEUED, error: '', uploadProgress: null })

  return new Promise((resolve, reject) => {
    queue.enqueue(slotId, async (signal) => {
      try {
        if (queue.isCancelled(slotId) || signal.aborted) {
          resolve()
          return
        }
        patch({ status: SLOT_STATUS.UPLOADING, error: '', uploadProgress: 0 })
        const result = await uploadDualImage({
          original: pending.original,
          thumb: pending.thumb,
          category,
          signal,
          onProgress: (pct) => {
            if (queue.isCancelled(slotId) || signal.aborted) return
            patch({ uploadProgress: pct })
          },
        })
        if (queue.isCancelled(slotId) || signal.aborted) {
          resolve()
          return
        }

        const donePatch = {
          thumb: result.thumbUrl || result.thumb || '',
          original: result.originalUrl || result.original || '',
          pending: null,
          status: SLOT_STATUS.DONE,
          error: '',
          localPreview: '',
          uploadProgress: null,
        }
        patch(donePatch)
        resolve(donePatch)
      } catch (err) {
        if (queue.isCancelled(slotId) || signal.aborted || err?.name === 'AbortError') {
          resolve()
          return
        }
        const message = err?.message || '上传失败'
        patch({
          status: SLOT_STATUS.ERROR,
          error: message,
          uploadProgress: null,
        })
        reject(err instanceof Error ? err : new Error(message))
      }
    })
  })
}

/**
 * Compress in compress pool and leave slot READY with pending blobs (no OSS upload).
 * Sets instant localPreview from the selected File unless skipLocalPreviewSetup.
 */
export function enqueueFileUpload({
  slotId,
  file,
  category = 'general',
  patch,
  prevLocalPreview = '',
  skipLocalPreviewSetup = false,
  compressQueue = getSharedCompressQueue(),
}) {
  if (!slotId || !file) return

  // Replace any in-flight work for this slot
  cancelSlotUpload(slotId)

  if (!skipLocalPreviewSetup) {
    if (prevLocalPreview) revokeObjectUrl(prevLocalPreview)
    patch({
      localPreview: '',
      status: SLOT_STATUS.QUEUED,
      error: '',
      pending: null,
      uploadProgress: null,
    })
    // Build a small preview off the full-size photo to keep memory bounded.
    const token = beginPreview(slotId)
    createPreviewObjectUrl(file)
      .then((url) => {
        if (!isPreviewCurrent(slotId, token) || compressQueue.isCancelled(slotId)) {
          revokeObjectUrl(url)
          return
        }
        patch({ localPreview: url })
      })
      .catch(() => {
        /* preview is optional; pending thumb will show when compression finishes */
      })
  } else {
    patch({
      status: SLOT_STATUS.QUEUED,
      error: '',
      pending: null,
      uploadProgress: null,
    })
  }

  compressQueue.enqueue(slotId, async (signal) => {
    try {
      if (compressQueue.isCancelled(slotId) || signal.aborted) return
      await yieldToPaint()
      if (compressQueue.isCancelled(slotId) || signal.aborted) return
      patch({ status: SLOT_STATUS.COMPRESSING, error: '' })
      const pair = await processDualImage(file, {
        watermark: category !== 'logo',
        onWatermarkStart: () => {
          if (compressQueue.isCancelled(slotId) || signal.aborted) return
          patch({ status: SLOT_STATUS.WATERMARKING })
        },
      })
      if (compressQueue.isCancelled(slotId) || signal.aborted) return

      patch({
        pending: pair,
        error: '',
        status: SLOT_STATUS.READY,
        uploadProgress: null,
      })
    } catch (err) {
      if (compressQueue.isCancelled(slotId) || signal.aborted || err?.name === 'AbortError') return
      patch({
        status: SLOT_STATUS.ERROR,
        error: err?.message || '图片处理失败',
      })
    }
  })
}

/**
 * Shared entry for gallery/cover: instant preview + compress only (upload deferred to save).
 */
export function startSlotPipeline(opts) {
  return enqueueFileUpload(opts)
}

/**
 * Upload all duals that still have pending blobs. Resolves when all succeed; rejects on first failure.
 * Successful slots remain DONE so a later retry only uploads remaining failures/pendings.
 *
 * @param {{ dual: object, category: string, patch: (partial: object) => void }[]} entries
 */
export async function flushPendingDualUploads(entries) {
  const list = (entries || []).filter(
    (e) => e?.dual?.slotId && e?.dual?.pending?.original && e?.dual?.pending?.thumb,
  )
  if (!list.length) return

  await Promise.all(
    list.map(({ dual, category, patch }) => {
      const prevLocalPreview = dual.localPreview
      return enqueuePendingUpload({
        slotId: dual.slotId,
        pending: dual.pending,
        category: category || 'general',
        patch,
      }).then((donePatch) => {
        // Parent-owned patches bypass component-level preview revocation on DONE.
        if (donePatch) revokeObjectUrl(prevLocalPreview)
      })
    }),
  )
}
