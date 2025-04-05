// 仪表盘组件
const Dashboard = {
    template: `
        <div>
            <el-row :gutter="20">
                <el-col :span="6">
                    <div class="dashboard-card">
                        <div class="dashboard-card-title">商品总数</div>
                        <div class="dashboard-card-value">{{ statistics.productCount }}</div>
                        <div class="dashboard-card-icon">
                            <el-icon><el-icon-goods /></el-icon>
                        </div>
                    </div>
                </el-col>
                
                <el-col :span="6">
                    <div class="dashboard-card">
                        <div class="dashboard-card-title">库存总量</div>
                        <div class="dashboard-card-value">{{ statistics.inventoryCount }}</div>
                        <div class="dashboard-card-icon">
                            <el-icon><el-icon-box /></el-icon>
                        </div>
                    </div>
                </el-col>
                
                <el-col :span="6">
                    <div class="dashboard-card">
                        <div class="dashboard-card-title">会员总数</div>
                        <div class="dashboard-card-value">{{ statistics.memberCount }}</div>
                        <div class="dashboard-card-icon">
                            <el-icon><el-icon-user /></el-icon>
                        </div>
                    </div>
                </el-col>
                
                <el-col :span="6">
                    <div class="dashboard-card">
                        <div class="dashboard-card-title">销售订单</div>
                        <div class="dashboard-card-value">{{ statistics.orderCount }}</div>
                        <div class="dashboard-card-icon">
                            <el-icon><el-icon-shopping-cart /></el-icon>
                        </div>
                    </div>
                </el-col>
            </el-row>
            
            <el-row :gutter="20" style="margin-top: 20px;">
                <el-col :span="12">
                    <div class="table-container">
                        <div class="table-header">
                            <div class="table-title">库存预警</div>
                            <el-button type="primary" size="small" @click="checkInventoryLevels">
                                检查库存水平
                            </el-button>
                        </div>
                        
                        <el-table :data="inventoryAlerts" style="width: 100%">
                            <el-table-column prop="product.name" label="商品名称"></el-table-column>
                            <el-table-column prop="store.name" label="店铺"></el-table-column>
                            <el-table-column prop="alertType" label="预警类型">
                                <template #default="scope">
                                    <el-tag :type="scope.row.alertType === 'low' ? 'danger' : 'warning'">
                                        {{ scope.row.alertType === 'low' ? '库存不足' : '库存过多' }}
                                    </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column prop="currentQty" label="当前库存"></el-table-column>
                            <el-table-column prop="threshold" label="阈值"></el-table-column>
                        </el-table>
                    </div>
                </el-col>
                
                <el-col :span="12">
                    <div class="table-container">
                        <div class="table-header">
                            <div class="table-title">最近销售</div>
                        </div>
                        
                        <el-table :data="recentSales" style="width: 100%">
                            <el-table-column prop="orderNumber" label="订单编号"></el-table-column>
                            <el-table-column prop="store.name" label="店铺"></el-table-column>
                            <el-table-column prop="actualAmount" label="金额">
                                <template #default="scope">
                                    ¥{{ scope.row.actualAmount.toFixed(2) }}
                                </template>
                            </el-table-column>
                            <el-table-column prop="status" label="状态">
                                <template #default="scope">
                                    <el-tag :type="getOrderStatusType(scope.row.status)">
                                        {{ getOrderStatusText(scope.row.status) }}
                                    </el-tag>
                                </template>
                            </el-table-column>
                            <el-table-column prop="createdAt" label="创建时间">
                                <template #default="scope">
                                    {{ formatDate(scope.row.createdAt) }}
                                </template>
                            </el-table-column>
                        </el-table>
                    </div>
                </el-col>
            </el-row>
        </div>
    `,
    setup() {
        // 响应式状态
        const statistics = Vue.reactive({
            productCount: 0,
            inventoryCount: 0,
            memberCount: 0,
            orderCount: 0
        });
        
        const inventoryAlerts = Vue.ref([]);
        const recentSales = Vue.ref([]);
        
        // 方法
        const loadStatistics = async () => {
            try {
                // 这里应该调用后端API获取统计数据
                // 暂时使用模拟数据
                statistics.productCount = 120;
                statistics.inventoryCount = 1580;
                statistics.memberCount = 350;
                statistics.orderCount = 68;
            } catch (error) {
                console.error('加载统计数据失败:', error);
                ElementPlus.ElMessage.error('加载统计数据失败');
            }
        };
        
        const loadInventoryAlerts = async () => {
            try {
                const inventoryService = new InventoryService();
                const response = await inventoryService.getInventoryAlerts({ status: 'active', limit: 5 });
                inventoryAlerts.value = response;
            } catch (error) {
                console.error('加载库存预警失败:', error);
                ElementPlus.ElMessage.error('加载库存预警失败');
            }
        };
        
        const loadRecentSales = async () => {
            try {
                // 这里应该调用后端API获取最近销售数据
                // 暂时使用模拟数据
                recentSales.value = [
                    {
                        orderNumber: 'SO20230405001',
                        store: { name: '总店' },
                        actualAmount: 1299.00,
                        status: 'completed',
                        createdAt: '2023-04-05T10:30:00Z'
                    },
                    {
                        orderNumber: 'SO20230404002',
                        store: { name: '分店1' },
                        actualAmount: 899.50,
                        status: 'completed',
                        createdAt: '2023-04-04T15:20:00Z'
                    },
                    {
                        orderNumber: 'SO20230404001',
                        store: { name: '分店2' },
                        actualAmount: 599.00,
                        status: 'processing',
                        createdAt: '2023-04-04T09:15:00Z'
                    },
                    {
                        orderNumber: 'SO20230403001',
                        store: { name: '总店' },
                        actualAmount: 1599.00,
                        status: 'completed',
                        createdAt: '2023-04-03T14:45:00Z'
                    },
                    {
                        orderNumber: 'SO20230402001',
                        store: { name: '分店1' },
                        actualAmount: 799.00,
                        status: 'completed',
                        createdAt: '2023-04-02T11:30:00Z'
                    }
                ];
            } catch (error) {
                console.error('加载最近销售失败:', error);
                ElementPlus.ElMessage.error('加载最近销售失败');
            }
        };
        
        const checkInventoryLevels = async () => {
            try {
                const inventoryService = new InventoryService();
                await inventoryService.checkInventoryLevels();
                await loadInventoryAlerts();
                ElementPlus.ElMessage.success('库存水平检查完成');
            } catch (error) {
                console.error('库存水平检查失败:', error);
                ElementPlus.ElMessage.error('库存水平检查失败');
            }
        };
        
        const getOrderStatusType = (status) => {
            switch (status) {
                case 'completed':
                    return 'success';
                case 'processing':
                    return 'primary';
                case 'pending':
                    return 'warning';
                case 'cancelled':
                    return 'danger';
                default:
                    return 'info';
            }
        };
        
        const getOrderStatusText = (status) => {
            switch (status) {
                case 'completed':
                    return '已完成';
                case 'processing':
                    return '处理中';
                case 'pending':
                    return '待处理';
                case 'cancelled':
                    return '已取消';
                default:
                    return '未知';
            }
        };
        
        const formatDate = (dateString) => {
            const date = new Date(dateString);
            return date.toLocaleString('zh-CN', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });
        };
        
        // 生命周期钩子
        Vue.onMounted(() => {
            loadStatistics();
            loadInventoryAlerts();
            loadRecentSales();
        });
        
        return {
            statistics,
            inventoryAlerts,
            recentSales,
            checkInventoryLevels,
            getOrderStatusType,
            getOrderStatusText,
            formatDate
        };
    }
};
