<template>
  <div class="product-list">
    <div class="page-header">
      <h1 class="page-title">商品管理</h1>
      <n-button type="primary" @click="handleAddProduct">
        添加商品
      </n-button>
    </div>

    <div class="page-content">
      <!-- 搜索工具栏 -->
      <div class="table-toolbar">
        <div class="table-search">
          <n-input
            v-model:value="searchForm.name"
            placeholder="商品名称"
            clearable
            style="width: 200px"
          />
          <n-input
            v-model:value="searchForm.sku"
            placeholder="SKU"
            clearable
            style="width: 150px"
          />
          <n-select
            v-model:value="searchForm.category"
            placeholder="类别"
            clearable
            :options="categoryOptions"
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

      <!-- 商品表格 -->
      <n-data-table
        ref="tableRef"
        :columns="columns"
        :data="products"
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
  NButton, NDataTable, NInput, NSelect, NSpace
} from 'naive-ui'

// 响应式状态
const loading = ref(false)
const products = ref([])
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
  sku: '',
  category: null
})

// 类别选项
const categoryOptions = [
  { label: '外套', value: '外套' },
  { label: '裤装', value: '裤装' },
  { label: '衬衫', value: '衬衫' },
  { label: 'T恤', value: 'T恤' },
  { label: '内衣', value: '内衣' },
  { label: '配饰', value: '配饰' }
]

// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 80
  },
  {
    title: 'SKU',
    key: 'sku',
    width: 150
  },
  {
    title: '商品名称',
    key: 'name',
    width: 200
  },
  {
    title: '类别',
    key: 'category',
    width: 100
  },
  {
    title: '颜色',
    key: 'color',
    width: 100
  },
  {
    title: '尺码',
    key: 'size',
    width: 100
  },
  {
    title: '季节',
    key: 'season',
    width: 100
  },
  {
    title: '零售价',
    key: 'retailPrice',
    width: 120,
    render(row) {
      return `¥${row.retailPrice?.toFixed(2) || '0.00'}`
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
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
const loadProducts = async () => {
  loading.value = true
  try {
    // 模拟数据，实际应该从API获取
    products.value = [
      {
        id: 1,
        sku: 'MS001',
        name: '男士休闲衬衫',
        category: '衬衫',
        color: '白色',
        size: 'L',
        season: '春季',
        costPrice: 89.00,
        retailPrice: 199.00
      },
      {
        id: 2,
        sku: 'WD001',
        name: '女士连衣裙',
        category: '裤装',
        color: '蓝色',
        size: 'M',
        season: '夏季',
        costPrice: 120.00,
        retailPrice: 299.00
      },
      {
        id: 3,
        sku: 'MT001',
        name: '男士T恤',
        category: 'T恤',
        color: '黑色',
        size: 'XL',
        season: '夏季',
        costPrice: 45.00,
        retailPrice: 99.00
      }
    ]
    pagination.itemCount = products.value.length
  } catch (error) {
    console.error('加载商品列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  loadProducts()
}

const resetSearch = () => {
  searchForm.name = ''
  searchForm.sku = ''
  searchForm.category = null
  pagination.page = 1
  loadProducts()
}

const handlePageChange = (page) => {
  pagination.page = page
  loadProducts()
}

const handlePageSizeChange = (pageSize) => {
  pagination.pageSize = pageSize
  pagination.page = 1
  loadProducts()
}

const handleAddProduct = () => {
  alert('添加商品功能尚未实现')
}

const handleEdit = (row) => {
  alert(`编辑商品: ${row.name}`)
}

const handleDelete = (row) => {
  alert(`删除商品: ${row.name}`)
}

// 生命周期钩子
onMounted(() => {
  loadProducts()
})
</script>

<style scoped>
.product-list {
  padding: 16px;
}
</style>
