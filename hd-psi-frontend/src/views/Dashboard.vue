<template>
  <div class="dashboard-container">
    <n-grid :cols="24" :x-gap="16">
      <!-- 统计卡片 -->
      <n-grid-item :span="6">
        <n-card class="stat-card">
          <div class="stat-icon">
            <n-icon size="48" color="#18a058">
              <CubeOutline />
            </n-icon>
          </div>
          <div class="stat-content">
            <div class="stat-title">商品总数</div>
            <div class="stat-value">{{ statistics.productCount }}</div>
          </div>
        </n-card>
      </n-grid-item>
      
      <n-grid-item :span="6">
        <n-card class="stat-card">
          <div class="stat-icon">
            <n-icon size="48" color="#2080f0">
              <StorefrontOutline />
            </n-icon>
          </div>
          <div class="stat-content">
            <div class="stat-title">库存总量</div>
            <div class="stat-value">{{ statistics.inventoryCount }}</div>
          </div>
        </n-card>
      </n-grid-item>
      
      <n-grid-item :span="6">
        <n-card class="stat-card">
          <div class="stat-icon">
            <n-icon size="48" color="#f0a020">
              <PeopleOutline />
            </n-icon>
          </div>
          <div class="stat-content">
            <div class="stat-title">会员总数</div>
            <div class="stat-value">{{ statistics.memberCount }}</div>
          </div>
        </n-card>
      </n-grid-item>
      
      <n-grid-item :span="6">
        <n-card class="stat-card">
          <div class="stat-icon">
            <n-icon size="48" color="#d03050">
              <CartOutline />
            </n-icon>
          </div>
          <div class="stat-content">
            <div class="stat-title">销售订单</div>
            <div class="stat-value">{{ statistics.orderCount }}</div>
          </div>
        </n-card>
      </n-grid-item>
      
      <!-- 库存预警 -->
      <n-grid-item :span="12">
        <n-card title="库存预警" class="data-card">
          <template #header-extra>
            <n-button size="small" @click="checkInventoryLevels">
              检查库存水平
            </n-button>
          </template>
          
          <n-data-table
            :columns="alertColumns"
            :data="inventoryAlerts"
            :pagination="{ pageSize: 5 }"
            :bordered="false"
            size="small"
          />
        </n-card>
      </n-grid-item>
      
      <!-- 最近销售 -->
      <n-grid-item :span="12">
        <n-card title="最近销售" class="data-card">
          <n-data-table
            :columns="salesColumns"
            :data="recentSales"
            :pagination="{ pageSize: 5 }"
            :bordered="false"
            size="small"
          />
        </n-card>
      </n-grid-item>
    </n-grid>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h } from 'vue'
import { 
  NGrid, NGridItem, NCard, NIcon, NDataTable, 
  NButton, NTag, useMessage
} from 'naive-ui'
import { 
  CubeOutline, StorefrontOutline, 
  PeopleOutline, CartOutline
} from '@vicons/ionicons5'

const message = useMessage()

// 统计数据
const statistics = reactive({
  productCount: 0,
  inventoryCount: 0,
  memberCount: 0,
  orderCount: 0
})

// 库存预警数据
const inventoryAlerts = ref([])

// 最近销售数据
const recentSales = ref([])

// 库存预警表格列
const alertColumns = [
  {
    title: '商品名称',
    key: 'productName',
    width: 150
  },
  {
    title: '店铺',
    key: 'storeName',
    width: 100
  },
  {
    title: '预警类型',
    key: 'alertType',
    width: 100,
    render(row) {
      return h(
        NTag,
        {
          type: row.alertType === 'low' ? 'error' : 'warning',
          size: 'small'
        },
        { default: () => row.alertType === 'low' ? '库存不足' : '库存过多' }
      )
    }
  },
  {
    title: '当前库存',
    key: 'currentQty',
    width: 100
  },
  {
    title: '阈值',
    key: 'threshold',
    width: 100
  }
]

