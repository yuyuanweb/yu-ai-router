<template>
  <a-layout-header class="morandi-header">
    <div class="header-right">
      <div v-if="loginUserStore.loginUser.id" class="user-profile">
        <a-dropdown placement="bottomRight">
          <div class="user-trigger">
            <a-avatar :src="loginUserStore.loginUser.userAvatar" class="user-avatar" />
            <span class="user-name">{{ loginUserStore.loginUser.userName ?? '无名' }}</span>
          </div>
          <template #overlay>
            <a-menu class="morandi-dropdown">
              <a-menu-item key="profile" @click="router.push('/user/profile')">
                <UserOutlined /> 个人中心
              </a-menu-item>
              <a-menu-divider />
              <a-menu-item key="logout" @click="doLogout">
                <LogoutOutlined /> 退出登录
              </a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
      <div v-else>
        <a-button type="primary" shape="round" href="/user/login">登录 / 注册</a-button>
      </div>
    </div>
  </a-layout-header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { useLoginUserStore } from '@/stores/loginUser.ts'
import { userLogout } from '@/api/userController.ts'
import { LogoutOutlined, UserOutlined } from '@ant-design/icons-vue'

const loginUserStore = useLoginUserStore()
const router = useRouter()

const doLogout = async () => {
  const res = await userLogout()
  if (res.data.code === 0) {
    loginUserStore.setLoginUser({ userName: '未登录' })
    message.success('退出登录成功')
    await router.push('/user/login')
  } else {
    message.error('退出登录失败，' + res.data.message)
  }
}
</script>

<style scoped>
.morandi-header {
  background: var(--bg-page) !important; /* 与背景融为一体 */
  height: 72px;
  padding: 0 32px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  border-bottom: 1px solid transparent;
}

.user-trigger {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 30px;
  transition: all 0.3s ease;
}

.user-trigger:hover {
  background: rgba(255, 255, 255, 0.5);
}

.user-avatar {
  border: 2px solid #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-main);
}

.morandi-dropdown {
  background: var(--bg-card);
  border-radius: 12px;
  box-shadow: var(--shadow-soft);
  border: 1px solid var(--border-color);
  padding: 8px;
}

:deep(.ant-dropdown-menu-item) {
  border-radius: 8px;
  transition: all 0.3s;
}
:deep(.ant-dropdown-menu-item:hover) {
  background-color: rgba(184, 169, 184, 0.1) !important;
  color: var(--color-primary) !important;
}
</style>
