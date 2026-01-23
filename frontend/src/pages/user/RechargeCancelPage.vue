<template>
  <div id="rechargeCancelPage">
    <a-result
      status="warning"
      title="支付已取消"
      sub-title="您取消了本次充值，如需充值请重新发起"
    >
      <template #extra>
        <a-space>
          <a-button type="primary" size="large" @click="goToProfile">
            返回个人中心
          </a-button>
          <a-button size="large" @click="retryRecharge">
            重新充值
          </a-button>
        </a-space>
      </template>

      <template #default>
        <a-alert
          message="温馨提示"
          description="如遇支付问题，请检查网络连接或联系管理员获取帮助"
          type="info"
          show-icon
          style="margin-top: 24px; max-width: 600px; margin-left: auto; margin-right: auto"
        />
      </template>
    </a-result>
  </div>
</template>

<script lang="ts" setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'

const router = useRouter()

// 跳转到个人中心
const goToProfile = () => {
  router.push('/user/profile')
}

// 重新充值
const retryRecharge = () => {
  router.push('/user/profile')
  message.info('请点击"立即充值"按钮重新发起充值')
}

onMounted(() => {
  message.warning('您已取消本次充值')
})
</script>

<style scoped>
#rechargeCancelPage {
  padding: 60px 24px;
  background: #f0f2f5;
  min-height: calc(100vh - 64px);
}
</style>