// 最近销售表格列
const salesColumns = [
  {
    title: '订单编号',
    key: 'orderNumber',
    width: 150
  },
  {
    title: '店铺',
    key: 'storeName',
    width: 100
  },
  {
    title: '金额',
    key: 'amount',
    width: 100,
    render(row) {
      return `¥${row.amount.toFixed(2)}`
    }
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      const statusMap = {
        completed: { type: 'success', text: '已完成' },
        processing: { type: 'info', text: '处理中' },
        pending: { type: 'warning', text: '待处理' },
        cancelled: { type: 'error', text: '已取消' }
      }
      const status = statusMap[row.status] || { type: 'default', text: '未知' }
      
      return h(
        NTag,
        {
          type: status.type,
          size: 'small'
        },
        { default: () => status.text }
      )
    }
  },
  {
    title: '创建时间',
    key: 'createdAt',
    width: 150
  }
]

// 加载统计数据
const loadStatistics = () => {
  // 模拟数据，实际应该从API获取
  statistics.productCount = 120
  statistics.inventoryCount = 1580
  statistics.memberCount = 350
  statistics.orderCount = 68
}

// 加载库存预警数据
const loadInventoryAlerts = () => {
  // 模拟数据，实际应该从API获取
  inventoryAlerts.value = [
    {
      id: 1,
      productName: '男士休闲衬衫',
      storeName: '总店',
      alertType: 'low',
      currentQty: 5,
      threshold: 10
    },
    {
      id: 2,
      productName: '女士连衣裙',
      storeName: '分店1',
      alertType: 'low',
      currentQty: 3,
      threshold: 8
    },
    {
      id: 3,
      productName: '儿童T恤',
      storeName: '总店',
      alertType: 'high',
      currentQty: 100,
      threshold: 50
    },
    {
      id: 4,
      productName: '男士牛仔裤',
      storeName: '分店2',
      alertType: 'low',
      currentQty: 7,
      threshold: 15
    },
    {
      id: 5,
      productName: '女士外套',
      storeName: '总店',
      alertType: 'high',
      currentQty: 80,
      threshold: 40
    }
  ]
}

// 加载最近销售数据
const loadRecentSales = () => {
  // 模拟数据，实际应该从API获取
  recentSales.value = [
    {
      id: 1,
      orderNumber: 'SO20230405001',
      storeName: '总店',
      amount: 1299.00,
      status: 'completed',
      createdAt: '2023-04-05 10:30:00'
    },
    {
      id: 2,
      orderNumber: 'SO20230404002',
      storeName: '分店1',
      amount: 899.50,
      status: 'completed',
      createdAt: '2023-04-04 15:20:00'
    },
    {
      id: 3,
      orderNumber: 'SO20230404001',
      storeName: '分店2',
      amount: 599.00,
      status: 'processing',
      createdAt: '2023-04-04 09:15:00'
    },
    {
      id: 4,
      orderNumber: 'SO20230403001',
      storeName: '总店',
      amount: 1599.00,
      status: 'completed',
      createdAt: '2023-04-03 14:45:00'
    },
    {
      id: 5,
      orderNumber: 'SO20230402001',
      storeName: '分店1',
      amount: 799.00,
      status: 'completed',
      createdAt: '2023-04-02 11:30:00'
    }
  ]
}

// 检查库存水平
const checkInventoryLevels = () => {
  // 实际应该调用API
  message.success('库存水平检查完成')
  loadInventoryAlerts()
}

// 生命周期钩子
onMounted(() => {
  loadStatistics()
  loadInventoryAlerts()
  loadRecentSales()
})
</script>

<style scoped>
.dashboard-container {
  padding: 16px;
}

.stat-card {
  height: 120px;
  display: flex;
  align-items: center;
}

.stat-icon {
  margin-right: 16px;
}

.stat-content {
  flex: 1;
}

.stat-title {
  font-size: 16px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.data-card {
  margin-top: 16px;
  height: 400px;
}
</style>
