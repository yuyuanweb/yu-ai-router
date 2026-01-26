import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import UserLoginPage from '@/pages/user/UserLoginPage.vue'
import UserRegisterPage from '@/pages/user/UserRegisterPage.vue'
import UserManagePage from '@/pages/admin/UserManagePage.vue'
import ModelManagePage from '@/pages/admin/ModelManagePage.vue'
import ProviderManagePage from '@/pages/admin/ProviderManagePage.vue'
import BlacklistManagePage from '@/pages/admin/BlacklistManagePage.vue'
import PluginManagePage from '@/pages/admin/PluginManagePage.vue'
import ChatPage from '@/pages/ChatPage.vue'
import ApiKeyPage from '@/pages/user/ApiKeyPage.vue'
import ProfilePage from '@/pages/user/ProfilePage.vue'
import CallHistoryPage from '@/pages/user/CallHistoryPage.vue'
import RechargeSuccessPage from '@/pages/user/RechargeSuccessPage.vue'
import RechargeCancelPage from '@/pages/user/RechargeCancelPage.vue'
import ImageGenerationPage from '@/pages/ImageGenerationPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: '主页',
      component: ProfilePage,
    },
    {
      path: '/chat',
      name: 'AI 对话',
      component: ChatPage,
    },
    {
      path: '/image',
      name: 'AI 绘图',
      component: ImageGenerationPage,
    },
    {
      path: '/user/apikey',
      name: 'API Key 管理',
      component: ApiKeyPage,
    },
    {
      path: '/user/profile',
      name: '个人中心',
      component: ProfilePage,
    },
    {
      path: '/user/history',
      name: '调用历史',
      component: CallHistoryPage,
    },
    {
      path: '/recharge/success',
      name: '充值成功',
      component: RechargeSuccessPage,
    },
    {
      path: '/recharge/cancel',
      name: '充值取消',
      component: RechargeCancelPage,
    },
    {
      path: '/user/login',
      name: '用户登录',
      component: UserLoginPage,
    },
    {
      path: '/user/register',
      name: '用户注册',
      component: UserRegisterPage,
    },
    {
      path: '/admin/userManage',
      name: '用户管理',
      component: UserManagePage,
    },
    {
      path: '/admin/modelManage',
      name: '模型管理',
      component: ModelManagePage,
    },
    {
      path: '/admin/providerManage',
      name: '提供者管理',
      component: ProviderManagePage,
    },
    {
      path: '/admin/blacklistManage',
      name: '黑名单管理',
      component: BlacklistManagePage,
    },
    {
      path: '/admin/pluginManage',
      name: '插件管理',
      component: PluginManagePage,
    },
  ],
})

export default router
