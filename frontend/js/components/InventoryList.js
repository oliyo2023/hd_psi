// 库存列表组件
const InventoryList = {
    template: `
        <div>
            <div class="table-container">
                <div class="table-header">
                    <div class="table-title">库存列表</div>
                    <el-button type="primary" @click="handleAdd">添加库存</el-button>
                </div>
                
                <el-form :inline="true" :model="searchForm" class="demo-form-inline" style="margin-bottom: 20px;">
                    <el-form-item label="商品名称">
                        <el-input v-model="searchForm.productName" placeholder="商品名称" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="SKU">
                        <el-input v-model="searchForm.sku" placeholder="SKU" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="店铺">
                        <el-select v-model="searchForm.storeId" placeholder="店铺" clearable>
                            <el-option
                                v-for="store in stores"
                                :key="store.id"
                                :label="store.name"
                                :value="store.id"
                            ></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="handleSearch">查询</el-button>
                        <el-button @click="resetSearch">重置</el-button>
                    </el-form-item>
                </el-form>
                
                <el-table
                    v-loading="loading"
                    :data="inventories"
                    style="width: 100%"
                    border
                >
                    <el-table-column prop="id" label="ID" width="80"></el-table-column>
                    <el-table-column prop="product.sku" label="SKU" width="150"></el-table-column>
                    <el-table-column prop="product.name" label="商品名称"></el-table-column>
                    <el-table-column prop="store.name" label="店铺"></el-table-column>
                    <el-table-column prop="quantity" label="库存数量"></el-table-column>
                    <el-table-column prop="location" label="库位"></el-table-column>
                    <el-table-column label="状态">
                        <template #default="scope">
                            <el-tag :type="getStatusType(scope.row.quantity)">
                                {{ getStatusText(scope.row.quantity) }}
                            </el-tag>
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
                                type="primary"
                                @click="handleTransaction(scope.row)"
                            >
                                调整
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
            
            <!-- 库存表单对话框 -->
            <el-dialog
                v-model="dialogVisible"
                :title="dialogType === 'add' ? '添加库存' : '编辑库存'"
                width="500px"
            >
                <el-form
                    ref="inventoryFormRef"
                    :model="inventoryForm"
                    :rules="inventoryRules"
                    label-width="100px"
                >
                    <el-form-item label="商品" prop="productId">
                        <el-select
                            v-model="inventoryForm.productId"
                            placeholder="请选择商品"
                            filterable
                            style="width: 100%;"
                        >
                            <el-option
                                v-for="product in products"
                                :key="product.id"
                                :label="product.name + ' (' + product.sku + ')'"
                                :value="product.id"
                            ></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="店铺" prop="storeId">
                        <el-select
                            v-model="inventoryForm.storeId"
                            placeholder="请选择店铺"
                            style="width: 100%;"
                        >
                            <el-option
                                v-for="store in stores"
                                :key="store.id"
                                :label="store.name"
                                :value="store.id"
                            ></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="库存数量" prop="quantity">
                        <el-input-number v-model="inventoryForm.quantity" :min="0" style="width: 100%;"></el-input-number>
                    </el-form-item>
                    <el-form-item label="库位" prop="location">
                        <el-input v-model="inventoryForm.location" placeholder="请输入库位"></el-input>
                    </el-form-item>
                </el-form>
                <template #footer>
                    <span class="dialog-footer">
                        <el-button @click="dialogVisible = false">取消</el-button>
                        <el-button type="primary" @click="submitForm">确定</el-button>
                    </span>
                </template>
            </el-dialog>
            
            <!-- 库存调整对话框 -->
            <el-dialog
                v-model="transactionDialogVisible"
                title="库存调整"
                width="500px"
            >
                <el-form
                    ref="transactionFormRef"
                    :model="transactionForm"
                    :rules="transactionRules"
                    label-width="100px"
                >
                    <el-form-item label="商品">
                        <el-input v-model="selectedInventory.product?.name" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="店铺">
                        <el-input v-model="selectedInventory.store?.name" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="当前库存">
                        <el-input v-model="selectedInventory.quantity" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="调整类型" prop="type">
                        <el-select v-model="transactionForm.type" placeholder="请选择调整类型" style="width: 100%;">
                            <el-option label="入库" value="in"></el-option>
                            <el-option label="出库" value="out"></el-option>
                            <el-option label="盘点调整" value="adjustment"></el-option>
                            <el-option label="损耗" value="loss"></el-option>
                            <el-option label="退货" value="return"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="调整数量" prop="quantity">
                        <el-input-number v-model="transactionForm.quantity" :min="1" style="width: 100%;"></el-input-number>
                    </el-form-item>
                    <el-form-item label="备注" prop="note">
                        <el-input v-model="transactionForm.note" type="textarea" placeholder="请输入备注"></el-input>
                    </el-form-item>
                </el-form>
                <template #footer>
                    <span class="dialog-footer">
                        <el-button @click="transactionDialogVisible = false">取消</el-button>
                        <el-button type="primary" @click="submitTransaction">确定</el-button>
                    </span>
                </template>
            </el-dialog>
        </div>
    `,
    setup() {
        // 响应式状态
        const loading = Vue.ref(false);
        const inventories = Vue.ref([]);
        const products = Vue.ref([]);
        const stores = Vue.ref([]);
        const currentPage = Vue.ref(1);
        const pageSize = Vue.ref(10);
        const total = Vue.ref(0);
        const dialogVisible = Vue.ref(false);
        const dialogType = Vue.ref('add'); // 'add' 或 'edit'
        const inventoryFormRef = Vue.ref(null);
        const transactionDialogVisible = Vue.ref(false);
        const transactionFormRef = Vue.ref(null);
        const selectedInventory = Vue.ref({});
        
        // 搜索表单
        const searchForm = Vue.reactive({
            productName: '',
            sku: '',
            storeId: ''
        });
        
        // 库存表单
        const inventoryForm = Vue.reactive({
            id: null,
            productId: null,
            storeId: null,
            quantity: 0,
            location: ''
        });
        
        // 库存调整表单
        const transactionForm = Vue.reactive({
            inventoryId: null,
            type: 'in',
            quantity: 1,
            note: ''
        });
        
        // 表单验证规则
        const inventoryRules = {
            productId: [
                { required: true, message: '请选择商品', trigger: 'change' }
            ],
            storeId: [
                { required: true, message: '请选择店铺', trigger: 'change' }
            ],
            quantity: [
                { required: true, message: '请输入库存数量', trigger: 'blur' }
            ]
        };
        
        // 库存调整表单验证规则
        const transactionRules = {
            type: [
                { required: true, message: '请选择调整类型', trigger: 'change' }
            ],
            quantity: [
                { required: true, message: '请输入调整数量', trigger: 'blur' }
            ]
        };
        
        // 方法
        const loadInventories = async () => {
            loading.value = true;
            try {
                const inventoryService = new InventoryService();
                const params = {
                    page: currentPage.value,
                    limit: pageSize.value,
                    ...searchForm
                };
                const response = await inventoryService.getInventories(params);
                inventories.value = response.items || [];
                total.value = response.total || 0;
            } catch (error) {
                console.error('加载库存列表失败:', error);
                ElementPlus.ElMessage.error('加载库存列表失败');
            } finally {
                loading.value = false;
            }
        };
        
        const loadProducts = async () => {
            try {
                const productService = new ProductService();
                const response = await productService.getProducts({ limit: 1000 });
                products.value = response.items || [];
            } catch (error) {
                console.error('加载商品列表失败:', error);
                ElementPlus.ElMessage.error('加载商品列表失败');
            }
        };
        
        const loadStores = async () => {
            try {
                // 这里应该调用后端API获取店铺列表
                // 暂时使用模拟数据
                stores.value = [
                    { id: 1, name: '总店' },
                    { id: 2, name: '分店1' },
                    { id: 3, name: '分店2' }
                ];
            } catch (error) {
                console.error('加载店铺列表失败:', error);
                ElementPlus.ElMessage.error('加载店铺列表失败');
            }
        };
        
        const handleSearch = () => {
            currentPage.value = 1;
            loadInventories();
        };
        
        const resetSearch = () => {
            Object.keys(searchForm).forEach(key => {
                searchForm[key] = '';
            });
            currentPage.value = 1;
            loadInventories();
        };
        
        const handleSizeChange = (val) => {
            pageSize.value = val;
            loadInventories();
        };
        
        const handleCurrentChange = (val) => {
            currentPage.value = val;
            loadInventories();
        };
        
        const handleAdd = () => {
            dialogType.value = 'add';
            resetForm();
            dialogVisible.value = true;
        };
        
        const handleEdit = (row) => {
            dialogType.value = 'edit';
            resetForm();
            inventoryForm.id = row.id;
            inventoryForm.productId = row.product?.id;
            inventoryForm.storeId = row.store?.id;
            inventoryForm.quantity = row.quantity;
            inventoryForm.location = row.location;
            dialogVisible.value = true;
        };
        
        const handleTransaction = (row) => {
            selectedInventory.value = row;
            resetTransactionForm();
            transactionForm.inventoryId = row.id;
            transactionDialogVisible.value = true;
        };
        
        const resetForm = () => {
            Object.keys(inventoryForm).forEach(key => {
                if (key === 'id' || key === 'productId' || key === 'storeId') {
                    inventoryForm[key] = null;
                } else if (key === 'quantity') {
                    inventoryForm[key] = 0;
                } else {
                    inventoryForm[key] = '';
                }
            });
            if (inventoryFormRef.value) {
                inventoryFormRef.value.resetFields();
            }
        };
        
        const resetTransactionForm = () => {
            transactionForm.inventoryId = null;
            transactionForm.type = 'in';
            transactionForm.quantity = 1;
            transactionForm.note = '';
            if (transactionFormRef.value) {
                transactionFormRef.value.resetFields();
            }
        };
        
        const submitForm = () => {
            inventoryFormRef.value.validate(async (valid) => {
                if (valid) {
                    try {
                        const inventoryService = new InventoryService();
                        if (dialogType.value === 'add') {
                            await inventoryService.createInventory(inventoryForm);
                            ElementPlus.ElMessage.success('添加成功');
                        } else {
                            await inventoryService.updateInventory(inventoryForm.id, inventoryForm);
                            ElementPlus.ElMessage.success('更新成功');
                        }
                        dialogVisible.value = false;
                        loadInventories();
                    } catch (error) {
                        console.error('保存库存失败:', error);
                        ElementPlus.ElMessage.error('保存库存失败');
                    }
                } else {
                    return false;
                }
            });
        };
        
        const submitTransaction = () => {
            transactionFormRef.value.validate(async (valid) => {
                if (valid) {
                    try {
                        const inventoryService = new InventoryService();
                        await inventoryService.createInventoryTransaction({
                            inventoryId: transactionForm.inventoryId,
                            type: transactionForm.type,
                            quantity: transactionForm.quantity,
                            note: transactionForm.note
                        });
                        ElementPlus.ElMessage.success('库存调整成功');
                        transactionDialogVisible.value = false;
                        loadInventories();
                    } catch (error) {
                        console.error('库存调整失败:', error);
                        ElementPlus.ElMessage.error('库存调整失败');
                    }
                } else {
                    return false;
                }
            });
        };
        
        const getStatusType = (quantity) => {
            if (quantity <= 0) {
                return 'danger';
            } else if (quantity <= 10) {
                return 'warning';
            } else {
                return 'success';
            }
        };
        
        const getStatusText = (quantity) => {
            if (quantity <= 0) {
                return '缺货';
            } else if (quantity <= 10) {
                return '低库存';
            } else {
                return '正常';
            }
        };
        
        // 生命周期钩子
        Vue.onMounted(() => {
            loadInventories();
            loadProducts();
            loadStores();
        });
        
        return {
            loading,
            inventories,
            products,
            stores,
            currentPage,
            pageSize,
            total,
            dialogVisible,
            dialogType,
            inventoryFormRef,
            transactionDialogVisible,
            transactionFormRef,
            selectedInventory,
            searchForm,
            inventoryForm,
            transactionForm,
            inventoryRules,
            transactionRules,
            handleSearch,
            resetSearch,
            handleSizeChange,
            handleCurrentChange,
            handleAdd,
            handleEdit,
            handleTransaction,
            submitForm,
            submitTransaction,
            getStatusType,
            getStatusText
        };
    }
};
