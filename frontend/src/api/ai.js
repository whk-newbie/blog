import http from './http'

export default {
  // Provider management
  listProviders() {
    return http.get('/admin/ai/providers')
  },
  getProvider(id) {
    return http.get(`/admin/ai/providers/${id}`)
  },
  createProvider(data) {
    return http.post('/admin/ai/providers', data)
  },
  updateProvider(id, data) {
    return http.put(`/admin/ai/providers/${id}`, data)
  },
  deleteProvider(id) {
    return http.delete(`/admin/ai/providers/${id}`)
  },
  checkProvider(id) {
    return http.post(`/admin/ai/providers/${id}/check`)
  },

  // Translation
  translateArticle(articleId, data) {
    return http.post(`/admin/ai/translate/${articleId}`, data)
  },

  // History
  getHistory(providerId) {
    return http.get(`/admin/ai/providers/${providerId}/history`)
  },
  clearHistory(providerId) {
    return http.post(`/admin/ai/providers/${providerId}/history/clear`)
  },
  saveMessage(data) {
    return http.post('/admin/ai/message', data)
  },

  // Chat (SSE streaming - uses fetch directly for stream support)
  chat(providerId, messages, signal) {
    const token = localStorage.getItem('token')
    const sessionId = window.__sessionId || ''
    return fetch('/api/v1/admin/ai/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        'X-Session-Id': sessionId,
      },
      body: JSON.stringify({ provider_id: providerId, messages }),
      signal,
    })
  },
}
