import http from './http'

/**
 * 上传API
 */
export default {
  /**
   * 上传图片
   */
  uploadImage(file, onProgress) {
    const formData = new FormData()
    formData.append('file', file)

    return http.post('/admin/upload/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(percent)
        }
      }
    })
  },

  /**
   * 上传文章图片
   */
  uploadArticleImage(file, onProgress) {
    const formData = new FormData()
    formData.append('file', file)

    return http.post('/admin/upload/article-image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const percent = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(percent)
        }
      }
    })
  }
}

