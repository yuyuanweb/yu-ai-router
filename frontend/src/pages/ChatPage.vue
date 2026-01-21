<template>
  <div id="chatPage">
    <!-- 对话内容包裹层 -->
    <div class="chat-wrapper">
      <!-- 消息列表区域（带滚动条） -->
      <div class="messages-container" ref="messagesContainer">
        <div v-if="messages.length === 0" class="empty-state">
          <h2>开始与 AI 对话</h2>
          <p>选择一个模型，然后输入您的问题</p>
        </div>

        <div v-for="(msg, index) in messages" :key="index" class="message-wrapper">
          <div class="message" :class="msg.role">
            <div class="message-header">
              <strong v-if="msg.role === 'user'">你</strong>
              <strong v-else>{{ getAIDisplayName() }}</strong>
            </div>
            <!-- AI 回答（可能包含思考内容） -->
            <div v-if="msg.role === 'assistant'" class="message-content">
              <div v-if="msg.thinking" class="thinking-section">
                <div class="thinking-toggle" @click="toggleThinking(index)">
                  <span class="thinking-title">已深度思考</span>
                  <DownOutlined :class="['thinking-arrow', { expanded: msg.thinkingExpanded }]" />
                </div>
                <div v-show="msg.thinkingExpanded" class="thinking-detail" v-html="renderMarkdown(msg.thinking)"></div>
              </div>
              <div class="answer-text" v-html="renderMarkdown(msg.answer || msg.content)"></div>
            </div>
            <!-- 用户消息 -->
            <div v-else class="message-content" v-html="renderMarkdown(msg.content)"></div>
          </div>
        </div>

        <!-- 流式响应中的消息 -->
        <div v-if="streamingContent || streamingThinking" class="message-wrapper">
          <div class="message assistant">
            <div class="message-header">
              <strong>{{ getAIDisplayName() }}</strong>
            </div>
            <div class="message-content">
              <!-- 思考内容（如果有） -->
              <div v-if="streamingThinking" class="thinking-section">
                <div class="thinking-toggle">
                  <span class="thinking-title">思考中...</span>
                </div>
                <div class="thinking-detail" v-html="renderMarkdown(streamingThinking)"></div>
              </div>
              <!-- 最终答案 -->
              <div v-if="streamingContent" class="answer-text" v-html="renderMarkdown(streamingContent)"></div>
            </div>
          </div>
        </div>

        <!-- 加载中提示 -->
        <div v-if="loading && !streamingContent" class="message-wrapper">
          <div class="message assistant">
            <div class="message-header">
              <strong>{{ getAIDisplayName() }}</strong>
            </div>
            <div class="message-content">
              <a-spin /> 思考中...
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 -->
      <div class="input-container">
        <!-- 模型选择和深度思考选项 -->
        <div class="options-bar">
          <!-- 路由策略选择 -->
          <a-dropdown :trigger="['click']">
            <a-button class="strategy-select-button" size="large">
              <ThunderboltOutlined />
              <span class="strategy-label">{{ getCurrentStrategyLabel() }}</span>
              <DownOutlined class="dropdown-icon" />
            </a-button>
            <template #overlay>
              <a-menu @click="handleStrategyChange">
                <a-menu-item
                  v-for="option in routingStrategyOptions"
                  :key="option.value"
                  :class="{ 'strategy-active': routingStrategy === option.value }"
                >
                  <div class="strategy-option">
                    <div class="strategy-option-header">
                      <span class="strategy-option-label">{{ option.label }}</span>
                      <a-tag v-if="routingStrategy === option.value" color="blue" size="small">当前</a-tag>
                    </div>
                    <div class="strategy-option-desc">{{ option.description }}</div>
                  </div>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>

          <!-- 模型选择按钮（仅固定模型策略时显示） -->
          <a-button v-if="routingStrategy === 'fixed'" @click="showModelSelector" size="large" class="model-select-button">
            <template v-if="selectedModel">
              <AppstoreOutlined />
              <span class="current-model">{{ selectedModel.modelName }}</span>
              <a-tag v-if="selectedModel.supportReasoning === 1" color="purple" size="small">思考</a-tag>
            </template>
            <template v-else>
              <AppstoreOutlined />
              选择模型
            </template>
            <DownOutlined class="dropdown-icon" />
          </a-button>

          <!-- 深度思考开关（仅固定模型策略且模型支持时显示） -->
          <div v-if="routingStrategy === 'fixed' && selectedModel?.supportReasoning === 1" class="reasoning-switch">
            <BulbOutlined class="reasoning-icon" />
            <span class="reasoning-label">深度思考</span>
            <a-switch v-model:checked="enableReasoning" />
          </div>
        </div>

        <div class="input-wrapper">
          <a-textarea
            v-model:value="inputMessage"
            :placeholder="canSendMessage ? '输入您的问题...' : '请先选择一个模型'"
            :auto-size="{ minRows: 2, maxRows: 6 }"
            :disabled="!canSendMessage || loading"
            @pressEnter="handleEnter"
            class="chat-input"
          />
          <div class="input-actions">
            <a-button
              type="primary"
              :loading="loading"
              :disabled="!canSendMessage || !inputMessage.trim()"
              @click="sendMessage"
              size="large"
            >
              发送
            </a-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 模型选择弹窗 -->
    <a-modal
      v-model:open="modelSelectorVisible"
      title="选择模型"
      :footer="null"
      width="900px"
      :body-style="{ maxHeight: '70vh', overflowY: 'auto' }"
    >
      <a-tabs v-model:activeKey="activeCategory" type="card">
        <a-tab-pane key="all" tab="全部模型">
          <div class="model-grid">
            <div
              v-for="model in allModels"
              :key="model.id"
              class="model-card"
              :class="{ active: selectedModel?.id === model.id }"
              @click="selectModelAndClose(model)"
            >
              <div class="model-header">
                <span class="model-name">{{ model.modelName }}</span>
                <a-tag v-if="model.supportReasoning === 1" color="purple" size="small">思考</a-tag>
              </div>
              <div class="model-info">
                <span class="model-provider">{{ model.providerDisplayName }}</span>
                <span class="model-price">¥{{ model.inputPrice }}/¥{{ model.outputPrice }}</span>
              </div>
              <div class="model-desc">{{ model.description }}</div>
            </div>
          </div>
        </a-tab-pane>

        <a-tab-pane key="fast" tab="快速模型">
          <div class="model-grid">
            <div
              v-for="model in fastModels"
              :key="model.id"
              class="model-card"
              :class="{ active: selectedModel?.id === model.id }"
              @click="selectModelAndClose(model)"
            >
              <div class="model-header">
                <span class="model-name">{{ model.modelName }}</span>
              </div>
              <div class="model-info">
                <span class="model-provider">{{ model.providerDisplayName }}</span>
                <span class="model-price">¥{{ model.inputPrice }}/¥{{ model.outputPrice }}</span>
              </div>
              <div class="model-desc">{{ model.description }}</div>
            </div>
          </div>
        </a-tab-pane>

        <a-tab-pane key="reasoning" tab="深度思考">
          <div class="model-grid">
            <div
              v-for="model in reasoningModels"
              :key="model.id"
              class="model-card"
              :class="{ active: selectedModel?.id === model.id }"
              @click="selectModelAndClose(model)"
            >
              <div class="model-header">
                <span class="model-name">{{ model.modelName }}</span>
                <a-tag color="purple" size="small">思考</a-tag>
              </div>
              <div class="model-info">
                <span class="model-provider">{{ model.providerDisplayName }}</span>
                <span class="model-price">¥{{ model.inputPrice }}/¥{{ model.outputPrice }}</span>
              </div>
              <div class="model-desc">{{ model.description }}</div>
            </div>
          </div>
        </a-tab-pane>

      </a-tabs>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { listActiveModels } from '@/api/modelController'
