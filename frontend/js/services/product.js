// 商品服务
class ProductService extends ApiService {
    constructor() {
        super();
    }

    // 获取商品列表
    async getProducts(params = {}) {
        try {
            const response = await this.get('/api/products', params);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取单个商品
    async getProduct(id) {
        try {
            const response = await this.get(`/api/products/${id}`);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 创建商品
    async createProduct(productData) {
        try {
            const response = await this.post('/api/products', productData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 更新商品
    async updateProduct(id, productData) {
        try {
            const response = await this.put(`/api/products/${id}`, productData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 删除商品
    async deleteProduct(id) {
        try {
            const response = await this.delete(`/api/products/${id}`);
            return response;
        } catch (error) {
            throw error;
        }
    }
}
