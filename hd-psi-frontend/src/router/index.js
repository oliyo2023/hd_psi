import { createRouter, createWebHistory } from 'vue-router'

// 导入视图组件
import Dashboard from '../views/Dashboard.vue'
import Login from '../views/Login.vue'
import ProductList from '../views/ProductList.vue'
import InventoryList from '../views/InventoryList.vue'
import MemberList from '../views/MemberList.vue'
import PurchaseList from '../views/PurchaseList.vue'
import PurchaseDetail from '../views/PurchaseDetail.vue'
import PurchaseCreate from '../views/PurchaseCreate.vue'
import SupplierList from '../views/SupplierList.vue'
import SupplierDetail from '../views/SupplierDetail.vue'
import NotFound from '../views/NotFound.vue'

// 路由配置
const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/products',
    name: 'ProductList',
    component: ProductList,
    meta: { requiresAuth: true }
  },
  {
    path: '/inventory',
    name: 'InventoryList',
    component: InventoryList,
    meta: { requiresAuth: true }
  },
  {
    path: '/members',
    name: 'MemberList',
    component: MemberList,
    meta: { requiresAuth: true }
  },
  {
    path: '/purchases',
    name: 'PurchaseList',
    component: PurchaseList,
    meta: { requiresAuth: true }
  },
  {
    path: '/purchases/create',
    name: 'PurchaseCreate',
    component: PurchaseCreate,
    meta: { requiresAuth: true }
  },
  {
    path: '/purchases/:id',
    name: 'PurchaseDetail',
    component: PurchaseDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/suppliers',
    name: 'SupplierList',
    component: SupplierList,
    meta: { requiresAuth: true }
  },
  {
    path: '/suppliers/:id',
    name: 'SupplierDetail',
    component: SupplierDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
