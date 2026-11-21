import imageCompression from 'browser-image-compression'

const ORIGINAL_MAX_MB = 3
const THUMB_MAX_MB = 0.5
const ORIGINAL_MAX = ORIGINAL_MAX_MB * 1024 * 1024
const THUMB_MAX = THUMB_MAX_MB * 1024 * 1024

export const WATERMARK_TEXT = 'MUHU'

const WATERMARK_ANGLE = (-30 * Math.PI) / 180
const WATERMARK_FONT_RATIO = 0.04
const WATERMARK_MIN_FONT = 15
const WATERMARK_STEP_X_RATIO = 0.4
const WATERMARK_STEP_Y_RATIO = 0.3
const WATERMARK_OPACITY = 0.18
const WATERMARK_SHADOW = 'rgba(0, 0, 0, 0.18)'
const JPEG_QUALITY_STEPS = [0.9, 0.8, 0.7, 0.6, 0.5]

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

function loadImageElement(blob) {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(blob)
    const img = new Image()
    img.onload = () => {
      URL.revokeObjectURL(url)
      resolve(img)
    }
    img.onerror = () => {
      URL.revokeObjectURL(url)
      reject(new Error('图片解码失败'))
    }
    img.src = url
  })
}

/** Decode a blob with EXIF orientation applied when the browser supports it. */
async function decodeForCanvas(blob) {
  if (typeof createImageBitmap === 'function') {
    try {
      return await createImageBitmap(blob, { imageOrientation: 'from-image' })
    } catch {
      /* fall through to <img> decoding */
    }
  }
  return loadImageElement(blob)
}

function drawTiledWatermark(ctx, width, height) {
  const shortEdge = Math.min(width, height)
  const fontSize = Math.max(WATERMARK_MIN_FONT, Math.round(shortEdge * WATERMARK_FONT_RATIO))
  const stepY = shortEdge * WATERMARK_STEP_Y_RATIO
  const radius = Math.sqrt(width * width + height * height) / 2

  ctx.save()
  ctx.font = `bold ${fontSize}px sans-serif`
  ctx.textAlign = 'center'
  ctx.textBaseline = 'middle'
  ctx.fillStyle = `rgba(255, 255, 255, ${WATERMARK_OPACITY})`
  ctx.shadowColor = WATERMARK_SHADOW
  ctx.shadowBlur = 2

  const stepX = ctx.measureText(WATERMARK_TEXT).width + shortEdge * WATERMARK_STEP_X_RATIO

  ctx.translate(width / 2, height / 2)
  ctx.rotate(WATERMARK_ANGLE)

  let row = 0
  for (let y = -radius; y <= radius; y += stepY, row += 1) {
    const offset = row % 2 === 0 ? 0 : stepX / 2
    for (let x = -radius - stepX + offset; x <= radius + stepX; x += stepX) {
      ctx.fillText(WATERMARK_TEXT, x, y)
    }
  }
  ctx.restore()
}

function canvasToBlob(canvas, type, quality) {
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (blob) resolve(blob)
        else reject(new Error('图片编码失败'))
      },
      type,
      quality,
    )
  })
}

async function encodeWatermarked(canvas, type, maxBytes) {
  if (type === 'image/png') {
    const png = await canvasToBlob(canvas, 'image/png')
    if (png.size <= maxBytes) return png
  } else {
    for (const quality of JPEG_QUALITY_STEPS) {
      const jpeg = await canvasToBlob(canvas, 'image/jpeg', quality)
      if (jpeg.size <= maxBytes) return jpeg
    }
  }

  // Last resort: reuse the existing compressor to downscale until under the cap.
  const fallbackType = type === 'image/png' ? 'image/png' : 'image/jpeg'
  const fallback = await canvasToBlob(canvas, fallbackType, 0.9)
  try {
    const wrapped = asNamedBlob(fallback, 'watermarked', fallback)
    return await compressUnder(wrapped, maxBytes / (1024 * 1024), Math.max(canvas.width, canvas.height))
  } catch {
    throw new Error(`水印后图片无法压缩到 ${(maxBytes / 1024).toFixed(0)}KB 以内`)
  }
}

/**
 * Draw the tiled MUHU watermark over a compressed blob and re-encode it.
 * Returns a blob within maxBytes or throws.
 */
export async function watermarkBlob(blob, maxBytes) {
  if (!blob) throw new Error('缺少待加水印的图片')

  const source = await decodeForCanvas(blob)
  const width = source.width || source.naturalWidth || 0
  const height = source.height || source.naturalHeight || 0
  if (!width || !height) {
    source.close?.()
    throw new Error('无法读取图片尺寸')
  }

  const canvas = document.createElement('canvas')
  canvas.width = width
  canvas.height = height
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    source.close?.()
    throw new Error('浏览器不支持图片画布')
  }

  try {
    ctx.drawImage(source, 0, 0, width, height)
    drawTiledWatermark(ctx, width, height)
  } finally {
    source.close?.()
  }

  const type = blob.type === 'image/png' ? 'image/png' : 'image/jpeg'
  return encodeWatermarked(canvas, type, maxBytes)
}

/**
 * Produce original (≤3MiB) + thumb (≤0.5MiB) blobs from a user-selected File.
 * When `watermark` is true, content images get the tiled MUHU mark before READY.
 */
export async function processDualImage(file, { watermark = false } = {}) {
  if (!file || !file.type?.startsWith('image/')) {
    throw new Error('请选择图片文件')
  }

  let originalBlob =
    file.size <= ORIGINAL_MAX
      ? file
      : await compressUnder(file, ORIGINAL_MAX_MB, 4096)

  let thumbBlob = await compressUnder(file, THUMB_MAX_MB, 1280)

  if (originalBlob.size > ORIGINAL_MAX) {
    throw new Error('原图超过 3MB，请换一张较小的图片')
  }
  if (thumbBlob.size > THUMB_MAX) {
    throw new Error('缩略图超过 0.5MB，请换一张或降低分辨率')
  }

  if (watermark) {
    originalBlob = await watermarkBlob(originalBlob, ORIGINAL_MAX)
    thumbBlob = await watermarkBlob(thumbBlob, THUMB_MAX)
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