import { message } from 'ant-design-vue'
import { marked } from 'marked'
import { AppstoreOutlined, DownOutlined, BulbOutlined, ThunderboltOutlined } from '@ant-design/icons-vue'

// 路由策略选项
const routingStrategyOptions = [
  { value: 'auto', label: '自动路由', description: '综合成本、延迟、优先级智能选择' },
  { value: 'cost_first', label: '成本优先', description: '选择费用最低的模型' },
  { value: 'latency_first', label: '延迟优先', description: '选择响应最快的模型' },
  { value: 'fixed', label: '固定模型', description: '使用选定的指定模型' },
]

// 模型列表
const allModels = ref<API.ModelVO[]>([])
const selectedModel = ref<API.ModelVO | null>(null)
const activeCategory = ref('all')
const modelSelectorVisible = ref(false)

// 对话消息
interface Message {
  role: string
  content: string
  thinking?: string
  answer?: string
  thinkingExpanded?: boolean
}

const messages = ref<Message[]>([])
const inputMessage = ref('')
const loading = ref(false)
const messagesContainer = ref<HTMLElement>()
const streamingContent = ref<string>('')
const streamingThinking = ref<string>('')

// 是否启用深度思考
const enableReasoning = ref<boolean>(false)

// 路由策略
const routingStrategy = ref<'auto' | 'cost_first' | 'latency_first' | 'fixed'>('auto')

