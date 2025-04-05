// 路由配置
const routes = [
    {
        path: '/',
        redirect: '/dashboard'
    },
    {
        path: '/login',
        component: LoginForm,
        meta: { requiresAuth: false }
    },
    {
        path: '/dashboard',
        component: Dashboard,
        meta: { requiresAuth: true }
    },
    {
        path: '/products',
        component: ProductList,
        meta: { requiresAuth: true }
    },
    {
        path: '/inventory',
        component: InventoryList,
        meta: { requiresAuth: true }
    },
    {
        path: '/members',
        component: MemberList,
        meta: { requiresAuth: true }
    }
];

// 创建路由实例
const router = VueRouter.createRouter({
    history: VueRouter.createWebHashHistory(),
    routes
});

// 路由守卫
router.beforeEach((to, from, next) => {
    const authService = new AuthService();
    const isAuthenticated = authService.isAuthenticated();
    
    if (to.meta.requiresAuth && !isAuthenticated) {
        // 需要认证但未登录，重定向到登录页
        next('/login');
    } else if (to.path === '/login' && isAuthenticated) {
        // 已登录但访问登录页，重定向到首页
        next('/dashboard');
    } else {
        // 正常导航
        next();
    }
});
