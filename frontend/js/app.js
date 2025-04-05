// 创建Vue应用
const app = Vue.createApp({
    setup() {
        const zhCn = ElementPlusLocaleZhCn;
        return {
            zhCn
        };
    }
});

// 注册组件
app.component('app-layout', AppLayout);
app.component('login-form', LoginForm);
app.component('dashboard', Dashboard);
app.component('product-list', ProductList);
app.component('inventory-list', InventoryList);
app.component('member-list', MemberList);

// 使用路由
app.use(router);

// 使用Element Plus
app.use(ElementPlus);

// 挂载应用
app.mount('#app');
