<template>
  <div id="rechargeSuccessPage">
    <a-result
      status="success"
      title="充值成功！"
      :sub-title="subTitle"
    >
      <template #extra>
        <a-space>
          <a-button type="primary" size="large" @click="goToProfile">
            查看余额
          </a-button>
          <a-button size="large" @click="goToChatPage">
            开始对话
          </a-button>
        </a-space>
      </template>

      <template #default>
        <a-descriptions bordered :column="1" style="margin-top: 24px; max-width: 600px; margin-left: auto; margin-right: auto">
          <a-descriptions-item label="充值金额" v-if="rechargeInfo.amount">
            <span style="color: #52c41a; font-size: 18px; font-weight: bold">
              ¥{{ rechargeInfo.amount?.toFixed(2) }}
            </span>
          </a-descriptions-item>
          <a-descriptions-item label="支付方式" v-if="rechargeInfo.paymentMethod">
            <a-tag color="blue">{{ rechargeInfo.paymentMethod }}</a-tag>
          </a-descriptions-item>
          <a-descriptions-item label="支付ID" v-if="rechargeInfo.sessionId">
            <a-typography-text copyable>{{ rechargeInfo.sessionId }}</a-typography-text>
          </a-descriptions-item>
          <a-descriptions-item label="充值时间" v-if="rechargeInfo.createTime">
            {{ rechargeInfo.createTime }}
          </a-descriptions-item>
        </a-descriptions>
      </template>
    </a-result>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()

const rechargeInfo = ref<any>({})
const loading = ref(false)

const subTitle = computed(() => {
  if (rechargeInfo.value.amount) {
    return `¥${rechargeInfo.value.amount.toFixed(2)} 已成功充值到您的账户，您可以继续使用 AI 对话服务`
  }
  return '您的充值已成功处理，余额已更新'
})

// 跳转到个人中心
const goToProfile = () => {
  router.push('/user/profile')
}

// 跳转到对话页面
const goToChatPage = () => {
  router.push('/chat')
}

// 处理支付成功回调
const handlePaymentSuccess = async () => {
  const sessionId = route.query.session_id as string
  
  if (!sessionId) {
    message.warning('未找到支付信息')
    setTimeout(() => {
      router.push('/user/profile')
    }, 2000)
    return
  }

  loading.value = true
  try {
    // 调用后端处理支付成功
    const response = await fetch(`/api/recharge/stripe/success?session_id=${sessionId}`, {
      credentials: 'include',
    })
    
    const result = await response.json()
    
    if (result.code === 0) {
      // 从session_id提取信息（仅用于显示）
      rechargeInfo.value = {
        sessionId: sessionId,
        paymentMethod: 'Stripe',
        createTime: dayjs().format('YYYY-MM-DD HH:mm:ss'),
        // 金额信息会从实际充值记录中获取，这里先不显示
      }
      message.success('充值成功！余额已更新')
    } else {
      message.error('充值处理失败：' + result.message)
      setTimeout(() => {
        router.push('/user/profile')
      }, 2000)
    }
  } catch (error: any) {
    message.error('处理支付回调失败')
    setTimeout(() => {
      router.push('/user/profile')
    }, 2000)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  handlePaymentSuccess()
})
</script>

<style scoped>
#rechargeSuccessPage {
  padding: 60px 24px;
  background: #f0f2f5;
  min-height: calc(100vh - 64px);
}
</style>
