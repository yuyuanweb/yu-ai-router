<template>
  <div class="chat-page">
    <a-card title="AI 对话" style="margin-bottom: 16px">
      <template #extra>
        <a-space>
          <a-select v-model:value="selectedModel" style="width: 200px">
            <a-select-option value="qwen-plus">通义千问-Plus</a-select-option>
          </a-select>
          <a-button @click="clearChat">清空对话</a-button>
        </a-space>
      </template>

      <!-- API Key 选择 -->
      <a-alert
        v-if="!selectedApiKeyId"
        message="请先选择一个 API Key"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
      />
      <a-select
        v-model:value="selectedApiKeyId"
        placeholder="选择 API Key"
        style="width: 100%; margin-bottom: 16px"
        @change="loadTokenStats"
      >
        <a-select-option v-for="key in apiKeys" :key="key.id" :value="key.id">
          {{ key.keyName || '未命名' }} ({{ key.keyValue }})
        </a-select-option>
      </a-select>

      <!-- Token 统计 -->
      <a-statistic
        title="已消耗 Token 总数"
        :value="totalTokens"
        style="margin-bottom: 16px"
      />

      <!-- 消息列表 -->
      <div class="message-list" ref="messageListRef">
        <div
          v-for="(msg, index) in messages"
          :key="index"
          :class="['message-item', msg.role]"
        >
          <div class="message-role">{{ getRoleLabel(msg.role || '') }}</div>
          <div class="message-content" v-html="formatContent(msg.content || '')"></div>
        </div>
        <!-- 加载中 -->
        <div v-if="loading" class="message-item assistant">
          <div class="message-role">AI</div>
          <div class="message-content">
            <a-spin /> <span v-html="formatContent(streamingContent || '')"></span>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-area">
        <a-textarea
          v-model:value="userInput"
          placeholder="请输入消息..."
          :auto-size="{ minRows: 3, maxRows: 6 }"
          :disabled="loading || !selectedApiKeyId"
          @pressEnter="handleSend"
        />
        <a-button
          type="primary"
          :loading="loading"
          :disabled="!userInput.trim() || !selectedApiKeyId"
          @click="handleSend"
          style="margin-top: 8px"
        >
          发送
        </a-button>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { message } from 'ant-design-vue'
import { listMyApiKeys } from '@/api/apiKeyController'
import { getMyTokenStats } from '@/api/statsController'

const selectedModel = ref('qwen-plus')
const selectedApiKeyId = ref<number>()
const apiKeys = ref<API.ApiKeyVO[]>([])
const messages = ref<API.ChatMessage[]>([])
const userInput = ref('')
const loading = ref(false)
const streamingContent = ref('')
const totalTokens = ref(0)
const messageListRef = ref()

// 加载 API Keys
const loadApiKeys = async () => {
  try {
    const res = await listMyApiKeys({
      pageNum: 1,
      pageSize: 10,
    })
    if (res.data.code === 0 && res.data.data && res.data.data.records) {
      apiKeys.value = res.data.data.records
      if (res.data.data.records.length > 0) {
        selectedApiKeyId.value = res.data.data.records[0].id
        await loadTokenStats()
      }
    }
  } catch (error: any) {
    message.error('加载 API Keys 失败：' + error.message)
  }
}

// 加载 Token 统计
const loadTokenStats = async () => {
  try {
    const res = await getMyTokenStats()
    if (res.data.code === 0 && res.data.data) {
      totalTokens.value = res.data.data.totalTokens || 0
    }
  } catch (error: any) {
    message.error('加载统计数据失败：' + error.message)
  }
}

