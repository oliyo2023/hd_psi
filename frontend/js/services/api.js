// API服务基类
class ApiService {
    constructor() {
        this.baseURL = 'http://localhost:8080';
        this.axios = axios.create({
            baseURL: this.baseURL,
            timeout: 10000,
            headers: {
                'Content-Type': 'application/json'
            }
        });

        // 请求拦截器 - 添加认证令牌
        this.axios.interceptors.request.use(
            config => {
                const token = localStorage.getItem('token');
                if (token) {
                    config.headers['Authorization'] = `Bearer ${token}`;
                }
                return config;
            },
            error => {
                return Promise.reject(error);
            }
        );

        // 响应拦截器 - 处理错误
        this.axios.interceptors.response.use(
            response => {
                return response;
            },
            error => {
                if (error.response) {
                    // 服务器返回错误状态码
                    if (error.response.status === 401) {
                        // 未授权，清除令牌并重定向到登录页
                        localStorage.removeItem('token');
                        localStorage.removeItem('user');
                        window.location.href = '/#/login';
                    }

                    // 显示错误消息
                    const errorMessage = error.response.data.error || '请求失败';
                    console.error(errorMessage);
                    // 错误消息将由调用者处理
                } else if (error.request) {
                    // 请求发送但未收到响应
                    console.error('服务器无响应，请检查网络连接');
                    // 错误消息将由调用者处理
                } else {
                    // 请求设置时出错
                    console.error('请求配置错误');
                    // 错误消息将由调用者处理
                }

                return Promise.reject(error);
            }
        );
    }

    // 通用GET请求
    async get(url, params = {}) {
        try {
            const response = await this.axios.get(url, { params });
            return response.data;
        } catch (error) {
            console.error('GET请求错误:', error);
            throw error;
        }
    }

    // 通用POST请求
    async post(url, data = {}) {
        try {
            const response = await this.axios.post(url, data);
            return response.data;
        } catch (error) {
            console.error('POST请求错误:', error);
            throw error;
        }
    }

    // 通用PUT请求
    async put(url, data = {}) {
        try {
            const response = await this.axios.put(url, data);
            return response.data;
        } catch (error) {
            console.error('PUT请求错误:', error);
            throw error;
        }
    }

    // 通用DELETE请求
    async delete(url) {
        try {
            const response = await this.axios.delete(url);
            return response.data;
        } catch (error) {
            console.error('DELETE请求错误:', error);
            throw error;
        }
    }
}
