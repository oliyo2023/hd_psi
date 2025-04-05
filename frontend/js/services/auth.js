// 认证服务
class AuthService extends ApiService {
    constructor() {
        super();
    }

    // 用户登录
    async login(username, password) {
        try {
            const response = await this.post('/auth/login', { username, password });

            // 保存令牌和用户信息到本地存储
            localStorage.setItem('token', response.token);
            localStorage.setItem('user', JSON.stringify(response.user));

            return response;
        } catch (error) {
            throw error;
        }
    }

    // 用户注册
    async register(userData) {
        try {
            const response = await this.post('/auth/register', userData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取当前用户信息
    async getProfile() {
        try {
            const response = await this.get('/api/profile');
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 更新用户信息
    async updateProfile(userData) {
        try {
            const response = await this.put('/api/profile', userData);

            // 更新本地存储中的用户信息
            const currentUser = JSON.parse(localStorage.getItem('user') || '{}');
            const updatedUser = { ...currentUser, ...userData };
            localStorage.setItem('user', JSON.stringify(updatedUser));

            return response;
        } catch (error) {
            throw error;
        }
    }

    // 修改密码
    async changePassword(oldPassword, newPassword) {
        try {
            const response = await this.put('/api/change-password', {
                old_password: oldPassword,
                new_password: newPassword
            });
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 用户登出
    logout() {
        // 清除本地存储中的认证信息
        localStorage.removeItem('token');
        localStorage.removeItem('user');

        // 重定向到登录页
        window.location.href = '/#/login';
    }

    // 检查用户是否已认证
    isAuthenticated() {
        return !!localStorage.getItem('token');
    }

    // 获取当前用户
    getCurrentUser() {
        const userJson = localStorage.getItem('user');
        return userJson ? JSON.parse(userJson) : null;
    }

    // 获取认证令牌
    getToken() {
        return localStorage.getItem('token');
    }
}
