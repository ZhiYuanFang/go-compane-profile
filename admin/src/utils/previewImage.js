/**
 * Build a small object-URL preview (long edge ≤ 480px) from a selected file.
 * Full-size previews of many photos exhaust memory and make the desktop tab jank.
 */

const PREVIEW_MAX_EDGE = 480
const PREVIEW_QUALITY = 0.72

async function decodeSource(file) {
  if (typeof createImageBitmap === 'function') {
    try {
      return await createImageBitmap(file, { imageOrientation: 'from-image' })
    } catch {
      /* fall through to <img> decoding */
    }
  }
  const url = URL.createObjectURL(file)
  try {
    return await new Promise((resolve, reject) => {
      const img = new Image()
      img.onload = () => resolve(img)
      img.onerror = () => reject(new Error('图片解码失败'))
      img.src = url
    })
  } finally {
    URL.revokeObjectURL(url)
  }
}

async function encodeCanvas(canvas, type) {
  if (typeof canvas.convertToBlob === 'function') {
    return canvas.convertToBlob({ type, quality: PREVIEW_QUALITY })
  }
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => (blob ? resolve(blob) : reject(new Error('预览生成失败'))),
      type,
      PREVIEW_QUALITY,
    )
  })
}

/**
 * @param {File} file
 * @returns {Promise<string>} object URL of the small preview
 */
export async function createPreviewObjectUrl(file) {
  const source = await decodeSource(file)
  const width = source.width || source.naturalWidth || 0
  const height = source.height || source.naturalHeight || 0
  if (!width || !height) {
    source.close?.()
    throw new Error('无法读取图片尺寸')
  }

  const scale = Math.min(1, PREVIEW_MAX_EDGE / Math.max(width, height))
  const targetWidth = Math.max(1, Math.round(width * scale))
  const targetHeight = Math.max(1, Math.round(height * scale))
  const type = file.type === 'image/png' ? 'image/png' : 'image/jpeg'

  try {
    if (typeof OffscreenCanvas === 'function') {
      const canvas = new OffscreenCanvas(targetWidth, targetHeight)
      const ctx = canvas.getContext('2d')
      if (!ctx) throw new Error('浏览器不支持图片画布')
      ctx.drawImage(source, 0, 0, targetWidth, targetHeight)
      return URL.createObjectURL(await encodeCanvas(canvas, type))
    }

    const canvas = document.createElement('canvas')
    canvas.width = targetWidth
    canvas.height = targetHeight
    const ctx = canvas.getContext('2d')
    if (!ctx) throw new Error('浏览器不支持图片画布')
    ctx.drawImage(source, 0, 0, targetWidth, targetHeight)
    return URL.createObjectURL(await encodeCanvas(canvas, type))
  } finally {
    source.close?.()
  }
}