<template>
  <div id="userRegisterPage">
    <h2 class="title">用户注册</h2>
    <a-form :model="formState" name="basic" autocomplete="off" @finish="handleSubmit">
      <a-form-item name="userAccount" :rules="[{ required: true, message: '请输入账号' }]">
        <a-input v-model:value="formState.userAccount" placeholder="请输入账号" />
      </a-form-item>
      <a-form-item
        name="userPassword"
        :rules="[
          { required: true, message: '请输入密码' },
          { min: 8, message: '密码不能小于 8 位' },
        ]"
      >
        <a-input-password v-model:value="formState.userPassword" placeholder="请输入密码" />
      </a-form-item>
      <a-form-item
        name="checkPassword"
        :rules="[
          { required: true, message: '请确认密码' },
          { min: 8, message: '密码不能小于 8 位' },
          { validator: validateCheckPassword },
        ]"
      >
        <a-input-password v-model:value="formState.checkPassword" placeholder="请确认密码" />
      </a-form-item>
      <div class="tips">
        已有账号？
        <RouterLink to="/user/login">去登录</RouterLink>
      </div>
      <a-form-item>
        <a-button type="primary" html-type="submit" style="width: 100%">注册</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { userRegister } from '@/api/userController.ts'
import { message } from 'ant-design-vue'
import { reactive } from 'vue'

const router = useRouter()

const formState = reactive<API.UserRegisterRequest>({
  userAccount: '',
  userPassword: '',
  checkPassword: '',
})

/**
 * 验证确认密码
 * @param rule
 * @param value
 * @param callback
 */
const validateCheckPassword = (rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value && value !== formState.userPassword) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

/**
 * 提交表单
 * @param values
 */
const handleSubmit = async (values: API.UserRegisterRequest) => {
  const res = await userRegister(values)
  // 注册成功，跳转到登录页面
  if (res.data.code === 0) {
    message.success('注册成功')
    router.push({
      path: '/user/login',
      replace: true,
    })
  } else {
    message.error('注册失败，' + res.data.message)
  }
}
</script>

/* 替换原有 style */
<style scoped>
#userRegisterPage {
  background: var(--bg-card); /* 奶白卡片 */
  max-width: 480px; /* 注册页不需要太宽 */
  padding: 40px;
  margin: 60px auto;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-soft); /* 极淡的阴影 */
  border: 1px solid var(--border-color);
}

.title {
  text-align: center;
  margin-bottom: 32px;
  color: var(--text-main);
  font-weight: 600;
}

.tips {
  margin-bottom: 24px;
  color: var(--text-secondary);
  font-size: 13px;
  text-align: right;
}

.tips a {
  color: var(--color-primary);
  transition: opacity 0.3s;
}

.tips a:hover {
  opacity: 0.8;
}
</style>