// 计算不同类别的模型
const fastModels = computed(() => {
  return allModels.value.filter(m => m.supportReasoning !== 1
  )
})

const reasoningModels = computed(() => {
  return allModels.value.filter(m => m.supportReasoning === 1)
})

// 判断是否可以发送消息
const canSendMessage = computed(() => {
  // 非固定模型策略时，只要有模型列表就可以发送（后端自动选择）
  if (routingStrategy.value !== 'fixed') {
    return allModels.value.length > 0
  }
  // 固定模型策略时，需要选择模型
  return selectedModel.value !== null
})

// 加载模型列表
const loadModels = async () => {
  try {
    const res = await listActiveModels()
    if (res.data.data) {
      allModels.value = res.data.data
      // 默认选择第一个模型
      if (allModels.value.length > 0) {
        selectedModel.value = allModels.value[0]
      }
    }
  } catch (err) {
    console.error('加载模型列表失败', err)
    message.error('加载模型列表失败')
  }
}

// 显示模型选择弹窗
const showModelSelector = () => {
  modelSelectorVisible.value = true
}

// 选择模型并关闭弹窗
const selectModelAndClose = (model: API.ModelVO) => {
  selectedModel.value = model
  // 如果模型不支持深度思考，关闭深度思考
  if (model.supportReasoning !== 1) {
    enableReasoning.value = false
  }
  modelSelectorVisible.value = false
}

// 获取当前策略的显示标签
const getCurrentStrategyLabel = () => {
  const option = routingStrategyOptions.find(opt => opt.value === routingStrategy.value)
  return option?.label || '自动路由'
}

// 获取 AI 显示名称
const getAIDisplayName = () => {
  if (routingStrategy.value === 'fixed' && selectedModel.value) {
    return selectedModel.value.modelName
  }
  // 非固定策略时显示策略名称
  const option = routingStrategyOptions.find(opt => opt.value === routingStrategy.value)
  return option?.label || 'AI'
}

// 处理路由策略变更
const handleStrategyChange = ({ key }: { key: string }) => {
  routingStrategy.value = key as 'auto' | 'cost_first' | 'latency_first' | 'fixed'
}

// 发送消息
const sendMessage = async () => {
  if (!canSendMessage.value || !inputMessage.value.trim()) {
    return
  }

  const userMessage = inputMessage.value.trim()

  // 添加用户消息到对话列表
  messages.value.push({
    role: 'user',
    content: userMessage,
  })

  loading.value = true
  streamingContent.value = ''
  streamingThinking.value = ''

  // 滚动到底部
  nextTick(() => scrollToBottom())

  try {
    // 构建请求
    const chatRequest: Record<string, unknown> = {
      messages: messages.value.map(m => ({
        role: m.role,
        content: m.content,
      })),
      stream: true,
      routing_strategy: routingStrategy.value,
    }

    // 固定模型策略时，传递选择的模型和深度思考配置
    if (routingStrategy.value === 'fixed' && selectedModel.value) {
      chatRequest.model = selectedModel.value.modelKey
      chatRequest.enable_reasoning = selectedModel.value.supportReasoning === 1 ? enableReasoning.value : false
    }

    // 调用流式API
    await streamChat(chatRequest)

    // 成功后清空输入框
    inputMessage.value = ''
  } catch (err) {
    const errorMsg = err instanceof Error ? err.message : '未知错误'
    message.error('对话失败: ' + errorMsg)
    // 移除失败的用户消息
    messages.value.pop()
    streamingContent.value = ''
    streamingThinking.value = ''
  } finally {
    loading.value = false
  }
}

