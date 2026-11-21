/**
 * Fetch an already-uploaded CDN image through the admin media proxy and wrap it
 * as a File so it can go through the normal compress/watermark pipeline.
 * @param {string} url CDN URL (or OSS object key) of the existing image
 * @returns {Promise<File>}
 */
export async function fetchRemoteImageFile(url) {
  if (!url) throw new Error('没有可处理的线上图片')

  const res = await fetch(`/admin/api/media/proxy?url=${encodeURIComponent(url)}`, {
    credentials: 'include',
  })
  if (!res.ok) {
    let message = `获取线上图片失败（${res.status}）`
    try {
      const text = (await res.text()).trim()
      if (text) message = text
    } catch {
      /* ignore */
    }
    throw new Error(message)
  }

  // GoFrame auth errors still use HTTP 200 with a JSON body.
  if ((res.headers.get('content-type') || '').includes('application/json')) {
    const data = await res.json().catch(() => null)
    throw new Error(data?.message || '未登录或没有权限')
  }

  const blob = await res.blob()
  if (!blob.type.startsWith('image/')) {
    throw new Error('线上文件不是图片')
  }

  const ext = blob.type.includes('png') ? 'png' : blob.type.includes('webp') ? 'webp' : 'jpg'
  return new File([blob], `remote.${ext}`, { type: blob.type })
}