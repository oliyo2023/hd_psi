<template>
  <div class="product-list">
    <div class="page-header">
      <h1 class="page-title">商品管理</h1>
      <n-button type="primary" @click="handleAddProduct">
        <template #icon>
          <n-icon><AddOutline /></n-icon>
        </template>
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
            <template #icon>
              <n-icon><SearchOutline /></n-icon>
            </template>
            查询
          </n-button>
          <n-button @click="resetSearch">
            <template #icon>
              <n-icon><RefreshOutline /></n-icon>
            </template>
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
    
    <!-- 商品表单对话框 -->
    <n-modal
      v-model:show="showModal"
      :title="formTitle"
      preset="card"
      style="width: 600px"
      :mask-closable="false"
    >
      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-placement="left"
        label-width="100px"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="SKU" path="sku">
          <n-input v-model:value="formData.sku" placeholder="请输入SKU" />
        </n-form-item>
        
        <n-form-item label="商品名称" path="name">
          <n-input v-model:value="formData.name" placeholder="请输入商品名称" />
        </n-form-item>
        
        <n-form-item label="类别" path="category">
          <n-select
            v-model:value="formData.category"
            placeholder="请选择类别"
            :options="categoryOptions"
          />
        </n-form-item>
        
        <n-form-item label="颜色" path="color">
          <n-input v-model:value="formData.color" placeholder="请输入颜色" />
        </n-form-item>
        
        <n-form-item label="尺码" path="size">
          <n-input v-model:value="formData.size" placeholder="请输入尺码" />
        </n-form-item>
        
        <n-form-item label="季节" path="season">
          <n-select
            v-model:value="formData.season"
            placeholder="请选择季节"
            :options="seasonOptions"
          />
        </n-form-item>
        
        <n-form-item label="成本价" path="costPrice">
          <n-input-number
            v-model:value="formData.costPrice"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </n-form-item>
        
        <n-form-item label="零售价" path="retailPrice">
          <n-input-number
            v-model:value="formData.retailPrice"
            :min="0"
            :precision="2"
            style="width: 100%"
          />
        </n-form-item>
        
        <n-form-item label="图片" path="image">
          <n-input v-model:value="formData.image" placeholder="请输入图片URL" />
        </n-form-item>
      </n-form>
      
      <template #footer>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed, h, onMounted } from 'vue'
import { 
  NButton, NDataTable, NInput, NSelect, NSpace, 
  NModal, NForm, NFormItem, NInputNumber, NIcon,
  useMessage, useDialog
} from 'naive-ui'
import { 
  AddOutline, SearchOutline, RefreshOutline,
  CreateOutline, TrashOutline
} from '@vicons/ionicons5'
import productService from '../services/product'

const message = useMessage()
const dialog = useDialog()

// 响应式状态
const loading = ref(false)
const products = ref([])
const showModal = ref(false)
const formMode = ref('add') // 'add' 或 'edit'
const formRef = ref(null)
const formData = reactive({
  id: null,
  sku: '',
  name: '',
  category: '',
  color: '',
  size: '',
  season: '',
  costPrice: 0,
  retailPrice: 0,
  image: ''
})
const submitting = ref(false)
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

// 计算属性
const formTitle = computed(() => {
  return formMode.value === 'add' ? '添加商品' : '编辑商品'
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

// 季节选项
const seasonOptions = [
  { label: '春季', value: '春季' },
  { label: '夏季', value: '夏季' },
  { label: '秋季', value: '秋季' },
  { label: '冬季', value: '冬季' },
  { label: '四季', value: '四季' }
]

// 表单验证规则
const rules = {
  sku: [
    { required: true, message: '请输入SKU', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择类别', trigger: 'change' }
  ],
  color: [
    { required: true, message: '请输入颜色', trigger: 'blur' }
  ],
  size: [
    { required: true, message: '请输入尺码', trigger: 'blur' }
  ],
  season: [
    { required: true, message: '请选择季节', trigger: 'change' }
  ],
  retailPrice: [
    { required: true, type: 'number', message: '请输入零售价', trigger: 'blur' }
  ]
}

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
      return `¥${row.retailPrice.toFixed(2)}`
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
            { icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              onClick: () => handleDelete(row)
            },
            { icon: () => h(NIcon, null, { default: () => h(TrashOutline) }) }
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
    const params = {
      page: pagination.page,
      limit: pagination.pageSize,
      ...searchForm
    }
    const response = await productService.getProducts(params)
    products.value = response.items || []
    pagination.itemCount = response.total || 0
  } catch (error) {
    console.error('加载商品列表失败:', error)
    message.error('加载商品列表失败')
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
  formMode.value = 'add'
  resetForm()
  showModal.value = true
}

const handleEdit = (row) => {
  formMode.value = 'edit'
  resetForm()
  Object.keys(formData).forEach(key => {
    if (key in row) {
      formData[key] = row[key]
    }
  })
  showModal.value = true
}

const handleDelete = (row) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除商品 "${row.name}" 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await productService.deleteProduct(row.id)
        message.success('删除成功')
        loadProducts()
      } catch (error) {
        console.error('删除商品失败:', error)
        message.error('删除商品失败')
      }
    }
  })
}

const resetForm = () => {
  formData.id = null
  formData.sku = ''
  formData.name = ''
  formData.category = ''
  formData.color = ''
  formData.size = ''
  formData.season = ''
  formData.costPrice = 0
  formData.retailPrice = 0
  formData.image = ''
  
  if (formRef.value) {
    formRef.value.restoreValidation()
  }
}

const handleSubmit = () => {
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      submitting.value = true
      try {
        if (formMode.value === 'add') {
          await productService.createProduct(formData)
          message.success('添加成功')
        } else {
          await productService.updateProduct(formData.id, formData)
          message.success('更新成功')
        }
        showModal.value = false
        loadProducts()
      } catch (error) {
        console.error('保存商品失败:', error)
        message.error('保存商品失败')
      } finally {
        submitting.value = false
      }
    }
  })
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
