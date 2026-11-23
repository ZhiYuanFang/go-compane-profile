/**
 * Wrap an image Blob as a named File for the compress/watermark pipeline.
 * @param {Blob} blob
 * @param {string} [basename='image']
 * @returns {File}
 */
export function blobAsImageFile(blob, basename = 'image') {
  if (!blob) throw new Error('缺少图片数据')
  const type = blob.type || 'image/jpeg'
  if (type && !type.startsWith('image/')) {
    throw new Error('文件不是图片')
  }
  const ext = type.includes('png') ? 'png' : type.includes('webp') ? 'webp' : 'jpg'
  return new File([blob], `${basename}.${ext}`, {
    type,
    lastModified: Date.now(),
  })
}

/**
 * Fetch an already-uploaded CDN image directly (CORS) and wrap it as a File
 * so it can go through the normal compress/watermark pipeline.
 * @param {string} url Absolute CDN URL of the existing image
 * @returns {Promise<File>}
 */
export async function fetchRemoteImageFile(url) {
  const trimmed = typeof url === 'string' ? url.trim() : ''
  if (!trimmed || !/^https?:\/\//i.test(trimmed)) {
    throw new Error('没有可处理的线上图片')
  }

  let res
  try {
    res = await fetch(trimmed)
  } catch {
    throw new Error('获取线上图片失败（网络或跨域被拒绝）')
  }
  if (!res.ok) {
    throw new Error(`获取线上图片失败（${res.status}）`)
  }

  const blob = await res.blob()
  if (!blob.type.startsWith('image/')) {
    throw new Error('线上文件不是图片')
  }

  return blobAsImageFile(blob, 'remote')
}
