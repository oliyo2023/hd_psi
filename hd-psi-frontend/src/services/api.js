import axios from 'axios'
import { useMessage } from 'naive-ui'

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
    const message = useMessage()
    
    if (error.response) {
      // 服务器返回错误状态码
      const status = error.response.status
      const errorMsg = error.response.data?.error || '请求失败'
      
      if (status === 401) {
        // 未授权，清除令牌并重定向到登录页
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        window.location.href = '/login'
        message.error('登录已过期，请重新登录')
      } else {
        message.error(errorMsg)
      }
    } else if (error.request) {
      // 请求发送但未收到响应
      message.error('服务器无响应，请检查网络连接')
    } else {
      // 请求设置时出错
      message.error('请求配置错误')
    }
    
    return Promise.reject(error)
  }
)

export default api
