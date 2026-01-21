import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import UserLoginPage from '@/pages/user/UserLoginPage.vue'
import UserRegisterPage from '@/pages/user/UserRegisterPage.vue'
import UserManagePage from '@/pages/admin/UserManagePage.vue'
import ModelManagePage from '@/pages/admin/ModelManagePage.vue'
import ProviderManagePage from '@/pages/admin/ProviderManagePage.vue'
import ChatPage from '@/pages/ChatPage.vue'
import ApiKeyPage from '@/pages/user/ApiKeyPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: '主页',
      component: HomePage,
    },
    {
      path: '/chat',
      name: 'AI 对话',
      component: ChatPage,
    },
    {
      path: '/user/apikey',
      name: 'API Key 管理',
      component: ApiKeyPage,
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
  ],
})

export default router
