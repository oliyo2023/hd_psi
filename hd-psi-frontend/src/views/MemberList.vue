<template>
  <div class="member-list">
    <div class="page-header">
      <h1 class="page-title">会员管理</h1>
      <n-button type="primary" @click="handleAddMember">
        添加会员
      </n-button>
    </div>
    
    <div class="page-content">
      <!-- 搜索工具栏 -->
      <div class="table-toolbar">
        <div class="table-search">
          <n-input
            v-model:value="searchForm.name"
            placeholder="会员姓名"
            clearable
            style="width: 200px"
          />
          <n-input
            v-model:value="searchForm.phone"
            placeholder="手机号码"
            clearable
            style="width: 150px"
          />
          <n-select
            v-model:value="searchForm.level"
            placeholder="会员等级"
            clearable
            :options="levelOptions"
            style="width: 150px"
          />
          <n-button type="primary" @click="handleSearch">
            查询
          </n-button>
          <n-button @click="resetSearch">
            重置
          </n-button>
        </div>
      </div>
      
      <!-- 会员表格 -->
      <n-data-table
        ref="tableRef"
        :columns="columns"
        :data="members"
        :loading="loading"
        :pagination="pagination"
        :row-key="row => row.id"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, h } from 'vue'
import { 
  NButton, NDataTable, NInput, NSelect, NSpace, NTag
} from 'naive-ui'

// 响应式状态
const loading = ref(false)
const members = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  pageSizes: [10, 20, 50, 100],
  showSizePicker: true,
  prefix({ itemCount }) {
    return `共 ${itemCount} 条`
  }
})

// 搜索表单
const searchForm = reactive({
  name: '',
  phone: '',
  level: null
})

// 会员等级选项
const levelOptions = [
  { label: '普通会员', value: '普通会员' },
  { label: '银卡会员', value: '银卡会员' },
  { label: '金卡会员', value: '金卡会员' },
  { label: '钻石会员', value: '钻石会员' }
]

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: '会员姓名',
    key: 'name',
    width: 120
  },
  {
    title: '手机号码',
    key: 'phone',
    width: 150
  },
  {
    title: '会员等级',
    key: 'level',
    width: 120,
    render(row) {
      const levelMap = {
        '普通会员': { type: 'default', text: '普通会员' },
        '银卡会员': { type: 'info', text: '银卡会员' },
        '金卡会员': { type: 'warning', text: '金卡会员' },
        '钻石会员': { type: 'success', text: '钻石会员' }
      }
      
      const level = levelMap[row.level] || levelMap['普通会员']
      
      return h(NTag, { type: level.type }, { default: () => level.text })
    }
  },
  {
    title: '积分',
    key: 'points',
    width: 100
  },
  {
    title: '累计消费',
    key: 'totalSpent',
    width: 120,
    render(row) {
      return `¥${row.totalSpent.toFixed(2)}`
    }
  },
  {
    title: '生日',
    key: 'birthday',
    width: 120
  },
  {
    title: '操作',
    key: 'actions',
    width: 250,
    fixed: 'right',
    render(row) {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              onClick: () => handleEdit(row)
            },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'info',
              onClick: () => handlePoints(row)
            },
            { default: () => '积分' }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => handleDelete(row)
            },
            { default: () => '删除' }
          )
        ]
      })
    }
  }
]

// 方法
const loadMembers = async () => {
  loading.value = true
  try {
    // 模拟数据，实际应该从API获取
    members.value = [
      {
        id: 1,
        name: '张三',
        phone: '13800138001',
        level: '金卡会员',
        points: 2500,
        totalSpent: 15000.00,
        birthday: '1990-01-15'
      },
      {
        id: 2,
        name: '李四',
        phone: '13800138002',
        level: '普通会员',
        points: 500,
        totalSpent: 2000.00,
        birthday: '1985-05-20'
      },
      {
        id: 3,
        name: '王五',
        phone: '13800138003',
        level: '银卡会员',
        points: 1200,
        totalSpent: 8000.00,
        birthday: '1992-11-08'
      },
      {
        id: 4,
        name: '赵六',
        phone: '13800138004',
        level: '钻石会员',
        points: 5000,
        totalSpent: 30000.00,
        birthday: '1988-07-30'
      }
    ]
    pagination.itemCount = members.value.length
  } catch (error) {
    console.error('加载会员列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadMembers()
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.phone = ''
  searchForm.level = null
  pagination.page = 1
  loadMembers()
}

const handlePageChange = (page) => {
  pagination.page = page
  loadMembers()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.page = 1
  loadMembers()
}

const handleAddMember = () => {
  alert('添加会员功能尚未实现')
}

const handleEdit = (row) => {
  alert(`编辑会员: ${row.name}`)
}

const handlePoints = (row) => {
  alert(`会员积分操作: ${row.name}`)
}

const handleDelete = (row) => {
  alert(`删除会员: ${row.name}`)
}

// 生命周期钩子
onMounted(() => {
  loadMembers()
})
</script>

<style scoped>
.member-list {
  padding: 16px;
}
</style>
