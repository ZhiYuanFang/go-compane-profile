import { uploadDualImage } from '@/api/client'
import { processDualImage, revokeObjectUrl } from '@/utils/imageProcess'
import { SLOT_STATUS } from '@/utils/dualSlot'

const DEFAULT_CONCURRENCY = 3

/**
 * Shared concurrency-limited queue keyed by slotId (used for compress and upload).
 */
export function createUploadQueue(concurrency = DEFAULT_CONCURRENCY) {
  const limit = Math.max(1, concurrency)
  /** @type {{ slotId: string, task: (signal: AbortSignal) => Promise<void> }[]} */
  const queue = []
  /** @type {Map<string, AbortController>} */
  const controllers = new Map()
  const cancelled = new Set()
  let active = 0

  function pump() {
    while (active < limit && queue.length) {
      const job = queue.shift()
      if (!job) break
      if (cancelled.has(job.slotId)) continue

      active += 1
      const ac = new AbortController()
      controllers.set(job.slotId, ac)

      Promise.resolve()
        .then(() => job.task(ac.signal))
        .catch(() => {})
        .finally(() => {
          controllers.delete(job.slotId)
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
  if (!sharedUploadQueue) sharedUploadQueue = createUploadQueue(DEFAULT_CONCURRENCY)
  return sharedUploadQueue
}

export function getSharedCompressQueue() {
  if (!sharedCompressQueue) sharedCompressQueue = createUploadQueue(DEFAULT_CONCURRENCY)
  return sharedCompressQueue
}

/** Cancel compress + upload work for a slot. */
export function cancelSlotUpload(slotId) {
  if (!slotId) return
  getSharedCompressQueue().cancel(slotId)
  getSharedUploadQueue().cancel(slotId)
}

/**
 * HTTP-only upload of existing pending blobs (upload pool concurrency 3).
 */
export function enqueuePendingUpload({
  slotId,
  pending,
  category,
  patch,
  queue = getSharedUploadQueue(),
}) {
  if (!slotId || !pending?.original || !pending?.thumb) return

  patch({ status: SLOT_STATUS.QUEUED, error: '' })

  queue.enqueue(slotId, async (signal) => {
    try {
      if (queue.isCancelled(slotId) || signal.aborted) return
      patch({ status: SLOT_STATUS.UPLOADING, error: '' })
      const result = await uploadDualImage({
        original: pending.original,
        thumb: pending.thumb,
        category,
        signal,
      })
      if (queue.isCancelled(slotId) || signal.aborted) return

      const donePatch = {
        thumb: result.thumbUrl || result.thumb || '',
        original: result.originalUrl || result.original || '',
        pending: null,
        status: SLOT_STATUS.DONE,
        error: '',
        localPreview: '',
      }
      patch(donePatch)
    } catch (err) {
      if (queue.isCancelled(slotId) || signal.aborted || err?.name === 'AbortError') return
      patch({
        status: SLOT_STATUS.ERROR,
        error: err?.message || '上传失败',
      })
    }
  })
}

/**
 * Compress in compress pool (concurrency 3), then enqueue HTTP upload.
 * Sets instant localPreview from the selected File.
 */
export function enqueueFileUpload({
  slotId,
  file,
  category,
  patch,
  prevLocalPreview = '',
  compressQueue = getSharedCompressQueue(),
}) {
  if (!slotId || !file) return

  // Replace any in-flight work for this slot
  cancelSlotUpload(slotId)

  if (prevLocalPreview) revokeObjectUrl(prevLocalPreview)
  const localPreview = URL.createObjectURL(file)

  patch({
    localPreview,
    status: SLOT_STATUS.QUEUED,
    error: '',
    pending: null,
  })

  compressQueue.enqueue(slotId, async (signal) => {
    try {
      if (compressQueue.isCancelled(slotId) || signal.aborted) return
      patch({ status: SLOT_STATUS.COMPRESSING, error: '' })
      const pair = await processDualImage(file)
      if (compressQueue.isCancelled(slotId) || signal.aborted) return

      patch({ pending: pair, error: '' })
      enqueuePendingUpload({
        slotId,
        pending: pair,
        category,
        patch,
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
 * Shared entry for gallery/cover: instant preview + compress pool + upload pool.
 */
export function startSlotPipeline(opts) {
  return enqueueFileUpload(opts)
}
