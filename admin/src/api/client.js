const baseURL = ''

async function request(path, options = {}) {
  const headers = { ...(options.headers || {}) }
  const isForm = typeof FormData !== 'undefined' && options.body instanceof FormData
  if (options.body && !isForm && !headers['Content-Type']) {
    headers['Content-Type'] = 'application/json'
  }

  const res = await fetch(`${baseURL}${path}`, {
    credentials: 'include',
    ...options,
    headers,
    body:
      options.body && !isForm && typeof options.body === 'object'
        ? JSON.stringify(options.body)
        : options.body,
  })

  const text = await res.text()
  let data = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = text
    }
  }

  if (!res.ok) {
    const message =
      (data && (data.message || data.msg || data.error)) ||
      (typeof data === 'string' ? data : null) ||
      `Request failed (${res.status})`
    const err = new Error(message)
    err.status = res.status
    err.data = data
    throw err
  }

  // GoFrame style: { code, message, data }
  if (data && typeof data === 'object' && 'code' in data && 'data' in data) {
    if (data.code !== 0 && data.code !== 200) {
      const err = new Error(data.message || data.msg || 'API error')
      err.status = res.status
      err.data = data
      throw err
    }
    return data.data
  }

  return data
}

export function login(password) {
  return request('/admin/api/login', {
    method: 'POST',
    body: { password },
  })
}

export function logout() {
  return request('/admin/api/logout', { method: 'POST' })
}

export function getCompany() {
  return request('/admin/api/company')
}

export function putCompany(payload) {
  return request('/admin/api/company', { method: 'PUT', body: payload })
}

export function listPortfolios(category) {
  const q = new URLSearchParams({ category })
  return request(`/admin/api/portfolios?${q}`)
}

export function getPortfolio(id) {
  return request(`/admin/api/portfolios/${id}`)
}

export function createPortfolio(payload) {
  return request('/admin/api/portfolios', { method: 'POST', body: payload })
}

export function updatePortfolio(id, payload) {
  return request(`/admin/api/portfolios/${id}`, { method: 'PUT', body: payload })
}

export function deletePortfolio(id) {
  return request(`/admin/api/portfolios/${id}`, { method: 'DELETE' })
}

export function putPortfolioImages(id, payload) {
  return request(`/admin/api/portfolios/${id}/gallery`, {
    method: 'PUT',
    body: payload,
  })
}

export function reorderPortfolios(payload) {
  const ids = (payload?.ids || []).map((id) => String(id))
  return request('/admin/api/portfolios/reorder', {
    method: 'POST',
    body: { category: payload.category, ids },
  })
}

export function listActivities() {
  return request('/admin/api/activities')
}

export function getActivity(id) {
  return request(`/admin/api/activities/${id}`)
}

export function createActivity(payload) {
  return request('/admin/api/activities', { method: 'POST', body: payload })
}

export function updateActivity(id, payload) {
  return request(`/admin/api/activities/${id}`, { method: 'PUT', body: payload })
}

export function deleteActivity(id) {
  return request(`/admin/api/activities/${id}`, { method: 'DELETE' })
}

export function reorderActivities(payload) {
  const ids = (payload?.ids || []).map((id) => String(id))
  return request('/admin/api/activities/reorder', {
    method: 'POST',
    body: { ids },
  })
}

export async function uploadDualImage({ original, thumb, category, signal, onProgress }) {
  const form = new FormData()
  form.append('original', original, original.name || 'original.jpg')
  form.append('thumb', thumb, thumb.name || 'thumb.jpg')
  form.append('category', category || 'general')
  return uploadDualImageStream(form, { signal, onProgress })
}

/**
 * POST multipart to /upload/stream and consume NDJSON progress/done/error events.
 * Progress reflects OSS-weighted put progress from the server.
 */
async function uploadDualImageStream(formData, { signal, onProgress } = {}) {
  const res = await fetch(`${baseURL}/admin/api/upload/stream`, {
    method: 'POST',
    body: formData,
    credentials: 'include',
    signal,
  })

  if (!res.ok) {
    let message = `Request failed (${res.status})`
    try {
      const text = await res.text()
      const line = text.split('\n').find((l) => l.trim())
      if (line) {
        const ev = JSON.parse(line)
        if (ev?.message) message = ev.message
      }
    } catch {
      /* ignore */
    }
    const err = new Error(message)
    err.status = res.status
    throw err
  }

  if (!res.body) {
    throw new Error('浏览器不支持流式上传响应')
  }

  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buf = ''
  let result = null
  let lastPct = -1

  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buf += decoder.decode(value, { stream: true })
    let nl
    while ((nl = buf.indexOf('\n')) >= 0) {
      const line = buf.slice(0, nl).trim()
      buf = buf.slice(nl + 1)
      if (!line) continue
      let ev
      try {
        ev = JSON.parse(line)
      } catch {
        continue
      }
      if (ev.type === 'progress') {
        const pct = typeof ev.pct === 'number' ? ev.pct : 0
        if (onProgress && pct !== lastPct) {
          lastPct = pct
          onProgress(pct)
        }
      } else if (ev.type === 'done') {
        result = {
          original: ev.original || '',
          thumb: ev.thumb || '',
          originalUrl: ev.original || '',
          thumbUrl: ev.thumb || '',
        }
      } else if (ev.type === 'error') {
        throw new Error(ev.message || '上传失败')
      }
    }
  }

  if (!result) {
    throw new Error('上传未返回结果')
  }
  return result
}

/** Resolve pending DualImageField value → { thumb, original } via upload if needed */
export async function resolveDualImage(value, category) {
  if (!value) return { thumb: '', original: '' }
  if (value.pending && value.pending.original && value.pending.thumb) {
    const result = await uploadDualImage({
      original: value.pending.original,
      thumb: value.pending.thumb,
      category,
    })
    return {
      thumb: result.thumbUrl || result.thumb || '',
      original: result.originalUrl || result.original || '',
    }
  }
  return {
    thumb: value.thumb || '',
    original: value.original || '',
  }
}
