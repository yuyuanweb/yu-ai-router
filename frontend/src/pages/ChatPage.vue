<template>
  <div id="chatPage" class="morandi-chat-page">
    <!-- 对话内容包裹层 -->
    <div class="chat-wrapper">
      <!-- 消息列表区域（带滚动条） -->
      <div class="messages-container" ref="messagesContainer">
        <div v-if="messages.length === 0" class="empty-state">
          <h2>开始与 AI 对话</h2>
          <p>选择一个模型，然后输入您的问题</p>
        </div>

        <div v-for="(msg, index) in messages" :key="index" class="message-wrapper" :class="msg.role">
          <!-- 角色标签 -->
          <div v-if="msg.role === 'assistant'" class="role-label assistant-label">
            <span>{{ getAIDisplayName() }}</span>
          </div>
          <div v-else class="role-label user-label">
            <span>你</span>
          </div>

          <div class="message" :class="msg.role">
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
            <div v-else class="message-content">
              <!-- 图片缩略图 -->
              <div v-if="msg.file && msg.file.type.startsWith('image/') && msg.filePreviewUrl" class="image-preview">
                <img :src="msg.filePreviewUrl" :alt="msg.file.name" class="image-thumbnail" />
                <div class="image-info">
                  <PictureOutlined style="font-size: 14px; color: rgba(255,255,255,0.85)" />
                  <span class="image-name">{{ msg.file.name }}</span>
                  <span class="image-size">({{ formatFileSize(msg.file.size) }})</span>
                </div>
              </div>
              <!-- PDF 文件信息 -->
              <div v-else-if="msg.file" class="message-file">
                <FilePdfOutlined style="font-size: 16px; color: #ff7a45" />
                <span class="file-name">{{ msg.file.name }}</span>
                <span class="file-size">({{ formatFileSize(msg.file.size) }})</span>
              </div>
              <div v-html="renderMarkdown(msg.content)"></div>
            </div>
          </div>
        </div>

        <!-- 流式响应中的消息 -->
        <div v-if="streamingContent || streamingThinking" class="message-wrapper assistant">
          <div class="role-label assistant-label">
            <span>{{ getAIDisplayName() }}</span>
          </div>
          <div class="message assistant">
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
        <div v-if="loading && !streamingContent" class="message-wrapper assistant">
          <div class="role-label assistant-label">
            <span>{{ getAIDisplayName() }}</span>
          </div>
          <div class="message assistant">
            <div class="message-content loading-message">
              <a-spin /> <span>思考中...</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 输入区域 (悬浮在底部) -->
      <div class="input-container">
        <!-- 模型选择和深度思考选项 -->
        <div class="options-bar">
          <!-- 路由策略选择 -->
          <a-dropdown :trigger="['click']">
            <a-button class="strategy-select-button capsule-btn" size="large">
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
          <a-button v-if="routingStrategy === 'fixed'" @click="showModelSelector" size="large" class="model-select-button capsule-btn">
            <template v-if="selectedModel">
              <AppstoreOutlined />
              <span class="current-model">{{ selectedModel.modelName }}</span>
              <a-tag v-if="selectedModel.supportReasoning === 1" color="purple" size="small" class="soft-tag">思考</a-tag>
            </template>
            <template v-else>
              <AppstoreOutlined />
              选择模型
            </template>
            <DownOutlined class="dropdown-icon" />
          </a-button>

          <!-- 深度思考开关 -->
          <a-button
            v-if="routingStrategy === 'fixed' && selectedModel?.supportReasoning === 1"
            :type="enableReasoning ? 'primary' : 'default'"
            class="capsule-btn"
            @click="enableReasoning = !enableReasoning"
            size="large"
          >
            <BulbOutlined />
            深度思考
          </a-button>

          <!-- 联网搜索开关 -->
          <a-button
            :type="enableWebSearch ? 'primary' : 'default'"
            class="capsule-btn"
            @click="enableWebSearch = !enableWebSearch"
            size="large"
          >
            <GlobalOutlined />
            联网搜索
          </a-button>

          <!-- PDF 解析按钮 -->
          <a-button
            size="large"
            :type="selectedPluginKey === 'pdf_parser' ? 'primary' : 'default'"
            @click="togglePlugin('pdf_parser')"
            class="capsule-btn"
          >
            <FilePdfOutlined />
            PDF 解析
          </a-button>

          <!-- 图片识别按钮 -->
          <a-button
            size="large"
            :type="selectedPluginKey === 'image_recognition' ? 'primary' : 'default'"
            @click="togglePlugin('image_recognition')"
            class="capsule-btn"
          >
            <PictureOutlined />
            图片识别
          </a-button>
        </div>

        <!-- 文件信息显示（在输入框上方） -->
        <div v-if="uploadedFile" class="file-info-bar">
          <div v-if="uploadedFile.type.startsWith('image/') && uploadedFilePreviewUrl" class="file-preview-container">
            <img :src="uploadedFilePreviewUrl" :alt="uploadedFile.name" class="file-preview-image" />
            <div class="file-info">
              <a-space>
                <PictureOutlined style="font-size: 16px; color: #52c41a" />
                <span class="file-name">{{ uploadedFile.name }}</span>
                <span class="file-size">{{ formatFileSize(uploadedFile.size) }}</span>
                <DeleteOutlined class="delete-icon" @click="clearFile" />
              </a-space>
            </div>
          </div>
          <div v-else>
            <a-space>
              <FilePdfOutlined style="font-size: 18px; color: #ff7a45" />
              <span class="file-name">{{ uploadedFile.name }}</span>
              <span class="file-size">{{ formatFileSize(uploadedFile.size) }}</span>
              <DeleteOutlined class="delete-icon" @click="clearFile" />
            </a-space>
          </div>
        </div>

        <div class="input-wrapper">
          <a-textarea
            v-model:value="inputMessage"
            :placeholder="getInputPlaceholder()"
            :auto-size="{ minRows: 1, maxRows: 6 }"
            :disabled="!canSendMessage || loading"
            @pressEnter="handleEnter"
            @paste="handlePaste"
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
      class="morandi-modal"
    >
      <div class="model-category-filter">
        <div class="category-btn" :class="{ active: activeCategory === 'all' }" @click="activeCategory = 'all'">全部模型</div>
        <div class="category-btn" :class="{ active: activeCategory === 'fast' }" @click="activeCategory = 'fast'">快速模型</div>
        <div class="category-btn" :class="{ active: activeCategory === 'reasoning' }" @click="activeCategory = 'reasoning'">深度思考</div>
      </div>

      <div class="model-grid">
        <div
          v-for="model in (activeCategory === 'all' ? allModels : activeCategory === 'fast' ? fastModels : reasoningModels)"
          :key="model.id"
          class="model-card"
          :class="{ active: selectedModel?.id === model.id }"
          @click="selectModelAndClose(model)"
        >
          <div class="model-header">
            <span class="model-name">{{ model.modelName }}</span>
            <a-tag v-if="model.supportReasoning === 1" color="purple" size="small" class="soft-tag">思考</a-tag>
          </div>
          <div class="model-info">
            <span class="model-provider">{{ model.providerDisplayName }}</span>
            <span class="model-price">¥{{ model.inputPrice }}/¥{{ model.outputPrice }}</span>
          </div>
          <div class="model-desc">{{ model.description }}</div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { listActiveModels } from '@/api/modelController'
