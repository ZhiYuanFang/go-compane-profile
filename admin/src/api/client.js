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

export function listPortfolios() {
  return request('/admin/api/portfolios')
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
    body: { ids },
  })
}

export function getPricing() {
  return request('/admin/api/pricing')
}

export function putPricing(payload) {
  return request('/admin/api/pricing', { method: 'PUT', body: payload })
}

export async function uploadDualImage({ original, thumb, category }) {
  const form = new FormData()
  form.append('original', original, original.name || 'original.jpg')
  form.append('thumb', thumb, thumb.name || 'thumb.jpg')
  form.append('category', category || 'general')
  return request('/admin/api/upload', { method: 'POST', body: form })
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
