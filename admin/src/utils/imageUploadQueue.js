import { uploadDualImage } from '@/api/client'
import { processDualImage } from '@/utils/imageProcess'
import { SLOT_STATUS } from '@/utils/dualSlot'

const DEFAULT_CONCURRENCY = 3

/**
 * Shared concurrency-limited upload queue keyed by slotId.
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

let sharedQueue = null

export function getSharedUploadQueue() {
  if (!sharedQueue) sharedQueue = createUploadQueue(DEFAULT_CONCURRENCY)
  return sharedQueue
}

/**
 * Compress + upload a file for a slot; patch(partial) updates the dual model.
 */
export function enqueueFileUpload({
  slotId,
  file,
  category,
  patch,
  queue = getSharedUploadQueue(),
}) {
  if (!slotId || !file) return

  patch({ status: SLOT_STATUS.QUEUED, error: '' })

  queue.enqueue(slotId, async (signal) => {
    try {
      if (queue.isCancelled(slotId) || signal.aborted) return
      patch({ status: SLOT_STATUS.COMPRESSING, error: '' })
      const pair = await processDualImage(file)
      if (queue.isCancelled(slotId) || signal.aborted) return

      patch({ pending: pair, status: SLOT_STATUS.UPLOADING, error: '' })
      const result = await uploadDualImage({
        original: pair.original,
        thumb: pair.thumb,
        category,
        signal,
      })
      if (queue.isCancelled(slotId) || signal.aborted) return

      patch({
        thumb: result.thumbUrl || result.thumb || '',
        original: result.originalUrl || result.original || '',
        pending: null,
        status: SLOT_STATUS.DONE,
        error: '',
      })
    } catch (err) {
      if (queue.isCancelled(slotId) || signal.aborted || err?.name === 'AbortError') return
      patch({
        status: SLOT_STATUS.ERROR,
        error: err?.message || '上传失败',
      })
    }
  })
}

export function cancelSlotUpload(slotId, queue = getSharedUploadQueue()) {
  if (!slotId) return
  queue.cancel(slotId)
}

/**
 * Re-upload from existing pending blobs (retry after failure).
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
      patch({
        thumb: result.thumbUrl || result.thumb || '',
        original: result.originalUrl || result.original || '',
        pending: null,
        status: SLOT_STATUS.DONE,
        error: '',
      })
    } catch (err) {
      if (queue.isCancelled(slotId) || signal.aborted || err?.name === 'AbortError') return
      patch({
        status: SLOT_STATUS.ERROR,
        error: err?.message || '上传失败',
      })
    }
  })
}