import { listEnabledPlugins } from '@/api/pluginController'
import { message } from 'ant-design-vue'
import { marked } from 'marked'
import {
  AppstoreOutlined,
  DownOutlined,
  BulbOutlined,
  ThunderboltOutlined,
  GlobalOutlined,
  FilePdfOutlined,
  PictureOutlined,
  DeleteOutlined,
} from '@ant-design/icons-vue'

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
  file?: { name: string; size: number; type: string }
  filePreviewUrl?: string
  pluginKey?: string
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

// 插件相关
const enabledPlugins = ref<API.PluginConfigVO[]>([])
const enableWebSearch = ref(false)
const selectedPluginKey = ref<string>('')

// 文件上传
const uploadedFile = ref<File | null>(null)
const uploadedFilePreviewUrl = ref<string>('')

// 计算不同类别的模型
const fastModels = computed(() => allModels.value.filter(m => m.supportReasoning !== 1))
const reasoningModels = computed(() => allModels.value.filter(m => m.supportReasoning === 1))

// 判断是否可以发送消息
const canSendMessage = computed(() => {
  if (routingStrategy.value !== 'fixed') return allModels.value.length > 0
  return selectedModel.value !== null
})

// 加载模型列表
const loadModels = async () => {
  try {
    const res = await listActiveModels()
    if (res.data.data) {
      allModels.value = res.data.data
      if (allModels.value.length > 0) {
        selectedModel.value = allModels.value[0]
      }
    }
  } catch (err) {
    console.error('加载模型列表失败', err)
    message.error('加载模型列表失败')
  }
}

