import imageCompression from 'browser-image-compression'
import {
  JPEG_QUALITY_STEPS,
  WATERMARK_TEXT,
  drawTiledWatermark,
} from './watermarkCanvas'

export { WATERMARK_TEXT }

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

let watermarkWorker = null
let watermarkWorkerUnavailable = false
let watermarkRequestSeq = 0
/** @type {Map<number, { resolve: (blob: Blob | null) => void, reject: (err: Error) => void }>} */
const watermarkPending = new Map()

function disableWatermarkWorker(reason) {
  if (!watermarkWorkerUnavailable) {
    console.warn('[watermark] Worker 不可用，回退到主线程：', reason || 'unknown')
  }
  watermarkWorkerUnavailable = true
  try {
    watermarkWorker?.terminate()
  } catch {
    /* ignore */
  }
  watermarkWorker = null
  for (const entry of watermarkPending.values()) entry.resolve(null)
  watermarkPending.clear()
}

function getWatermarkWorker() {
  if (watermarkWorkerUnavailable) return null
  if (typeof Worker !== 'function') {
    watermarkWorkerUnavailable = true
    console.warn('[watermark] 当前浏览器不支持 Worker，回退到主线程')
    return null
  }
  if (!watermarkWorker) {
    try {
      watermarkWorker = new Worker(new URL('./watermarkWorker.js', import.meta.url), {
        type: 'module',
      })
      watermarkWorker.onmessage = (event) => {
        const data = event.data || {}
        const entry = watermarkPending.get(data.id)
        if (!entry) return
        watermarkPending.delete(data.id)
        if (data.fallback) {
          disableWatermarkWorker('worker canvas unsupported')
          entry.resolve(null)
        } else if (data.error) entry.reject(new Error(data.error))
        else entry.resolve(data.blob)
      }
      watermarkWorker.onerror = (event) => {
        // Worker unavailable (CSP, unsupported browser, bundling issue): use main thread.
        disableWatermarkWorker(event?.message || 'worker error')
      }
    } catch {
      watermarkWorkerUnavailable = true
      watermarkWorker = null
    }
  }
  return watermarkWorker
}

function watermarkInWorker(blob, maxBytes) {
  const worker = getWatermarkWorker()
  if (!worker) return Promise.resolve(null)
  const id = ++watermarkRequestSeq
  const type = blob.type === 'image/png' ? 'image/png' : 'image/jpeg'
  return new Promise((resolve, reject) => {
    watermarkPending.set(id, { resolve, reject })
    try {
      worker.postMessage({ id, blob, type, maxBytes })
    } catch (err) {
      watermarkPending.delete(id)
      disableWatermarkWorker(err)
      resolve(null)
    }
  })
}

/**
 * Draw the tiled MUHU watermark over a compressed blob and re-encode it.
 * Prefers a worker so large batches do not block the UI; falls back to canvas
 * on the main thread when worker canvas support is missing.
 * Returns a blob within maxBytes or throws.
 */
export async function watermarkBlob(blob, maxBytes) {
  if (!blob) throw new Error('缺少待加水印的图片')
  const fromWorker = await watermarkInWorker(blob, maxBytes)
  if (fromWorker) return fromWorker
  return watermarkBlobOnMainThread(blob, maxBytes)
}

async function watermarkBlobOnMainThread(blob, maxBytes) {
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
export async function processDualImage(file, { watermark = false, onWatermarkStart } = {}) {
  if (!file || !file.type?.startsWith('image/')) {
    throw new Error('请选择图片文件')
  }

  // Process one large blob at a time: parallel full-size decodes of several
  // images exhaust memory and stall the whole desktop tab.
  let originalBlob =
    file.size <= ORIGINAL_MAX
      ? file
      : await compressUnder(file, ORIGINAL_MAX_MB, 4096)

  // Derive the thumb from the already-compressed original so the source image
  // is not decoded again at full resolution.
  let thumbBlob = await compressUnder(originalBlob, THUMB_MAX_MB, 1280)

  if (originalBlob.size > ORIGINAL_MAX) {
    throw new Error('原图超过 3MB，请换一张较小的图片')
  }
  if (thumbBlob.size > THUMB_MAX) {
    throw new Error('缩略图超过 0.5MB，请换一张或降低分辨率')
  }

  if (watermark) {
    onWatermarkStart?.()
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