// 应用布局组件
const AppLayout = {
    template: `
        <div v-if="isLoggedIn" class="app-container">
            <el-container>
                <el-aside width="auto">
                    <div class="sidebar">
                        <div class="sidebar-logo">
                            <span v-if="!isCollapse">服装进销存系统</span>
                            <span v-else>HD-PSI</span>
                        </div>
                        <el-menu
                            :default-active="activeMenu"
                            class="el-menu-vertical"
                            :collapse="isCollapse"
                            background-color="#304156"
                            text-color="#bfcbd9"
                            active-text-color="#409EFF"
                            router
                        >
                            <el-menu-item index="/dashboard">
                                <el-icon><el-icon-odometer /></el-icon>
                                <template #title>仪表盘</template>
                            </el-menu-item>
                            
                            <el-sub-menu index="1">
                                <template #title>
                                    <el-icon><el-icon-goods /></el-icon>
                                    <span>商品管理</span>
                                </template>
                                <el-menu-item index="/products">商品列表</el-menu-item>
                            </el-sub-menu>
                            
                            <el-sub-menu index="2">
                                <template #title>
                                    <el-icon><el-icon-box /></el-icon>
                                    <span>库存管理</span>
                                </template>
                                <el-menu-item index="/inventory">库存列表</el-menu-item>
                                <el-menu-item index="/inventory-transactions">库存交易</el-menu-item>
                                <el-menu-item index="/inventory-alerts">库存预警</el-menu-item>
                                <el-menu-item index="/inventory-checks">库存盘点</el-menu-item>
                            </el-sub-menu>
                            
                            <el-sub-menu index="3">
                                <template #title>
                                    <el-icon><el-icon-user /></el-icon>
                                    <span>会员管理</span>
                                </template>
                                <el-menu-item index="/members">会员列表</el-menu-item>
                            </el-sub-menu>
                            
                            <el-sub-menu index="4">
                                <template #title>
                                    <el-icon><el-icon-shopping-cart /></el-icon>
                                    <span>销售管理</span>
                                </template>
                                <el-menu-item index="/sales">销售订单</el-menu-item>
                                <el-menu-item index="/returns">退换货管理</el-menu-item>
                            </el-sub-menu>
                            
                            <el-sub-menu index="5">
                                <template #title>
                                    <el-icon><el-icon-office-building /></el-icon>
                                    <span>店铺管理</span>
                                </template>
                                <el-menu-item index="/stores">店铺列表</el-menu-item>
                            </el-sub-menu>
                            
                            <el-sub-menu index="6">
                                <template #title>
                                    <el-icon><el-icon-setting /></el-icon>
                                    <span>系统设置</span>
                                </template>
                                <el-menu-item index="/profile">个人信息</el-menu-item>
                                <el-menu-item index="/users">用户管理</el-menu-item>
                            </el-sub-menu>
                        </el-menu>
                    </div>
                </el-aside>
                
                <el-container>
                    <el-header style="height: 60px;">
                        <div style="display: flex; justify-content: space-between; align-items: center; height: 100%;">
                            <div>
                                <el-icon 
                                    :size="24" 
                                    style="cursor: pointer; margin-right: 20px;"
                                    @click="toggleSidebar"
                                >
                                    <el-icon-fold v-if="!isCollapse" />
                                    <el-icon-expand v-else />
                                </el-icon>
                            </div>
                            
                            <div>
                                <el-dropdown trigger="click" @command="handleCommand">
                                    <span class="el-dropdown-link" style="cursor: pointer;">
                                        {{ currentUser.name }}
                                        <el-icon class="el-icon--right"><el-icon-arrow-down /></el-icon>
                                    </span>
                                    <template #dropdown>
                                        <el-dropdown-menu>
                                            <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                                            <el-dropdown-item command="changePassword">修改密码</el-dropdown-item>
                                            <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                                        </el-dropdown-menu>
                                    </template>
                                </el-dropdown>
                            </div>
                        </div>
                    </el-header>
                    
                    <el-main>
                        <router-view></router-view>
                    </el-main>
                    
                    <el-footer style="height: 30px; line-height: 30px; text-align: center; font-size: 12px; color: #909399;">
                        © 2023 服装进销存系统 (HD-PSI)
                    </el-footer>
                </el-container>
            </el-container>
        </div>
        <router-view v-else></router-view>
    `,
    setup() {
        const router = VueRouter.useRouter();
        const route = VueRouter.useRoute();
        
        // 响应式状态
        const isCollapse = Vue.ref(false);
        const isLoggedIn = Vue.ref(false);
        const currentUser = Vue.ref({});
        
        // 计算属性
        const activeMenu = Vue.computed(() => {
            return route.path;
        });
        
        // 方法
        const toggleSidebar = () => {
            isCollapse.value = !isCollapse.value;
        };
        
        const handleCommand = (command) => {
            if (command === 'logout') {
                const authService = new AuthService();
                authService.logout();
            } else if (command === 'profile') {
                router.push('/profile');
            } else if (command === 'changePassword') {
                router.push('/change-password');
            }
        };
        
        // 生命周期钩子
        Vue.onMounted(() => {
            const authService = new AuthService();
            isLoggedIn.value = authService.isAuthenticated();
            currentUser.value = authService.getCurrentUser() || {};
        });
        
        // 监听路由变化
        Vue.watch(
            () => route.path,
            () => {
                const authService = new AuthService();
                isLoggedIn.value = authService.isAuthenticated();
                currentUser.value = authService.getCurrentUser() || {};
            }
        );
        
        return {
            isCollapse,
            isLoggedIn,
            currentUser,
            activeMenu,
            toggleSidebar,
            handleCommand
        };
    }
};