// 加载启用的插件列表
const loadEnabledPlugins = async () => {
  try {
    const res = await listEnabledPlugins()
    if (res.data.code === 0 && res.data.data) {
      enabledPlugins.value = res.data.data
    }
  } catch (err) {
    console.error('加载插件列表失败', err)
  }
}

// 切换插件
const togglePlugin = (pluginKey: string) => {
  if (selectedPluginKey.value === pluginKey) {
    selectedPluginKey.value = ''
    clearFile()
  } else {
    enableWebSearch.value = false
    selectedPluginKey.value = pluginKey
    if (uploadedFile.value) {
      const isValid = validateFileForPlugin(uploadedFile.value, pluginKey)
      if (!isValid) clearFile()
    }
  }
}

// 监听联网搜索开关变化
watch(enableWebSearch, (newValue) => {
  if (newValue) {
    selectedPluginKey.value = ''
    clearFile()
  }
})

// 处理粘贴事件
const handlePaste = (e: ClipboardEvent) => {
  const items = e.clipboardData?.items
  if (!items) return
  if (!selectedPluginKey.value || selectedPluginKey.value === 'web_search') return
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.kind === 'file') {
      e.preventDefault()
      const file = item.getAsFile()
      if (file) handleFileSelect(file)
      break
    }
  }
}

// 处理文件选择
const handleFileSelect = (file: File) => {
  if (!validateFileForPlugin(file, selectedPluginKey.value)) return
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    message.error('文件大小不能超过 10MB')
    return
  }
  uploadedFile.value = file
  if (file.type.startsWith('image/')) {
    uploadedFilePreviewUrl.value = URL.createObjectURL(file)
  }
  message.success(`文件已添加: ${file.name}`)
}

// 验证文件是否匹配插件
const validateFileForPlugin = (file: File, pluginKey: string): boolean => {
  const isImage = file.type.startsWith('image/')
  const isPdf = file.type === 'application/pdf'
  if (pluginKey === 'pdf_parser' && !isPdf) {
    message.error('PDF 解析只能上传 PDF 文件')
    return false
  } else if (pluginKey === 'image_recognition' && !isImage) {
    message.error('图片识别只能上传图片文件')
    return false
  }
  return true
}

// 清除文件
const clearFile = () => {
  if (uploadedFilePreviewUrl.value) {
    URL.revokeObjectURL(uploadedFilePreviewUrl.value)
    uploadedFilePreviewUrl.value = ''
  }
  uploadedFile.value = null
}

// 获取输入框提示文本
const getInputPlaceholder = () => {
  if (!canSendMessage.value) return '请先选择一个模型'
  if (selectedPluginKey.value === 'pdf_parser') return '粘贴 PDF 文件，然后输入您的问题...'
  if (selectedPluginKey.value === 'image_recognition') return '粘贴图片，然后输入您的问题...'
  return '输入您的问题... (Shift + Enter 换行)'
}

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

// 显示模型选择弹窗
const showModelSelector = () => {
  modelSelectorVisible.value = true
}

// 选择模型并关闭弹窗
const selectModelAndClose = (model: API.ModelVO) => {
  selectedModel.value = model
  if (model.supportReasoning !== 1) enableReasoning.value = false
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
  const option = routingStrategyOptions.find(opt => opt.value === routingStrategy.value)
  return option?.label || 'AI'
}

// 处理路由策略变更
const handleStrategyChange = ({ key }: { key: string }) => {
  routingStrategy.value = key as 'auto' | 'cost_first' | 'latency_first' | 'fixed'
}

// 发送消息
const sendMessage = async () => {
  if (!canSendMessage.value || !inputMessage.value.trim()) return

  const userMessage = inputMessage.value.trim()
  const userMsg: Message = { role: 'user', content: userMessage }

  if (uploadedFile.value) {
    userMsg.file = {
      name: uploadedFile.value.name,
      size: uploadedFile.value.size,
      type: uploadedFile.value.type,
    }
    if (uploadedFilePreviewUrl.value) userMsg.filePreviewUrl = uploadedFilePreviewUrl.value
  }

  if (enableWebSearch.value) {
    userMsg.pluginKey = 'web_search'
  } else if (selectedPluginKey.value) {
    userMsg.pluginKey = selectedPluginKey.value
  }

  messages.value.push(userMsg)
  loading.value = true
  streamingContent.value = ''
  streamingThinking.value = ''

  nextTick(() => scrollToBottom())

  try {
    if (uploadedFile.value) {
      await sendMessageWithFile()
    } else {
      await sendNormalMessage()
    }
    inputMessage.value = ''
    clearFile()
    selectedPluginKey.value = ''
  } catch (err) {
    const errorMsg = err instanceof Error ? err.message : '未知错误'
    message.error('对话失败: ' + errorMsg)
    messages.value.pop()
    streamingContent.value = ''
    streamingThinking.value = ''
  } finally {
    loading.value = false
  }
}