// 流式调用
const streamChat = async (chatRequest: Record<string, unknown>) => {
  const url = `/api/internal/chat/completions`

  const response = await fetch(url, {
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

  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    const result = await response.json()
    if (result.code !== 0) {
      throw new Error(result.message || '请求失败')
    }
  }

  const reader = response.body?.getReader()
  const decoder = new TextDecoder()

  if (!reader) {
    throw new Error('无法获取响应流')
  }

  let buffer = ''
  let assistantContent = ''
  let thinkingContent = ''
  let inThinking = false

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
          // 移除 'data:' 前缀
          const content = line.substring(5).trim()
          if (content && content !== '[DONE]') {
            // 反转义换行符
            let unescapedContent = content.replace(/\\n/g, '\n')

            // 处理思考内容标识
            while (unescapedContent.length > 0) {
              if (unescapedContent.startsWith('[THINKING]')) {
                // 开始思考
                inThinking = true
                unescapedContent = unescapedContent.substring(10)
                continue
              }

              if (unescapedContent.startsWith('[/THINKING]')) {
                // 结束思考
                inThinking = false
                unescapedContent = unescapedContent.substring(11)
                continue
              }

              // 检查是否包含标识
              const thinkingStartIndex = unescapedContent.indexOf('[THINKING]')
              const thinkingEndIndex = unescapedContent.indexOf('[/THINKING]')

              if (inThinking) {
                if (thinkingEndIndex !== -1) {
                  // 思考内容结束
                  thinkingContent += unescapedContent.substring(0, thinkingEndIndex)
                  unescapedContent = unescapedContent.substring(thinkingEndIndex)
                  continue
                } else {
                  // 继续累积思考内容
                  thinkingContent += unescapedContent
                  break
                }
              } else {
                if (thinkingStartIndex !== -1) {
                  // 先处理思考开始前的内容
                  if (thinkingStartIndex > 0) {
                    assistantContent += unescapedContent.substring(0, thinkingStartIndex)
                  }
                  unescapedContent = unescapedContent.substring(thinkingStartIndex)
                  continue
                } else {
                  // 普通答案内容
                  assistantContent += unescapedContent
                  break
                }
              }
            }

            // 更新流式显示
            streamingThinking.value = thinkingContent
            streamingContent.value = assistantContent
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
        const content = line.substring(5).trim()
        if (content && content !== '[DONE]') {
          const unescapedContent = content.replace(/\\n/g, '\n')
          if (inThinking) {
            thinkingContent += unescapedContent
          } else {
            assistantContent += unescapedContent
          }
          streamingThinking.value = thinkingContent
          streamingContent.value = assistantContent
        }
      }
    }
  }

  // 完成后添加到消息列表
  const message: Message = {
    role: 'assistant',
    content: assistantContent,
  }

  // 如果有思考内容，添加到消息中
  if (thinkingContent) {
    message.thinking = thinkingContent
    message.answer = assistantContent
    message.thinkingExpanded = false  // 默认收起
  }

  messages.value.push(message)
  streamingContent.value = ''
  streamingThinking.value = ''
}

// 切换思考内容的展开/收起
const toggleThinking = (index: number) => {
  if (messages.value[index].thinkingExpanded !== undefined) {
    messages.value[index].thinkingExpanded = !messages.value[index].thinkingExpanded
  }
}

// 渲染 Markdown
const renderMarkdown = (content: string) => {
  return marked(content)
}

// 滚动到底部
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

// 处理回车发送
const handleEnter = (e: KeyboardEvent) => {
  if (!e.shiftKey) {
    e.preventDefault()
    sendMessage()
  }
}

// 监听流式内容变化，自动滚动
watch([streamingContent, streamingThinking], () => {
  nextTick(() => scrollToBottom())
})

onMounted(() => {
  loadModels()
})
</script>

<style scoped>
#chatPage {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px);
  background: #f5f5f5;
  overflow: hidden;
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 12px;
  padding: 8px 0;
}

.model-card {
  padding: 16px;
  border: 2px solid #e8e8e8;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
}

.model-card:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
}

.model-card.active {
  border-color: #1890ff;
  background: #e6f7ff;
}

.model-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.model-name {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.model-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-size: 12px;
  color: #8c8c8c;
}

.model-provider {
  color: #1890ff;
}

.model-price {
  color: #52c41a;
  font-weight: 500;
}

.model-desc {
  font-size: 12px;
  color: #595959;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 对话内容包裹层 */
.chat-wrapper {
  flex: 1;
  max-width: 900px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 20px;
}

/* 消息列表区域 */
.messages-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 0;
  min-height: 0;
  margin-bottom: 16px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #8c8c8c;
}

