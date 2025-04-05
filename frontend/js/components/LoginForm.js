// 登录表单组件
const LoginForm = {
    template: `
        <div class="login-container">
            <el-form
                ref="loginFormRef"
                :model="loginForm"
                :rules="loginRules"
                class="login-form"
                @submit.prevent="handleLogin"
            >
                <h2 class="login-title">服装进销存系统</h2>
                
                <el-form-item prop="username">
                    <el-input
                        v-model="loginForm.username"
                        placeholder="用户名"
                        prefix-icon="el-icon-user"
                    ></el-input>
                </el-form-item>
                
                <el-form-item prop="password">
                    <el-input
                        v-model="loginForm.password"
                        type="password"
                        placeholder="密码"
                        prefix-icon="el-icon-lock"
                        show-password
                    ></el-input>
                </el-form-item>
                
                <el-form-item>
                    <el-button
                        type="primary"
                        :loading="loading"
                        style="width: 100%;"
                        @click="handleLogin"
                    >
                        登录
                    </el-button>
                </el-form-item>
            </el-form>
        </div>
    `,
    setup() {
        const router = VueRouter.useRouter();
        const loginFormRef = Vue.ref(null);
        const loading = Vue.ref(false);
        
        // 登录表单数据
        const loginForm = Vue.reactive({
            username: '',
            password: ''
        });
        
        // 表单验证规则
        const loginRules = {
            username: [
                { required: true, message: '请输入用户名', trigger: 'blur' }
            ],
            password: [
                { required: true, message: '请输入密码', trigger: 'blur' },
                { min: 6, message: '密码长度不能少于6个字符', trigger: 'blur' }
            ]
        };
        
        // 处理登录
        const handleLogin = () => {
            loginFormRef.value.validate(async (valid) => {
                if (valid) {
                    loading.value = true;
                    try {
                        const authService = new AuthService();
                        await authService.login(loginForm.username, loginForm.password);
                        
                        // 登录成功，跳转到首页
                        router.push('/dashboard');
                        
                        // 显示成功消息
                        ElementPlus.ElMessage.success('登录成功');
                    } catch (error) {
                        console.error('登录失败:', error);
                    } finally {
                        loading.value = false;
                    }
                } else {
                    return false;
                }
            });
        };
        
        return {
            loginFormRef,
            loginForm,
            loginRules,
            loading,
            handleLogin
        };
    }
};
