/** Shared watermark drawing used by the main thread and the watermark worker. */

export const WATERMARK_TEXT = 'MUHU'

export const WATERMARK_ANGLE = (-30 * Math.PI) / 180
const WATERMARK_FONT_RATIO = 0.04
const WATERMARK_MIN_FONT = 15
const WATERMARK_STEP_X_RATIO = 0.4
const WATERMARK_STEP_Y_RATIO = 0.3
const WATERMARK_OPACITY = 0.18
const WATERMARK_SHADOW = 'rgba(0, 0, 0, 0.18)'
export const JPEG_QUALITY_STEPS = [0.9, 0.8, 0.7, 0.6, 0.5]

export function drawTiledWatermark(ctx, width, height) {
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