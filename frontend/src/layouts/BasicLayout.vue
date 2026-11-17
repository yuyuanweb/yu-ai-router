<template>
  <a-layout class="morandi-layout">
    <!-- 左侧侧边栏 -->
    <a-layout-sider class="morandi-sider" width="240" theme="light">
      <div class="logo-area">
        <img class="logo" src="@/assets/logo.png" alt="Logo" />
        <span class="site-title">摇光 AI 网关</span>
      </div>
      <a-menu
        v-model:selectedKeys="selectedKeys"
        mode="inline"
        :items="menuItems"
        class="morandi-menu"
        @click="handleMenuClick"
      />
    </a-layout-sider>

    <!-- 右侧主内容区 -->
    <a-layout class="main-layout">
      <GlobalHeader />
      <a-layout-content class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { type MenuProps } from 'ant-design-vue'
import { useLoginUserStore } from '@/stores/loginUser.ts'
import GlobalHeader from '@/components/GlobalHeader.vue'
import {
  HomeOutlined, CommentOutlined, PictureOutlined, KeyOutlined,
  HistoryOutlined, SafetyOutlined, AppstoreOutlined, CloudServerOutlined,
  StopOutlined, ApiOutlined, TeamOutlined
} from '@ant-design/icons-vue'

const loginUserStore = useLoginUserStore()
const router = useRouter()
const selectedKeys = ref<string[]>(['/'])

// 路由变化监听
router.afterEach((to) => {
  selectedKeys.value = [to.path]
})

// 菜单配置（合并了用户管理和后台管理，更清爽）
const originItems = [
  { key: '/', icon: () => h(HomeOutlined), label: '主页' },
  { key: '/chat', icon: () => h(CommentOutlined), label: 'AI 对话' },
  { key: '/image', icon: () => h(PictureOutlined), label: 'AI 绘图' },
  { type: 'divider' },
  { key: '/user/apikey', icon: () => h(KeyOutlined), label: 'API Key' },
  { key: '/user/mykeys', icon: () => h(SafetyOutlined), label: '我的密钥' },
  { key: '/user/history', icon: () => h(HistoryOutlined), label: '调用历史' },
  { type: 'divider' },
  { key: '/admin/userManage', icon: () => h(TeamOutlined), label: '用户管理' },
  { key: '/admin/modelManage', icon: () => h(AppstoreOutlined), label: '模型管理' },
  { key: '/admin/providerManage', icon: () => h(CloudServerOutlined), label: '提供者管理' },
  { key: '/admin/blacklistManage', icon: () => h(StopOutlined), label: '黑名单' },
  { key: '/admin/pluginManage', icon: () => h(ApiOutlined), label: '插件管理' },
]

// 过滤管理员菜单
const menuItems = computed<MenuProps['items']>(() => {
  return originItems.filter((menu) => {
    if (menu.key?.toString().startsWith('/admin')) {
      return loginUserStore.loginUser?.userRole === 'admin'
    }
    return true
  })
})

const handleMenuClick: MenuProps['onClick'] = (e) => {
  if (e.key.toString().startsWith('/')) {
    router.push(e.key.toString())
  }
}
</script>

<style scoped>
.morandi-layout {
  height: 100vh;
  overflow: hidden;
  background: var(--bg-page);
}

.morandi-sider {
  background: var(--bg-card) !important;
  border-right: 1px dashed var(--border-color); /* 细虚线分割，更柔和 */
  display: flex;
  flex-direction: column;
}

.logo-area {
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 24px;
}

.logo {
  width: 36px;
  height: 36px;
  border-radius: 10px;
}

.site-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-main);
  letter-spacing: 1px;
}

/* 侧边栏菜单深度美化 */
.morandi-menu {
  border-inline-end: none !important;
  background: transparent !important;
  padding: 12px;
}

:deep(.morandi-menu .ant-menu-item) {
  height: 48px;
  line-height: 48px;
  border-radius: 12px;
  margin-bottom: 4px;
  color: var(--text-secondary);
  font-weight: 500;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
}

:deep(.morandi-menu .ant-menu-item:hover) {
  background-color: rgba(184, 169, 184, 0.08) !important;
  color: var(--text-main) !important;
}

:deep(.morandi-menu .ant-menu-item-selected) {
  background-color: rgba(184, 169, 184, 0.15) !important;
  color: #8B7896 !important; /* 雾霾紫加深版 */
  font-weight: 600;
}

:deep(.morandi-menu .ant-menu-item-divider) {
  border-color: var(--border-color) !important;
  border-style: dashed !important;
  margin: 16px 0;
}

.main-layout {
  background: var(--bg-page);
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 24px;
  position: relative;
}

/* 页面切换动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.4s ease;
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(12px);
}
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}
</style>
