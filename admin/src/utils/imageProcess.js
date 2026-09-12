import imageCompression from 'browser-image-compression'

const ORIGINAL_MAX_MB = 3
const THUMB_MAX_MB = 0.5
const ORIGINAL_MAX = ORIGINAL_MAX_MB * 1024 * 1024
const THUMB_MAX = THUMB_MAX_MB * 1024 * 1024

function extFromType(type) {
  if (!type) return 'jpg'
  if (type.includes('png')) return 'png'
  if (type.includes('webp')) return 'webp'
  return 'jpg'
}

function asNamedBlob(blob, prefix, sourceFile) {
  const ext = extFromType(blob.type || sourceFile?.type)
  return new File([blob], `${prefix}.${ext}`, {
    type: blob.type || sourceFile?.type || 'image/jpeg',
    lastModified: Date.now(),
  })
}

async function compressUnder(file, maxSizeMB, maxWidthOrHeight) {
  let quality = 0.92
  let width = maxWidthOrHeight
  let blob = file

  for (let i = 0; i < 8; i++) {
    blob = await imageCompression(file, {
      maxSizeMB,
      maxWidthOrHeight: width,
      useWebWorker: true,
      initialQuality: quality,
      fileType: file.type === 'image/png' ? 'image/png' : 'image/jpeg',
    })
    if (blob.size <= maxSizeMB * 1024 * 1024) return blob
    quality = Math.max(0.4, quality - 0.1)
    width = Math.round(width * 0.85)
  }

  if (blob.size > maxSizeMB * 1024 * 1024) {
    throw new Error(`无法将图片压缩到 ${(maxSizeMB * 1024).toFixed(0)}KB 以内`)
  }
  return blob
}

/**
 * Produce original (≤3MiB) + thumb (≤0.5MiB) blobs from a user-selected File.
 */
export async function processDualImage(file) {
  if (!file || !file.type?.startsWith('image/')) {
    throw new Error('请选择图片文件')
  }

  const originalBlob =
    file.size <= ORIGINAL_MAX
      ? file
      : await compressUnder(file, ORIGINAL_MAX_MB, 4096)

  const thumbBlob = await compressUnder(file, THUMB_MAX_MB, 1280)

  if (originalBlob.size > ORIGINAL_MAX) {
    throw new Error('原图超过 3MB，请换一张较小的图片')
  }
  if (thumbBlob.size > THUMB_MAX) {
    throw new Error('缩略图超过 0.5MB，请换一张或降低分辨率')
  }

  return {
    original: asNamedBlob(originalBlob, 'original', file),
    thumb: asNamedBlob(thumbBlob, 'thumb', file),
  }
}

export function revokeObjectUrl(url) {
  if (url && String(url).startsWith('blob:')) {
    URL.revokeObjectURL(url)
  }
}
