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
  return xhrFormUpload('/admin/api/upload', form, { signal, onProgress })
}

/**
 * Multipart POST via XHR so upload progress is available (fetch cannot).
 */
function xhrFormUpload(path, formData, { signal, onProgress } = {}) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open('POST', `${baseURL}${path}`)
    xhr.withCredentials = true

    const onAbort = () => {
      xhr.abort()
    }
    if (signal) {
      if (signal.aborted) {
        const err = new Error('Aborted')
        err.name = 'AbortError'
        reject(err)
        return
      }
      signal.addEventListener('abort', onAbort)
    }

    let lastPct = -1
    xhr.upload.onprogress = (e) => {
      if (!onProgress || !e.lengthComputable || e.total <= 0) return
      const pct = Math.min(100, Math.round((e.loaded / e.total) * 100))
      if (pct === lastPct) return
      lastPct = pct
      onProgress(pct)
    }

    xhr.onload = () => {
      if (signal) signal.removeEventListener('abort', onAbort)
      let data = null
      const text = xhr.responseText || ''
      if (text) {
        try {
          data = JSON.parse(text)
        } catch {
          data = text
        }
      }

      if (xhr.status < 200 || xhr.status >= 300) {
        const message =
          (data && (data.message || data.msg || data.error)) ||
          (typeof data === 'string' ? data : null) ||
          `Request failed (${xhr.status})`
        const err = new Error(message)
        err.status = xhr.status
        err.data = data
        reject(err)
        return
      }

      if (data && typeof data === 'object' && 'code' in data && 'data' in data) {
        if (data.code !== 0 && data.code !== 200) {
          const err = new Error(data.message || data.msg || 'API error')
          err.status = xhr.status
          err.data = data
          reject(err)
          return
        }
        resolve(data.data)
        return
      }
      resolve(data)
    }

    xhr.onerror = () => {
      if (signal) signal.removeEventListener('abort', onAbort)
      reject(new Error('网络错误'))
    }

    xhr.onabort = () => {
      if (signal) signal.removeEventListener('abort', onAbort)
      const err = new Error('Aborted')
      err.name = 'AbortError'
      reject(err)
    }

    xhr.send(formData)
  })
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
