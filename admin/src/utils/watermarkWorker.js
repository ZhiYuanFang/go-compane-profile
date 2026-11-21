/**
 * Draws the tiled watermark and re-encodes the image off the main thread.
 * Falls back ({ fallback: true }) when the browser lacks worker canvas support.
 */
import { JPEG_QUALITY_STEPS, drawTiledWatermark } from './watermarkCanvas'

async function encode(canvas, type, maxBytes) {
  if (type === 'image/png') {
    return canvas.convertToBlob({ type: 'image/png' })
  }
  let out = null
  for (const quality of JPEG_QUALITY_STEPS) {
    out = await canvas.convertToBlob({ type: 'image/jpeg', quality })
    if (out.size <= maxBytes) break
  }
  return out
}

function canUseWorkerCanvas() {
  if (typeof OffscreenCanvas !== 'function') return false
  if (typeof createImageBitmap !== 'function') return false
  try {
    const probe = new OffscreenCanvas(1, 1)
    const ctx = probe.getContext('2d')
    return !!ctx && typeof probe.convertToBlob === 'function'
  } catch {
    return false
  }
}

async function watermark(blob, type, maxBytes) {
  const source = await createImageBitmap(blob, { imageOrientation: 'from-image' })
  const width = source.width
  const height = source.height
  if (!width || !height) {
    source.close?.()
    throw new Error('无法读取图片尺寸')
  }

  let canvas = new OffscreenCanvas(width, height)
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

  let out = await encode(canvas, type, maxBytes)
  while (out.size > maxBytes) {
    const scale = Math.max(0.5, Math.sqrt(maxBytes / out.size) * 0.95)
    if (scale >= 0.999) break
    const nextWidth = Math.max(1, Math.round(canvas.width * scale))
    const nextHeight = Math.max(1, Math.round(canvas.height * scale))
    const next = new OffscreenCanvas(nextWidth, nextHeight)
    const nextCtx = next.getContext('2d')
    if (!nextCtx) break
    nextCtx.drawImage(canvas, 0, 0, nextWidth, nextHeight)
    canvas = next
    out = await encode(canvas, type, maxBytes)
  }

  if (out.size > maxBytes) {
    throw new Error(`水印后图片无法压缩到 ${(maxBytes / 1024).toFixed(0)}KB 以内`)
  }
  return out
}

self.onmessage = async (event) => {
  const { id, blob, type, maxBytes } = event.data || {}
  if (!canUseWorkerCanvas()) {
    self.postMessage({ id, fallback: true })
    return
  }
  try {
    const out = await watermark(blob, type, maxBytes)
    self.postMessage({ id, blob: out })
  } catch (err) {
    self.postMessage({ id, error: err?.message || '图片编码失败' })
  }
}