.empty-state h2 {
  font-size: 24px;
  margin-bottom: 8px;
  color: #262626;
}

.message-wrapper {
  margin-bottom: 16px;
}

.message {
  padding: 16px 20px;
  border-radius: 12px;
  max-width: 100%;
}

.message.user {
  background: #1890ff;
  color: white;
  margin-left: auto;
  margin-right: 0;
  max-width: 80%;
}

.message.assistant {
  background: white;
  border: 1px solid #e8e8e8;
  max-width: 100%;
}

.message-header {
  margin-bottom: 8px;
  font-size: 14px;
}

.message.user .message-header {
  color: rgba(255, 255, 255, 0.85);
}

.message.assistant .message-header {
  color: #8c8c8c;
}

.message-content {
  line-height: 1.6;
  word-wrap: break-word;
}

.message-content :deep(pre) {
  background: #f6f6f6;
  padding: 12px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 8px 0;
}

.message-content :deep(code) {
  background: #f6f6f6;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Monaco', 'Consolas', monospace;
}

.message-content :deep(pre code) {
  background: transparent;
  padding: 0;
}

/* 思考内容样式 */
.thinking-section {
  margin-bottom: 16px;
  border-radius: 6px;
  background: #f8f8f8;
  overflow: hidden;
  border: 1px solid #e0e0e0;
}

.thinking-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  cursor: pointer;
  user-select: none;
  transition: background 0.2s;
}

.thinking-toggle:hover {
  background: #ececec;
}

.thinking-title {
  font-size: 13px;
  color: #888;
  font-weight: 400;
}

.thinking-arrow {
  font-size: 11px;
  color: #888;
  transition: transform 0.3s;
}

.thinking-arrow.expanded {
  transform: rotate(180deg);
}

.thinking-detail {
  padding: 12px 14px;
  background: #fafafa;
  border-top: 1px solid #e0e0e0;
  color: #666;
  font-size: 15px;
  line-height: 1.6;
}

.thinking-detail :deep(p) {
  margin: 8px 0;
  color: #666;
}

.thinking-detail :deep(ul),
.thinking-detail :deep(ol) {
  margin: 8px 0;
  padding-left: 20px;
  color: #666;
}

.thinking-detail :deep(li) {
  margin: 4px 0;
}

.thinking-detail :deep(strong),
.thinking-detail :deep(b) {
  color: #555;
  font-weight: 500;
}

.thinking-detail :deep(code) {
  background: #f0f0f0;
  color: #666;
  font-size: 14px;
}

.answer-text {
  color: #262626;
}

/* 输入区域 */
.input-container {
  flex-shrink: 0;
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.options-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
  flex-wrap: wrap;
}

.model-select-button {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  padding: 8px 16px;
  height: auto;
  font-size: 14px;
  transition: all 0.3s;
}

.model-select-button:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.current-model {
  font-weight: 500;
}

.dropdown-icon {
  font-size: 12px;
  opacity: 0.6;
}

.reasoning-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f5f5f5;
  border-radius: 8px;
}

.reasoning-icon {
  font-size: 16px;
  color: #722ed1;
}

.reasoning-label {
  font-size: 14px;
  color: #262626;
  margin-right: 4px;
}

/* 路由策略选择器样式 */
.strategy-select-button {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  padding: 8px 16px;
  height: auto;
  font-size: 14px;
  background: linear-gradient(135deg, #f0f5ff 0%, #e6f4ff 100%);
  transition: all 0.3s;
}

.strategy-select-button:hover {
  border-color: #1890ff;
  background: linear-gradient(135deg, #e6f4ff 0%, #d6e4ff 100%);
}

.strategy-label {
  font-weight: 500;
  color: #1890ff;
}

.strategy-option {
  padding: 4px 0;
  min-width: 200px;
}

.strategy-option-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.strategy-option-label {
  font-size: 14px;
  font-weight: 500;
  color: #262626;
}

.strategy-option-desc {
  font-size: 12px;
  color: #8c8c8c;
}

.strategy-active {
  background: #e6f4ff;
}


.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-input {
  border-radius: 8px;
  font-size: 15px;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

/* 滚动条样式 */
.messages-container::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background-color: rgba(0, 0, 0, 0.3);
}
</style>
