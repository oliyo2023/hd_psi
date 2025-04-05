// 商品列表组件
const ProductList = {
    template: `
        <div>
            <div class="table-container">
                <div class="table-header">
                    <div class="table-title">商品列表</div>
                    <el-button type="primary" @click="handleAdd">添加商品</el-button>
                </div>
                
                <el-form :inline="true" :model="searchForm" class="demo-form-inline" style="margin-bottom: 20px;">
                    <el-form-item label="商品名称">
                        <el-input v-model="searchForm.name" placeholder="商品名称" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="SKU">
                        <el-input v-model="searchForm.sku" placeholder="SKU" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="类别">
                        <el-select v-model="searchForm.category" placeholder="类别" clearable>
                            <el-option label="外套" value="外套"></el-option>
                            <el-option label="裤装" value="裤装"></el-option>
                            <el-option label="衬衫" value="衬衫"></el-option>
                            <el-option label="T恤" value="T恤"></el-option>
                            <el-option label="内衣" value="内衣"></el-option>
                            <el-option label="配饰" value="配饰"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="handleSearch">查询</el-button>
                        <el-button @click="resetSearch">重置</el-button>
                    </el-form-item>
                </el-form>
                
                <el-table
                    v-loading="loading"
                    :data="products"
                    style="width: 100%"
                    border
                >
                    <el-table-column prop="id" label="ID" width="80"></el-table-column>
                    <el-table-column prop="sku" label="SKU" width="150"></el-table-column>
                    <el-table-column prop="name" label="商品名称"></el-table-column>
                    <el-table-column prop="category" label="类别"></el-table-column>
                    <el-table-column prop="color" label="颜色"></el-table-column>
                    <el-table-column prop="size" label="尺码"></el-table-column>
                    <el-table-column prop="season" label="季节"></el-table-column>
                    <el-table-column prop="retailPrice" label="零售价">
                        <template #default="scope">
                            ¥{{ scope.row.retailPrice.toFixed(2) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="200">
                        <template #default="scope">
                            <el-button
                                size="small"
                                @click="handleEdit(scope.row)"
                            >
                                编辑
                            </el-button>
                            <el-button
                                size="small"
                                type="danger"
                                @click="handleDelete(scope.row)"
                            >
                                删除
                            </el-button>
                        </template>
                    </el-table-column>
                </el-table>
                
                <div style="margin-top: 20px; display: flex; justify-content: flex-end;">
                    <el-pagination
                        v-model:current-page="currentPage"
                        v-model:page-size="pageSize"
                        :page-sizes="[10, 20, 50, 100]"
                        layout="total, sizes, prev, pager, next, jumper"
                        :total="total"
                        @size-change="handleSizeChange"
                        @current-change="handleCurrentChange"
                    ></el-pagination>
                </div>
            </div>
            
            <!-- 商品表单对话框 -->
            <el-dialog
                v-model="dialogVisible"
                :title="dialogType === 'add' ? '添加商品' : '编辑商品'"
                width="500px"
            >
                <el-form
                    ref="productFormRef"
                    :model="productForm"
                    :rules="productRules"
                    label-width="100px"
                >
                    <el-form-item label="SKU" prop="sku">
                        <el-input v-model="productForm.sku" placeholder="请输入SKU"></el-input>
                    </el-form-item>
                    <el-form-item label="商品名称" prop="name">
                        <el-input v-model="productForm.name" placeholder="请输入商品名称"></el-input>
                    </el-form-item>
                    <el-form-item label="类别" prop="category">
                        <el-select v-model="productForm.category" placeholder="请选择类别" style="width: 100%;">
                            <el-option label="外套" value="外套"></el-option>
                            <el-option label="裤装" value="裤装"></el-option>
                            <el-option label="衬衫" value="衬衫"></el-option>
                            <el-option label="T恤" value="T恤"></el-option>
                            <el-option label="内衣" value="内衣"></el-option>
                            <el-option label="配饰" value="配饰"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="颜色" prop="color">
                        <el-input v-model="productForm.color" placeholder="请输入颜色"></el-input>
                    </el-form-item>
                    <el-form-item label="尺码" prop="size">
                        <el-input v-model="productForm.size" placeholder="请输入尺码"></el-input>
                    </el-form-item>
                    <el-form-item label="季节" prop="season">
                        <el-select v-model="productForm.season" placeholder="请选择季节" style="width: 100%;">
                            <el-option label="春季" value="春季"></el-option>
                            <el-option label="夏季" value="夏季"></el-option>
                            <el-option label="秋季" value="秋季"></el-option>
                            <el-option label="冬季" value="冬季"></el-option>
                            <el-option label="四季" value="四季"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="成本价" prop="costPrice">
                        <el-input-number v-model="productForm.costPrice" :precision="2" :step="0.01" :min="0" style="width: 100%;"></el-input-number>
                    </el-form-item>
                    <el-form-item label="零售价" prop="retailPrice">
                        <el-input-number v-model="productForm.retailPrice" :precision="2" :step="0.01" :min="0" style="width: 100%;"></el-input-number>
                    </el-form-item>
                    <el-form-item label="图片" prop="image">
                        <el-input v-model="productForm.image" placeholder="请输入图片URL"></el-input>
                    </el-form-item>
                </el-form>
                <template #footer>
                    <span class="dialog-footer">
                        <el-button @click="dialogVisible = false">取消</el-button>
                        <el-button type="primary" @click="submitForm">确定</el-button>
                    </span>
                </template>
            </el-dialog>
        </div>
    `,
    setup() {
        // 响应式状态
        const loading = Vue.ref(false);
        const products = Vue.ref([]);
        const currentPage = Vue.ref(1);
        const pageSize = Vue.ref(10);
        const total = Vue.ref(0);
        const dialogVisible = Vue.ref(false);
        const dialogType = Vue.ref('add'); // 'add' 或 'edit'
        const productFormRef = Vue.ref(null);
        
        // 搜索表单
        const searchForm = Vue.reactive({
            name: '',
            sku: '',
            category: ''
        });
        
        // 商品表单
        const productForm = Vue.reactive({
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
        });
        
        // 表单验证规则
        const productRules = {
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
                { required: true, message: '请输入零售价', trigger: 'blur' }
            ]
        };
        
        // 方法
        const loadProducts = async () => {
            loading.value = true;
            try {
                const productService = new ProductService();
                const params = {
                    page: currentPage.value,
                    limit: pageSize.value,
                    ...searchForm
                };
                const response = await productService.getProducts(params);
                products.value = response.items || [];
                total.value = response.total || 0;
            } catch (error) {
                console.error('加载商品列表失败:', error);
                ElementPlus.ElMessage.error('加载商品列表失败');
            } finally {
                loading.value = false;
            }
        };
        
        const handleSearch = () => {
            currentPage.value = 1;
            loadProducts();
        };
        
        const resetSearch = () => {
            Object.keys(searchForm).forEach(key => {
                searchForm[key] = '';
            });
            currentPage.value = 1;
            loadProducts();
        };
        
        const handleSizeChange = (val) => {
            pageSize.value = val;
            loadProducts();
        };
        
        const handleCurrentChange = (val) => {
            currentPage.value = val;
            loadProducts();
        };
        
        const handleAdd = () => {
            dialogType.value = 'add';
            resetForm();
            dialogVisible.value = true;
        };
        
        const handleEdit = (row) => {
            dialogType.value = 'edit';
            resetForm();
            Object.keys(productForm).forEach(key => {
                if (key in row) {
                    productForm[key] = row[key];
                }
            });
            dialogVisible.value = true;
        };
        
        const handleDelete = (row) => {
            ElementPlus.ElMessageBox.confirm(
                `确定要删除商品 "${row.name}" 吗？`,
                '警告',
                {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }
            ).then(async () => {
                try {
                    const productService = new ProductService();
                    await productService.deleteProduct(row.id);
                    ElementPlus.ElMessage.success('删除成功');
                    loadProducts();
                } catch (error) {
                    console.error('删除商品失败:', error);
                    ElementPlus.ElMessage.error('删除商品失败');
                }
            }).catch(() => {
                // 取消删除
            });
        };
        
        const resetForm = () => {
            Object.keys(productForm).forEach(key => {
                if (key === 'id') {
                    productForm[key] = null;
                } else if (key === 'costPrice' || key === 'retailPrice') {
                    productForm[key] = 0;
                } else {
                    productForm[key] = '';
                }
            });
            if (productFormRef.value) {
                productFormRef.value.resetFields();
            }
        };
        
        const submitForm = () => {
            productFormRef.value.validate(async (valid) => {
                if (valid) {
                    try {
                        const productService = new ProductService();
                        if (dialogType.value === 'add') {
                            await productService.createProduct(productForm);
                            ElementPlus.ElMessage.success('添加成功');
                        } else {
                            await productService.updateProduct(productForm.id, productForm);
                            ElementPlus.ElMessage.success('更新成功');
                        }
                        dialogVisible.value = false;
                        loadProducts();
                    } catch (error) {
                        console.error('保存商品失败:', error);
                        ElementPlus.ElMessage.error('保存商品失败');
                    }
                } else {
                    return false;
                }
            });
        };
        
        // 生命周期钩子
        Vue.onMounted(() => {
            loadProducts();
        });
        
        return {
            loading,
            products,
            currentPage,
            pageSize,
            total,
            dialogVisible,
            dialogType,
            productFormRef,
            searchForm,
            productForm,
            productRules,
            handleSearch,
            resetSearch,
            handleSizeChange,
            handleCurrentChange,
            handleAdd,
            handleEdit,
            handleDelete,
            submitForm
        };
    }
};
