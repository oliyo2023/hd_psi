import api from './api'

// 认证服务
export default {
  // 用户登录
  login(username, password) {
    return api.post('/auth/login', { username, password })
  },
  
  // 用户注册
  register(userData) {
    return api.post('/auth/register', userData)
  },
  
  // 获取当前用户信息
  getProfile() {
    return api.get('/api/profile')
  },
  
  // 更新用户信息
  updateProfile(userData) {
    return api.put('/api/profile', userData)
  },
  
  // 修改密码
  changePassword(oldPassword, newPassword) {
    return api.put('/api/change-password', {
      old_password: oldPassword,
      new_password: newPassword
    })
  },
  
  // 用户登出
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    window.location.href = '/login'
  },
  
  // 检查用户是否已认证
  isAuthenticated() {
    return !!localStorage.getItem('token')
  },
  
  // 获取当前用户
  getCurrentUser() {
    const userJson = localStorage.getItem('user')
    return userJson ? JSON.parse(userJson) : null
  },
  
  // 获取认证令牌
  getToken() {
    return localStorage.getItem('token')
  }
}