// 发送消息
const handleSend = async () => {
  if (!userInput.value.trim() || !selectedApiKeyId.value) return

  const userMessage: API.ChatMessage = {
    role: 'user',
    content: userInput.value.trim(),
  }

  messages.value.push(userMessage)
  // 先不清空输入，等请求成功后再清空
  loading.value = true
  streamingContent.value = ''

  // 滚动到底部
  await nextTick()
  scrollToBottom()

  try {
    const chatRequest: API.ChatRequest = {
      model: selectedModel.value,
      messages: messages.value,
      stream: true,
    }

    let assistantContent = ''

    // 使用 fetch 调用内部聊天接口（传递 API Key ID）
    const response = await fetch(`/api/internal/chat/completions?apiKeyId=${selectedApiKeyId.value}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include', // 携带 session cookie
      body: JSON.stringify(chatRequest),
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const reader = response.body?.getReader()
    const decoder = new TextDecoder()

    if (!reader) {
      throw new Error('无法获取响应流')
    }

    let buffer = '' // 用于存储未处理完的数据

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      // 解码当前块，追加到缓冲区
      buffer += decoder.decode(value, { stream: true })

      // SSE 格式：每个消息以 \n\n 分隔
      const parts = buffer.split('\n\n')
      // 保留最后一个可能不完整的部分
      buffer = parts.pop() || ''

      // 处理完整的消息
      for (const part of parts) {
        const lines = part.split('\n')
        for (const line of lines) {
          if (line.startsWith('data:')) {
            // 移除 'data:' 前缀，如果有空格也一起移除
            const content = line.substring(5).trim()
            if (content && content !== '[DONE]') {
              assistantContent += content
              streamingContent.value = assistantContent
              scrollToBottom()
            }
          }
        }
      }
    }

    // 处理缓冲区剩余数据
    if (buffer) {
      const lines = buffer.split('\n')
      for (const line of lines) {
        if (line.startsWith('data:')) {
          // 移除 'data:' 前缀，如果有空格也一起移除
          const content = line.substring(5).trim()
          if (content && content !== '[DONE]') {
            assistantContent += content
            streamingContent.value = assistantContent
          }
        }
      }
    }

    // 完成
    messages.value.push({
      role: 'assistant',
      content: assistantContent,
    })
    streamingContent.value = ''
    loading.value = false
    // 清空输入框（请求成功后）
    userInput.value = ''
    // 刷新 Token 统计
    await loadTokenStats()
  } catch (error: any) {
    message.error('对话失败：' + error.message)
    loading.value = false
    streamingContent.value = ''
    // 失败时不清空输入框，方便用户重新发送
  }
}

// 清空对话
const clearChat = () => {
  messages.value = []
  streamingContent.value = ''
}

// 获取角色标签
const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    user: '我',
    assistant: 'AI',
    system: '系统',
  }
  return labels[role] || role
}

// 格式化内容，处理换行符
const formatContent = (content: string) => {
  if (!content) return ''
  // 先将后端转义的 \\n 还原为真正的换行符
  const unescaped = content.replace(/\\n/g, '\n')
  // 转义 HTML 特殊字符，防止 XSS 攻击
  const escaped = unescaped
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;')
  // 将换行符转换为 <br> 标签
  return escaped.replace(/\n/g, '<br>')
}

// 滚动到底部
const scrollToBottom = () => {
  if (messageListRef.value) {
    messageListRef.value.scrollTop = messageListRef.value.scrollHeight
  }
}

onMounted(() => {
  loadApiKeys()
})
</script>

<style scoped>
.chat-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px;
}

.message-list {
  max-height: 500px;
  overflow-y: auto;
  margin-bottom: 16px;
  padding: 16px;
  background: #f5f5f5;
  border-radius: 8px;
}

.message-item {
  margin-bottom: 16px;
  padding: 12px;
  border-radius: 8px;
  background: white;
}

.message-item.user {
  margin-left: 20%;
  background: #e6f7ff;
}

.message-item.assistant {
  margin-right: 20%;
  background: #f6ffed;
}

.message-role {
  font-weight: bold;
  margin-bottom: 8px;
  color: #1890ff;
}

.message-item.assistant .message-role {
  color: #52c41a;
}

.message-content {
  white-space: pre-wrap;
  word-break: break-word;
}

.input-area {
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}
</style>
