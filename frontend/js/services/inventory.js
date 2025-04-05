// 库存服务
class InventoryService extends ApiService {
    constructor() {
        super();
    }

    // 获取库存列表
    async getInventories(params = {}) {
        try {
            const response = await this.get('/api/inventory', params);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取单个库存
    async getInventory(id) {
        try {
            const response = await this.get(`/api/inventory/${id}`);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 创建库存
    async createInventory(inventoryData) {
        try {
            const response = await this.post('/api/inventory', inventoryData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 更新库存
    async updateInventory(id, inventoryData) {
        try {
            const response = await this.put(`/api/inventory/${id}`, inventoryData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 删除库存
    async deleteInventory(id) {
        try {
            const response = await this.delete(`/api/inventory/${id}`);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取库存交易记录
    async getInventoryTransactions(params = {}) {
        try {
            const response = await this.get('/api/inventory-transactions', params);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 创建库存交易记录
    async createInventoryTransaction(transactionData) {
        try {
            const response = await this.post('/api/inventory-transactions', transactionData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取库存预警
    async getInventoryAlerts(params = {}) {
        try {
            const response = await this.get('/api/inventory-alerts', params);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 更新库存预警状态
    async updateAlertStatus(id, status) {
        try {
            const response = await this.put(`/api/inventory-alerts/${id}/status`, { status });
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 检查库存水平
    async checkInventoryLevels() {
        try {
            const response = await this.post('/api/inventory-alerts/check');
            return response;
        } catch (error) {
            throw error;
        }
    }
}
