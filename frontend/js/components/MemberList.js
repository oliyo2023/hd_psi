// 会员列表组件
const MemberList = {
    template: `
        <div>
            <div class="table-container">
                <div class="table-header">
                    <div class="table-title">会员列表</div>
                    <el-button type="primary" @click="handleAdd">添加会员</el-button>
                </div>
                
                <el-form :inline="true" :model="searchForm" class="demo-form-inline" style="margin-bottom: 20px;">
                    <el-form-item label="会员名称">
                        <el-input v-model="searchForm.name" placeholder="会员名称" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="手机号码">
                        <el-input v-model="searchForm.phone" placeholder="手机号码" clearable></el-input>
                    </el-form-item>
                    <el-form-item label="会员等级">
                        <el-select v-model="searchForm.level" placeholder="会员等级" clearable>
                            <el-option label="普通会员" value="普通会员"></el-option>
                            <el-option label="银卡会员" value="银卡会员"></el-option>
                            <el-option label="金卡会员" value="金卡会员"></el-option>
                            <el-option label="钻石会员" value="钻石会员"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="handleSearch">查询</el-button>
                        <el-button @click="resetSearch">重置</el-button>
                    </el-form-item>
                </el-form>
                
                <el-table
                    v-loading="loading"
                    :data="members"
                    style="width: 100%"
                    border
                >
                    <el-table-column prop="id" label="ID" width="80"></el-table-column>
                    <el-table-column prop="name" label="会员姓名"></el-table-column>
                    <el-table-column prop="phone" label="手机号码"></el-table-column>
                    <el-table-column prop="level" label="会员等级">
                        <template #default="scope">
                            <el-tag :type="getMemberLevelType(scope.row.level)">
                                {{ scope.row.level }}
                            </el-tag>
                        </template>
                    </el-table-column>
                    <el-table-column prop="points" label="积分"></el-table-column>
                    <el-table-column prop="totalSpent" label="累计消费">
                        <template #default="scope">
                            ¥{{ scope.row.totalSpent.toFixed(2) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="birthday" label="生日">
                        <template #default="scope">
                            {{ formatDate(scope.row.birthday) }}
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="250">
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
                                @click="handlePoints(scope.row)"
                            >
                                积分
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
            
            <!-- 会员表单对话框 -->
            <el-dialog
                v-model="dialogVisible"
                :title="dialogType === 'add' ? '添加会员' : '编辑会员'"
                width="500px"
            >
                <el-form
                    ref="memberFormRef"
                    :model="memberForm"
                    :rules="memberRules"
                    label-width="100px"
                >
                    <el-form-item label="会员姓名" prop="name">
                        <el-input v-model="memberForm.name" placeholder="请输入会员姓名"></el-input>
                    </el-form-item>
                    <el-form-item label="手机号码" prop="phone">
                        <el-input v-model="memberForm.phone" placeholder="请输入手机号码"></el-input>
                    </el-form-item>
                    <el-form-item label="性别" prop="gender">
                        <el-radio-group v-model="memberForm.gender">
                            <el-radio label="男">男</el-radio>
                            <el-radio label="女">女</el-radio>
                            <el-radio label="其他">其他</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-form-item label="生日" prop="birthday">
                        <el-date-picker
                            v-model="memberForm.birthday"
                            type="date"
                            placeholder="请选择生日"
                            style="width: 100%;"
                        ></el-date-picker>
                    </el-form-item>
                    <el-form-item label="电子邮箱" prop="email">
                        <el-input v-model="memberForm.email" placeholder="请输入电子邮箱"></el-input>
                    </el-form-item>
                    <el-form-item label="地址" prop="address">
                        <el-input v-model="memberForm.address" placeholder="请输入地址"></el-input>
                    </el-form-item>
                    <el-form-item label="备注" prop="note">
                        <el-input v-model="memberForm.note" type="textarea" placeholder="请输入备注"></el-input>
                    </el-form-item>
                </el-form>
                <template #footer>
                    <span class="dialog-footer">
                        <el-button @click="dialogVisible = false">取消</el-button>
                        <el-button type="primary" @click="submitForm">确定</el-button>
                    </span>
                </template>
            </el-dialog>
            
            <!-- 积分操作对话框 -->
            <el-dialog
                v-model="pointsDialogVisible"
                title="积分操作"
                width="500px"
            >
                <el-form
                    ref="pointsFormRef"
                    :model="pointsForm"
                    :rules="pointsRules"
                    label-width="100px"
                >
                    <el-form-item label="会员姓名">
                        <el-input v-model="selectedMember.name" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="当前积分">
                        <el-input v-model="selectedMember.points" disabled></el-input>
                    </el-form-item>
                    <el-form-item label="操作类型" prop="type">
                        <el-radio-group v-model="pointsForm.type">
                            <el-radio label="add">增加积分</el-radio>
                            <el-radio label="deduct">扣减积分</el-radio>
                        </el-radio-group>
                    </el-form-item>
                    <el-form-item label="积分数量" prop="points">
                        <el-input-number v-model="pointsForm.points" :min="1" style="width: 100%;"></el-input-number>
                    </el-form-item>
                    <el-form-item label="原因" prop="reason">
                        <el-input v-model="pointsForm.reason" placeholder="请输入积分变动原因"></el-input>
                    </el-form-item>
                </el-form>
                <template #footer>
                    <span class="dialog-footer">
                        <el-button @click="pointsDialogVisible = false">取消</el-button>
                        <el-button type="primary" @click="submitPoints">确定</el-button>
                    </span>
                </template>
            </el-dialog>
        </div>
    `,
    setup() {
        // 响应式状态
        const loading = Vue.ref(false);
        const members = Vue.ref([]);
        const currentPage = Vue.ref(1);
        const pageSize = Vue.ref(10);
        const total = Vue.ref(0);
        const dialogVisible = Vue.ref(false);
        const dialogType = Vue.ref('add'); // 'add' 或 'edit'
        const memberFormRef = Vue.ref(null);
        const pointsDialogVisible = Vue.ref(false);
        const pointsFormRef = Vue.ref(null);
        const selectedMember = Vue.ref({});
        
        // 搜索表单
        const searchForm = Vue.reactive({
            name: '',
            phone: '',
            level: ''
        });
        
        // 会员表单
        const memberForm = Vue.reactive({
            id: null,
            name: '',
            phone: '',
            gender: '男',
            birthday: '',
            email: '',
            address: '',
            note: ''
        });
        
        // 积分表单
        const pointsForm = Vue.reactive({
            memberId: null,
            type: 'add',
            points: 100,
            reason: ''
        });
        
        // 表单验证规则
        const memberRules = {
            name: [
                { required: true, message: '请输入会员姓名', trigger: 'blur' }
            ],
            phone: [
                { required: true, message: '请输入手机号码', trigger: 'blur' },
                { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
            ],
            gender: [
                { required: true, message: '请选择性别', trigger: 'change' }
            ],
            email: [
                { type: 'email', message: '请输入正确的电子邮箱地址', trigger: 'blur' }
            ]
        };
        
        // 积分表单验证规则
        const pointsRules = {
            type: [
                { required: true, message: '请选择操作类型', trigger: 'change' }
            ],
            points: [
                { required: true, message: '请输入积分数量', trigger: 'blur' }
            ],
            reason: [
                { required: true, message: '请输入积分变动原因', trigger: 'blur' }
            ]
        };
        
        // 方法
        const loadMembers = async () => {
            loading.value = true;
            try {
                const memberService = new MemberService();
                const params = {
                    page: currentPage.value,
                    limit: pageSize.value,
                    ...searchForm
                };
                const response = await memberService.getMembers(params);
                members.value = response.items || [];
                total.value = response.total || 0;
            } catch (error) {
                console.error('加载会员列表失败:', error);
                ElementPlus.ElMessage.error('加载会员列表失败');
            } finally {
                loading.value = false;
            }
        };
        
        const handleSearch = () => {
            currentPage.value = 1;
            loadMembers();
        };
        
        const resetSearch = () => {
            Object.keys(searchForm).forEach(key => {
                searchForm[key] = '';
            });
            currentPage.value = 1;
            loadMembers();
        };
        
        const handleSizeChange = (val) => {
            pageSize.value = val;
            loadMembers();
        };
        
        const handleCurrentChange = (val) => {
            currentPage.value = val;
            loadMembers();
        };
        
        const handleAdd = () => {
            dialogType.value = 'add';
            resetForm();
            dialogVisible.value = true;
        };
        
        const handleEdit = (row) => {
            dialogType.value = 'edit';
            resetForm();
            Object.keys(memberForm).forEach(key => {
                if (key in row) {
                    memberForm[key] = row[key];
                }
            });
            dialogVisible.value = true;
        };
        
        const handlePoints = (row) => {
            selectedMember.value = row;
            resetPointsForm();
            pointsForm.memberId = row.id;
            pointsDialogVisible.value = true;
        };
        
        const handleDelete = (row) => {
            ElementPlus.ElMessageBox.confirm(
                `确定要删除会员 "${row.name}" 吗？`,
                '警告',
                {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }
            ).then(async () => {
                try {
                    const memberService = new MemberService();
                    await memberService.deleteMember(row.id);
                    ElementPlus.ElMessage.success('删除成功');
                    loadMembers();
                } catch (error) {
                    console.error('删除会员失败:', error);
                    ElementPlus.ElMessage.error('删除会员失败');
                }
            }).catch(() => {
                // 取消删除
            });
        };
        
        const resetForm = () => {
            Object.keys(memberForm).forEach(key => {
                if (key === 'id') {
                    memberForm[key] = null;
                } else if (key === 'gender') {
                    memberForm[key] = '男';
                } else {
                    memberForm[key] = '';
                }
            });
            if (memberFormRef.value) {
                memberFormRef.value.resetFields();
            }
        };
        
        const resetPointsForm = () => {
            pointsForm.memberId = null;
            pointsForm.type = 'add';
            pointsForm.points = 100;
            pointsForm.reason = '';
            if (pointsFormRef.value) {
                pointsFormRef.value.resetFields();
            }
        };
        
        const submitForm = () => {
            memberFormRef.value.validate(async (valid) => {
                if (valid) {
                    try {
                        const memberService = new MemberService();
                        if (dialogType.value === 'add') {
                            await memberService.createMember(memberForm);
                            ElementPlus.ElMessage.success('添加成功');
                        } else {
                            await memberService.updateMember(memberForm.id, memberForm);
                            ElementPlus.ElMessage.success('更新成功');
                        }
                        dialogVisible.value = false;
                        loadMembers();
                    } catch (error) {
                        console.error('保存会员失败:', error);
                        ElementPlus.ElMessage.error('保存会员失败');
                    }
                } else {
                    return false;
                }
            });
        };
        
        const submitPoints = () => {
            pointsFormRef.value.validate(async (valid) => {
                if (valid) {
                    try {
                        const memberService = new MemberService();
                        const pointsData = {
                            points: pointsForm.points,
                            reason: pointsForm.reason
                        };
                        
                        if (pointsForm.type === 'add') {
                            await memberService.addPoints(pointsForm.memberId, pointsData);
                            ElementPlus.ElMessage.success('积分增加成功');
                        } else {
                            await memberService.deductPoints(pointsForm.memberId, pointsData);
                            ElementPlus.ElMessage.success('积分扣减成功');
                        }
                        
                        pointsDialogVisible.value = false;
                        loadMembers();
                    } catch (error) {
                        console.error('积分操作失败:', error);
                        ElementPlus.ElMessage.error('积分操作失败');
                    }
                } else {
                    return false;
                }
            });
        };
        
        const getMemberLevelType = (level) => {
            switch (level) {
                case '普通会员':
                    return '';
                case '银卡会员':
                    return 'info';
                case '金卡会员':
                    return 'warning';
                case '钻石会员':
                    return 'success';
                default:
                    return '';
            }
        };
        
        const formatDate = (dateString) => {
            if (!dateString) return '';
            const date = new Date(dateString);
            return date.toLocaleDateString('zh-CN');
        };
        
        // 生命周期钩子
        Vue.onMounted(() => {
            loadMembers();
        });
        
        return {
            loading,
            members,
            currentPage,
            pageSize,
            total,
            dialogVisible,
            dialogType,
            memberFormRef,
            pointsDialogVisible,
            pointsFormRef,
            selectedMember,
            searchForm,
            memberForm,
            pointsForm,
            memberRules,
            pointsRules,
            handleSearch,
            resetSearch,
            handleSizeChange,
            handleCurrentChange,
            handleAdd,
            handleEdit,
            handlePoints,
            handleDelete,
            submitForm,
            submitPoints,
            getMemberLevelType,
            formatDate
        };
    }
};