// 发送普通消息
const sendNormalMessage = async () => {
  const chatRequest: Record<string, unknown> = {
    messages: messages.value.map(m => ({ role: m.role, content: m.content })),
    stream: true,
    routing_strategy: routingStrategy.value,
  }

  if (routingStrategy.value === 'fixed' && selectedModel.value) {
    chatRequest.model = selectedModel.value.modelKey
    chatRequest.enable_reasoning = selectedModel.value.supportReasoning === 1 ? enableReasoning.value : false
  }

  if (enableWebSearch.value) chatRequest.plugin_key = 'web_search'
  await streamChat(chatRequest)
}

// 发送带文件的消息
const sendMessageWithFile = async () => {
  if (!uploadedFile.value) return
  const formData = new FormData()
  formData.append('file', uploadedFile.value)
  formData.append('messages', JSON.stringify(messages.value.map(m => ({ role: m.role, content: m.content }))))
  formData.append('stream', 'true')
  if (routingStrategy.value) formData.append('routing_strategy', routingStrategy.value)

  if (routingStrategy.value === 'fixed' && selectedModel.value && selectedModel.value.modelKey) {
    formData.append('model', selectedModel.value.modelKey)
    if (selectedModel.value.supportReasoning === 1) {
      formData.append('enable_reasoning', String(enableReasoning.value))
    }
  }
  if (selectedPluginKey.value) formData.append('plugin_key', selectedPluginKey.value)
  await streamChatWithFile(formData)
}

// 流式调用（带文件上传）
const streamChatWithFile = async (formData: FormData) => {
  const url = `/api/internal/chat/completions/upload`
  const response = await fetch(url, { method: 'POST', credentials: 'include', body: formData })
  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`)
  await processStream(response)
}

// 流式调用（普通消息）
const streamChat = async (chatRequest: Record<string, unknown>) => {
  const url = `/api/internal/chat/completions`
  const response = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(chatRequest),
  })
  if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`)
  const contentType = response.headers.get('content-type') || ''
  if (contentType.includes('application/json')) {
    const result = await response.json()
    if (result.code !== 0) throw new Error(result.message || '请求失败')
  }
  await processStream(response)
}

// 处理流式响应
const processStream = async (response: Response) => {
  const reader = response.body?.getReader()
  const decoder = new TextDecoder()
  if (!reader) throw new Error('无法获取响应流')

  let buffer = ''
  let assistantContent = ''
  let thinkingContent = ''

  while (true) {
    const { done, value } = await reader.read()
    if (done) break
    buffer += decoder.decode(value, { stream: true })
    const parts = buffer.split('\n\n')
    buffer = parts.pop() || ''
    for (const part of parts) {
      const lines = part.split('\n')
      for (const line of lines) {
        if (line.startsWith('data:')) {
          const data = line.substring(5).trim()
          if (!data) continue
          try {
            const streamResponse: API.StreamResponse = JSON.parse(data)
            if (streamResponse.choices && streamResponse.choices.length > 0) {
              const choice = streamResponse.choices[0]
              if (choice.finishReason === 'stop') continue
              const delta = choice.delta
              if (delta?.reasoningContent) {
                thinkingContent += delta.reasoningContent
                streamingThinking.value = thinkingContent
              }
              if (delta?.content) {
                assistantContent += delta.content
                streamingContent.value = assistantContent
              }
            }
          } catch (e) { console.error('解析 SSE 数据失败:', data, e) }
        }
      }
    }
  }

  if (buffer) {
    const lines = buffer.split('\n')
    for (const line of lines) {
      if (line.startsWith('data:')) {
        const data = line.substring(5).trim()
        if (data) {
          try {
            const streamResponse: API.StreamResponse = JSON.parse(data)
            if (streamResponse.choices && streamResponse.choices.length > 0) {
              const choice = streamResponse.choices[0]
              if (choice.finishReason === 'stop') continue
              const delta = choice.delta
              if (delta?.reasoningContent) {
                thinkingContent += delta.reasoningContent
                streamingThinking.value = thinkingContent
              }
              if (delta?.content) {
                assistantContent += delta.content
                streamingContent.value = assistantContent
              }
            }
          } catch (e) { console.error('解析剩余 SSE 数据失败:', data, e) }
        }
      }
    }
  }

  const msg: Message = { role: 'assistant', content: assistantContent }
  if (thinkingContent) {
    msg.thinking = thinkingContent
    msg.answer = assistantContent
    msg.thinkingExpanded = false
  }
  messages.value.push(msg)
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
  loadEnabledPlugins()
})
</script>

