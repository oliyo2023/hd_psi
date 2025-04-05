// 会员服务
class MemberService extends ApiService {
    constructor() {
        super();
    }

    // 获取会员列表
    async getMembers(params = {}) {
        try {
            const response = await this.get('/api/members', params);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取单个会员
    async getMember(id) {
        try {
            const response = await this.get(`/api/members/${id}`);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 创建会员
    async createMember(memberData) {
        try {
            const response = await this.post('/api/members', memberData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 更新会员
    async updateMember(id, memberData) {
        try {
            const response = await this.put(`/api/members/${id}`, memberData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 删除会员
    async deleteMember(id) {
        try {
            const response = await this.delete(`/api/members/${id}`);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取会员积分
    async getMemberPoints(id) {
        try {
            const response = await this.get(`/api/members/${id}/points`);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 获取会员积分交易记录
    async getPointsTransactions(id) {
        try {
            const response = await this.get(`/api/members/${id}/points/transactions`);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 增加会员积分
    async addPoints(id, pointsData) {
        try {
            const response = await this.post(`/api/members/${id}/points/add`, pointsData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 扣减会员积分
    async deductPoints(id, pointsData) {
        try {
            const response = await this.post(`/api/members/${id}/points/deduct`, pointsData);
            return response;
        } catch (error) {
            throw error;
        }
    }

    // 计算会员等级
    async calculateMemberLevel(id) {
        try {
            const response = await this.post(`/api/members/${id}/level/calculate`);
            return response;
        } catch (error) {
            throw error;
        }
    }
}
