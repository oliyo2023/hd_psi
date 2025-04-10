import axios from 'axios'

// 引入全局axios实例，用于重试请求

// 创建axios实例
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data
  },
  async error => {
    // 注意: 不要在拦截器中使用 useMessage，因为它需要在组件中使用
    console.error('API 请求错误:', error)

    const originalRequest = error.config

    if (error.response) {
      // 服务器返回错误状态码
      const status = error.response.status

      // 如果是401错误且不是刷新令牌的请求
      if (status === 401 && !originalRequest._retry) {
        // 标记该请求已尝试过刷新令牌
        originalRequest._retry = true

        try {
          // 尝试使用刷新令牌获取新的访问令牌
          const refreshToken = localStorage.getItem('refreshToken')
          if (refreshToken) {
            const rememberMe = localStorage.getItem('rememberMe') === 'true'

            // 引入auth服务
            const auth = await import('./auth').then(module => module.default)

            // 尝试刷新令牌
            const response = await auth.refreshToken(refreshToken, rememberMe)

            // 更新本地存储
            localStorage.setItem('token', response.token)
            localStorage.setItem('refreshToken', response.refresh_token)
            localStorage.setItem('tokenExpires', response.expires_at)
            localStorage.setItem('user', JSON.stringify(response.user))

            // 更新原始请求的认证信息
            originalRequest.headers['Authorization'] = `Bearer ${response.token}`

            // 重新发送原始请求
            return axios(originalRequest)
          }
        } catch (refreshError) {
          console.error('刷新令牌失败:', refreshError)

          // 刷新令牌失败，清除所有认证信息
          localStorage.removeItem('token')
          localStorage.removeItem('refreshToken')
          localStorage.removeItem('tokenExpires')
          localStorage.removeItem('user')

          // 重定向到登录页
          const currentPath = window.location.pathname
          if (currentPath !== '/login') {
            window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
          }
        }
      } else if (status === 401) {
        // 如果已经尝试过刷新令牌依然失败，清除信息并重定向
        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
        localStorage.removeItem('tokenExpires')
        localStorage.removeItem('user')

        // 重定向到登录页
        const currentPath = window.location.pathname
        if (currentPath !== '/login') {
          window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
        }
      }
    }

    return Promise.reject(error)
  }
)

export default api