<style scoped>
/* 整体页面 */
.morandi-chat-page {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 72px);
  background: radial-gradient(circle at 50% 0%, #FAF7F4 0%, #F5F2EF 100%);
  overflow: hidden;
  position: relative;
}

/* 核心包裹层 */
.chat-wrapper {
  display: flex;
  flex-direction: column;
  flex: 1;
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  height: 100%;
  min-height: 0;
  position: relative;
  /* 关键：底部留出空间给悬浮的输入框，确保滑动到底部时内容能完全露出 */
  padding: 24px 24px 240px 24px;
}

/* 消息列表区域 */
.messages-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 8px; /* 给滚动条留点空间 */
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 120px 20px;
}
.empty-state h2 {
  font-size: 32px;
  margin-bottom: 12px;
  color: var(--text-main);
  font-weight: 600;
  letter-spacing: 2px;
}
.empty-state p {
  font-size: 15px;
  color: var(--text-secondary);
}

/* 消息气泡 */
.message-wrapper {
  margin-bottom: 28px;
  animation: messageFadeIn 0.5s cubic-bezier(0.25, 0.8, 0.25, 1);
  display: flex;
  flex-direction: column;
  gap: 6px;
}
@keyframes messageFadeIn {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
.message-wrapper.user { align-items: flex-end; }
.message-wrapper.assistant { align-items: flex-start; }
.role-label { font-size: 12px; font-weight: 600; color: var(--text-secondary); padding: 0 4px; }

.message.user {
  background: linear-gradient(135deg, #B8A9B8 0%, #C5B8C5 100%);
  color: #fff;
  padding: 16px 20px;
  border-radius: 20px 20px 4px 20px;
  max-width: 75%;
  box-shadow: 0 8px 24px rgba(184, 169, 184, 0.25);
}

.message.assistant {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  color: var(--text-main);
  padding: 20px 24px;
  border-radius: 20px 20px 20px 4px;
  max-width: 90%;
  box-shadow: var(--shadow-soft);
}

/* Markdown 内容 */
.message-content :deep(pre) {
  background: #F2EFEC;
  padding: 16px;
  border-radius: 12px;
  overflow-x: auto;
  margin: 12px 0;
  border: 1px dashed var(--border-color);
}
.message-content :deep(code) {
  background: #F2EFEC;
  padding: 3px 8px;
  border-radius: 6px;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 13px;
  color: #8B7896;
}
.message-content :deep(blockquote) {
  border-left: 4px solid var(--color-secondary);
  padding-left: 16px;
  margin: 12px 0;
  color: var(--text-secondary);
  background: rgba(169, 180, 194, 0.1);
  padding: 12px 16px;
  border-radius: 0 12px 12px 0;
}

/* 思考区 */
.thinking-section {
  margin-bottom: 16px;
  border-radius: 12px;
  background: rgba(184, 169, 184, 0.05);
  border: 1px dashed rgba(184, 169, 184, 0.3);
  overflow: hidden;
}
.thinking-toggle { padding: 10px 16px; cursor: pointer; }
.thinking-title { font-size: 13px; color: var(--color-primary); font-weight: 500; }
.thinking-detail {
  padding: 16px;
  background: rgba(250, 247, 244, 0.5);
  border-top: 1px dashed rgba(184, 169, 184, 0.2);
  color: var(--text-secondary);
}

/* =========================================
   悬浮输入框区域 (绝对定位在底部)
   ========================================= */
.input-container {
  position: absolute;
  bottom: 24px; /* 悬浮高度 */
  left: 50%;
  transform: translateX(-50%);
  width: calc(100% - 48px); /* 左右留白 */
  max-width: 860px;
  background: rgba(250, 247, 244, 0.85);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 16px 24px;
  box-shadow: 0 12px 40px rgba(184, 169, 184, 0.15), 0 0 0 1px rgba(255, 255, 255, 0.5) inset;
  border: 1px solid var(--border-color);
  z-index: 10;
}

.options-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  padding-bottom: 12px;
  border-bottom: 1px dashed var(--border-color);
  flex-wrap: wrap;
}

/* 统一所有胶囊按钮的样式 */
.capsule-btn {
  border-radius: 30px !important;
  border: 1px solid var(--border-color) !important;
  background: rgba(255, 255, 255, 0.6) !important;
  color: var(--text-main) !important;
  font-size: 13px;
  height: 36px;
  padding: 0 16px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.capsule-btn:hover {
  border-color: var(--color-primary) !important;
  color: var(--color-primary) !important;
  background: #fff !important;
  transform: translateY(-1px);
}

/* 选中状态（Primary） */
.capsule-btn.ant-btn-primary {
  background: var(--color-primary) !important;
  border-color: var(--color-primary) !important;
  color: #fff !important;
  box-shadow: 0 4px 12px rgba(184, 169, 184, 0.3) !important;
}

.chat-input {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  font-size: 15px;
  color: var(--text-main);
  resize: none;
  padding: 8px 0;
}
.chat-input:hover, .chat-input:focus {
  background: transparent !important;
  box-shadow: none !important;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-top: 8px;
}

.input-actions :deep(.ant-btn-primary) {
  height: 40px;
  padding: 0 32px;
  font-size: 15px;
  font-weight: 500;
  border-radius: 20px;
  background: var(--color-primary);
  border: none;
  box-shadow: 0 4px 12px rgba(184, 169, 184, 0.3);
  transition: all 0.3s ease;
}

.input-actions :deep(.ant-btn-primary:hover) {
  background: #C5B8C5;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(184, 169, 184, 0.4);
}

/* 文件预览 */
.file-info-bar {
  padding: 12px 16px;
  background: rgba(184, 169, 184, 0.05);
  border: 1px dashed var(--border-color);
  border-radius: 12px;
  margin-bottom: 12px;
}
.file-preview-container { display: flex; align-items: center; gap: 12px; }
.file-preview-image { width: 60px; height: 60px; object-fit: cover; border-radius: 8px; }
.file-name { font-weight: 500; color: var(--text-main); font-size: 13px; }
.file-size { font-size: 12px; color: var(--text-secondary); margin-left: 8px; }
.delete-icon { font-size: 16px; color: #C48383; cursor: pointer; transition: all 0.3s; }
.delete-icon:hover { color: #ff4d4f; transform: scale(1.1); }

/* 模型选择弹窗 */
.model-category-filter {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px dashed var(--border-color);
}
.category-btn {
  padding: 8px 24px;
  border-radius: 30px;
  font-size: 14px;
  cursor: pointer;
  color: var(--text-secondary);
  background: rgba(184, 169, 184, 0.05);
  transition: all 0.3s ease;
}
.category-btn.active {
  background: var(--color-primary);
  color: white;
  box-shadow: 0 4px 12px rgba(184, 169, 184, 0.3);
}
.category-btn:hover:not(.active) { background: rgba(184, 169, 184, 0.15); }

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  padding: 8px 0;
}
.model-card {
  padding: 24px;
  border: 1px dashed var(--border-color);
  border-radius: 16px;
  background: var(--bg-card);
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.25, 0.8, 0.25, 1);
}
.model-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-soft);
  transform: translateY(-4px);
}
.model-card.active {
  border: 2px solid var(--color-primary);
  background: rgba(184, 169, 184, 0.05);
}
.model-header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 12px; gap: 8px; }
.model-name { font-size: 16px; font-weight: 600; color: var(--text-main); }
.model-info { display: flex; justify-content: space-between; align-items: center; margin: 12px 0; font-size: 12px; }
.model-provider { color: var(--color-secondary); background: rgba(169, 180, 194, 0.15); padding: 3px 10px; border-radius: 6px; }
.model-price { color: var(--color-primary); font-weight: 600; }
.model-desc { font-size: 13px; color: var(--text-secondary); line-height: 1.5; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }

/* 滚动条极简化 */
.messages-container::-webkit-scrollbar { width: 6px; }
.messages-container::-webkit-scrollbar-thumb { background: rgba(184, 169, 184, 0.2); border-radius: 10px; }
.messages-container::-webkit-scrollbar-thumb:hover { background: rgba(184, 169, 184, 0.4); }
</style>
