export const SLOT_STATUS = {
  IDLE: 'idle',
  COMPRESSING: 'compressing',
  QUEUED: 'queued',
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
  return {
    slotId: src?.slotId || newSlotId(),
    thumb,
    original,
    pending: src?.pending || null,
    status: src?.status || (hasRemote ? SLOT_STATUS.DONE : SLOT_STATUS.IDLE),
    error: src?.error || '',
  }
}

export function isEmptyDual(d) {
  if (!d) return true
  if (d.pending) return false
  if (d.status === SLOT_STATUS.COMPRESSING || d.status === SLOT_STATUS.QUEUED || d.status === SLOT_STATUS.UPLOADING) {
    return false
  }
  return !d.thumb && !d.original
}

export function dualStatusLabel(d) {
  if (!d) return ''
  switch (d.status) {
    case SLOT_STATUS.COMPRESSING:
      return '处理中…'
    case SLOT_STATUS.QUEUED:
      return '排队中'
    case SLOT_STATUS.UPLOADING:
      return '上传中'
    case SLOT_STATUS.DONE:
      return '已上传'
    case SLOT_STATUS.ERROR:
      return d.error ? `失败：${d.error}` : '失败'
    default:
      if (d.pending) return '排队中'
      if (d.thumb || d.original) return '已上传'
      return ''
  }
}

export function dualIsBlockingSave(d) {
  if (!d || isEmptyDual(d)) return false
  if (d.status === SLOT_STATUS.ERROR) return true
  if (
    d.status === SLOT_STATUS.COMPRESSING ||
    d.status === SLOT_STATUS.QUEUED ||
    d.status === SLOT_STATUS.UPLOADING
  ) {
    return true
  }
  return !!(d.pending && d.status !== SLOT_STATUS.DONE)
}

export function dualsBlockingReason(duals) {
  const list = (duals || []).filter(Boolean)
  if (list.some((d) => d.status === SLOT_STATUS.ERROR)) {
    return '有图片上传失败，请重试或删除后再保存'
  }
  if (list.some((d) => dualIsBlockingSave(d))) {
    return '仍有图片上传中，请稍候再保存'
  }
  return ''
}
