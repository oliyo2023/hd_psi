import axios from 'axios'

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
  error => {
    // 注意: 不要在拦截器中使用 useMessage，因为它需要在组件中使用
    console.error('API 请求错误:', error)

    if (error.response) {
      // 服务器返回错误状态码
      const status = error.response.status

      if (status === 401) {
        // 未授权，清除令牌并重定向到登录页
        localStorage.removeItem('token')
        localStorage.removeItem('user')

        // 使用当前路径作为重定向参数
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
