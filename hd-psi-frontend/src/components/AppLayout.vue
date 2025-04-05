<template>
  <n-config-provider :locale="zhCN" :theme="theme">
    <n-loading-bar-provider>
      <n-dialog-provider>
        <n-notification-provider>
          <n-message-provider>
            <div class="app-container">
              <template v-if="isLoggedIn">
                <n-layout has-sider>
                  <!-- 侧边栏 -->
                  <n-layout-sider
                    :collapsed="collapsed"
                    :collapsed-width="64"
                    :width="240"
                    show-trigger
                    collapse-mode="width"
                    @collapse="collapsed = true"
                    @expand="collapsed = false"
                    bordered
                    class="sidebar"
                  >
                    <div class="logo">
                      <h2 v-if="!collapsed">服装进销存系统</h2>
                      <h2 v-else>HD</h2>
                    </div>
                    <n-menu
                      :collapsed="collapsed"
                      :collapsed-width="64"
                      :collapsed-icon-size="22"
                      :options="menuOptions"
                      :render-label="renderMenuLabel"
                      :render-icon="renderMenuIcon"
                      :value="activeKey"
                      @update:value="handleMenuUpdate"
                    />
                  </n-layout-sider>
                  
                  <!-- 主内容区 -->
                  <n-layout>
                    <!-- 顶部导航栏 -->
                    <n-layout-header bordered class="header">
                      <div class="header-content">
                        <div class="left">
                          <n-button quaternary circle @click="collapsed = !collapsed">
                            <template #icon>
                              <n-icon size="20">
                                <MenuOutline v-if="collapsed" />
                                <CloseOutline v-else />
                              </n-icon>
                            </template>
                          </n-button>
                        </div>
                        <div class="right">
                          <n-dropdown :options="userOptions" @select="handleUserAction">
                            <n-button quaternary>
                              {{ currentUser.name || currentUser.username }}
                              <n-icon size="14" style="margin-left: 5px">
                                <ChevronDownOutline />
                              </n-icon>
                            </n-button>
                          </n-dropdown>
                        </div>
                      </div>
                    </n-layout-header>
                    
                    <!-- 内容区 -->
                    <n-layout-content class="content">
                      <router-view />
                    </n-layout-content>
                    
                    <!-- 页脚 -->
                    <n-layout-footer bordered class="footer">
                      <p>© {{ new Date().getFullYear() }} 服装进销存系统 (HD-PSI)</p>
                    </n-layout-footer>
                  </n-layout>
                </n-layout>
              </template>
              <template v-else>
                <router-view />
              </template>
            </div>
          </n-message-provider>
        </n-notification-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { 
  NConfigProvider, NLayout, NLayoutSider, NLayoutHeader, 
  NLayoutContent, NLayoutFooter, NMenu, NButton, NIcon,
  NDropdown, NMessageProvider, NNotificationProvider,
  NDialogProvider, NLoadingBarProvider
} from 'naive-ui'
import { zhCN, darkTheme } from 'naive-ui'
import { 
  HomeOutline, CartOutline, PeopleOutline, 
  PersonOutline, LogOutOutline, SettingsOutline,
  MenuOutline, CloseOutline, ChevronDownOutline,
  BagHandleOutline, StorefrontOutline, CubeOutline
} from '@vicons/ionicons5'
import auth from '../services/auth'

// 响应式状态
const router = useRouter()
const route = useRoute()
const collapsed = ref(false)
const isLoggedIn = ref(false)
const currentUser = ref({})
const theme = ref(null) // 默认使用亮色主题

// 计算属性
const activeKey = computed(() => {
  const path = route.path
  if (path.startsWith('/dashboard')) return 'dashboard'
  if (path.startsWith('/products')) return 'products'
  if (path.startsWith('/inventory')) return 'inventory'
  if (path.startsWith('/members')) return 'members'
  if (path.startsWith('/purchases')) return 'purchases'
  if (path.startsWith('/suppliers')) return 'suppliers'
  return ''
})

// 菜单配置
const menuOptions = [
  {
    label: '仪表盘',
    key: 'dashboard',
    icon: HomeOutline,
    path: '/dashboard'
  },
  {
    label: '商品管理',
    key: 'products',
    icon: CubeOutline,
    path: '/products'
  },
  {
    label: '库存管理',
    key: 'inventory',
    icon: StorefrontOutline,
    path: '/inventory'
  },
  {
    label: '采购管理',
    key: 'purchases',
    icon: BagHandleOutline,
    children: [
      {
        label: '采购单',
        key: 'purchase-orders',
        path: '/purchases'
      },
      {
        label: '供应商',
        key: 'suppliers',
        path: '/suppliers'
      }
    ]
  },
  {
    label: '会员管理',
    key: 'members',
    icon: PeopleOutline,
    path: '/members'
  },
  {
    label: '销售管理',
    key: 'sales',
    icon: CartOutline,
    path: '/sales'
  }
]

// 用户下拉菜单选项
const userOptions = [
  {
    label: '个人信息',
    key: 'profile',
    icon: () => h(NIcon, null, { default: () => h(PersonOutline) })
  },
  {
    label: '系统设置',
    key: 'settings',
    icon: () => h(NIcon, null, { default: () => h(SettingsOutline) })
  },
  {
    type: 'divider',
    key: 'd1'
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h(NIcon, null, { default: () => h(LogOutOutline) })
  }
]

// 方法
const renderMenuLabel = (option) => {
  return option.label
}

const renderMenuIcon = (option) => {
  return option.icon ? h(NIcon, null, { default: () => h(option.icon) }) : null
}

const handleMenuUpdate = (key) => {
  const findPath = (options, key) => {
    for (const option of options) {
      if (option.key === key) return option.path
      if (option.children) {
        const path = findPath(option.children, key)
        if (path) return path
      }
    }
    return null
  }
  
  const path = findPath(menuOptions, key)
  if (path) router.push(path)
}

const handleUserAction = (key) => {
  if (key === 'logout') {
    auth.logout()
  } else if (key === 'profile') {
    router.push('/profile')
  } else if (key === 'settings') {
    router.push('/settings')
  }
}

// 生命周期钩子
onMounted(() => {
  isLoggedIn.value = auth.isAuthenticated()
  currentUser.value = auth.getCurrentUser() || {}
})
</script>

<style scoped>
.app-container {
  height: 100vh;
  width: 100%;
}

.sidebar {
  background-color: #001529;
}

.logo {
  height: 64px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  overflow: hidden;
}

.header {
  height: 64px;
  padding: 0 16px;
}

.header-content {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.content {
  padding: 24px;
  min-height: calc(100vh - 64px - 48px);
  background-color: #f0f2f5;
}

.footer {
  height: 48px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
}
</style>
