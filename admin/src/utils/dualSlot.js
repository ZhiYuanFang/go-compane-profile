export const SLOT_STATUS = {
  IDLE: 'idle',
  DOWNLOADING: 'downloading',
  COMPRESSING: 'compressing',
  WATERMARKING: 'watermarking',
  QUEUED: 'queued',
  READY: 'ready',
  UPLOADING: 'uploading',
  DONE: 'done',
  ERROR: 'error',
}

export function newSlotId() {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `slot-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`
}

export function createEmptyDual(src) {
  const thumb = src?.thumb || ''
  const original = src?.original || ''
  const hasRemote = !!(thumb || original)
  const pending = src?.pending || null
  let status = src?.status
  if (!status) {
    if (pending?.original && pending?.thumb) status = SLOT_STATUS.READY
    else if (hasRemote) status = SLOT_STATUS.DONE
    else status = SLOT_STATUS.IDLE
  }
  return {
    slotId: src?.slotId || newSlotId(),
    thumb,
    original,
    pending,
    localPreview: src?.localPreview || '',
    uploadProgress: typeof src?.uploadProgress === 'number' ? src.uploadProgress : null,
    status,
    error: src?.error || '',
  }
}

export function isEmptyDual(d) {
  if (!d) return true
  if (d.pending || d.localPreview) return false
  if (
    d.status === SLOT_STATUS.DOWNLOADING ||
    d.status === SLOT_STATUS.COMPRESSING ||
    d.status === SLOT_STATUS.WATERMARKING ||
    d.status === SLOT_STATUS.QUEUED ||
    d.status === SLOT_STATUS.READY ||
    d.status === SLOT_STATUS.UPLOADING
  ) {
    return false
  }
  return !d.thumb && !d.original
}

export function dualStatusLabel(d) {
  if (!d) return ''
  switch (d.status) {
    case SLOT_STATUS.DOWNLOADING:
      return '下载原图…'
    case SLOT_STATUS.COMPRESSING:
      return '处理中…'
    case SLOT_STATUS.WATERMARKING:
      return '添加水印中…'
    case SLOT_STATUS.QUEUED:
      return '排队中'
    case SLOT_STATUS.READY:
      return '待保存'
    case SLOT_STATUS.UPLOADING:
      if (typeof d.uploadProgress === 'number' && d.uploadProgress >= 0) {
        return `上传中 ${d.uploadProgress}%`
      }
      return '上传中'
    case SLOT_STATUS.DONE:
      return '已上传'
    case SLOT_STATUS.ERROR:
      return d.error ? `失败：${d.error}` : '失败'
    default:
      if (d.pending) return '待保存'
      if (d.thumb || d.original) return '已上传'
      return ''
  }
}

/** Four-tier tone: wait | busy | done | error | '' */
export function dualStatusTone(d) {
  if (!d) return ''
  switch (d.status) {
    case SLOT_STATUS.QUEUED:
    case SLOT_STATUS.READY:
      return 'wait'
    case SLOT_STATUS.DOWNLOADING:
    case SLOT_STATUS.COMPRESSING:
    case SLOT_STATUS.WATERMARKING:
    case SLOT_STATUS.UPLOADING:
      return 'busy'
    case SLOT_STATUS.DONE:
      return 'done'
    case SLOT_STATUS.ERROR:
      return 'error'
    default:
      if (d.pending) return 'wait'
      if (d.thumb || d.original) return 'done'
      return ''
  }
}

export function dualIsBlockingSave(d) {
  if (!d || isEmptyDual(d)) return false
  if (d.status === SLOT_STATUS.ERROR) return true
  if (
    d.status === SLOT_STATUS.DOWNLOADING ||
    d.status === SLOT_STATUS.COMPRESSING ||
    d.status === SLOT_STATUS.WATERMARKING ||
    d.status === SLOT_STATUS.QUEUED ||
    d.status === SLOT_STATUS.UPLOADING
  ) {
    return true
  }
  // READY with pending is allowed — save will flush uploads
  return false
}

export function dualsBlockingReason(duals) {
  const list = (duals || []).filter(Boolean)
  if (list.some((d) => d.status === SLOT_STATUS.ERROR)) {
    return '有图片上传失败，请重试或删除后再保存'
  }
  if (list.some((d) => dualIsBlockingSave(d))) {
    return '仍有图片处理中，请稍候再保存'
  }
  return ''
}
