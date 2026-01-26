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

          <!-- 联网搜索开关 -->
          <div class="plugin-switch">
            <GlobalOutlined class="plugin-icon" />
            <span class="plugin-label">联网搜索</span>
            <a-switch v-model:checked="enableWebSearch" />
          </div>

          <!-- PDF 解析按钮 -->
          <a-button
            size="large"
            :type="selectedPluginKey === 'pdf_parser' ? 'primary' : 'default'"
            @click="togglePlugin('pdf_parser')"
            class="plugin-button"
          >
            <FilePdfOutlined />
            PDF 解析
          </a-button>

          <!-- 图片识别按钮 -->
          <a-button
            size="large"
            :type="selectedPluginKey === 'image_recognition' ? 'primary' : 'default'"
            @click="togglePlugin('image_recognition')"
            class="plugin-button"
          >
            <PictureOutlined />
            图片识别
          </a-button>
        </div>

        <!-- 文件信息显示（在输入框上方） -->
        <div v-if="uploadedFile" class="file-info-bar">
          <!-- 图片预览 -->
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
          <!-- PDF 文件信息 -->
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
            :auto-size="{ minRows: 2, maxRows: 6 }"
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
  file?: {
    name: string
    size: number
    type: string
  }
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

// 切换插件（PDF 解析 / 图片识别）
const togglePlugin = (pluginKey: string) => {
  if (selectedPluginKey.value === pluginKey) {
    // 如果已选中，则取消
    selectedPluginKey.value = ''
    clearFile()
  } else {
    // 选中新插件（互斥，关闭联网搜索）
    enableWebSearch.value = false
    selectedPluginKey.value = pluginKey
    // 如果已有文件，检查是否匹配
    if (uploadedFile.value) {
      const isValid = validateFileForPlugin(uploadedFile.value, pluginKey)
      if (!isValid) {
        clearFile()
      }
    }
  }
}

// 监听联网搜索开关变化
watch(enableWebSearch, (newValue) => {
  if (newValue) {
    // 启用联网搜索时，关闭其他插件
    selectedPluginKey.value = ''
    clearFile()
  }
})

// 处理粘贴事件
const handlePaste = (e: ClipboardEvent) => {
  const items = e.clipboardData?.items
  if (!items) return

  // 只有选择了 PDF 解析或图片识别才能粘贴文件
  if (!selectedPluginKey.value || selectedPluginKey.value === 'web_search') {
    return
  }

  for (let i = 0; i < items.length; i++) {
    const item = items[i]

    // 检查是否是文件
    if (item.kind === 'file') {
      e.preventDefault()
      const file = item.getAsFile()
      if (file) {
        handleFileSelect(file)
      }
      break
    }
  }
}

// 处理文件选择
const handleFileSelect = (file: File) => {
  // 验证文件
  if (!validateFileForPlugin(file, selectedPluginKey.value)) {
    return
  }

  // 检查文件大小（10MB）
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    message.error('文件大小不能超过 10MB')
    return
  }

  uploadedFile.value = file

  // 如果是图片，生成预览 URL
  if (file.type.startsWith('image/')) {
    uploadedFilePreviewUrl.value = URL.createObjectURL(file)
  }

  message.success(`文件已添加: ${file.name}`)
}

// 验证文件是否匹配插件
const validateFileForPlugin = (file: File, pluginKey: string): boolean => {
  const isImage = file.type.startsWith('image/')
  const isPdf = file.type === 'application/pdf'

  if (pluginKey === 'pdf_parser') {
    if (!isPdf) {
      message.error('PDF 解析只能上传 PDF 文件')
      return false
    }
  } else if (pluginKey === 'image_recognition') {
    if (!isImage) {
      message.error('图片识别只能上传图片文件')
      return false
    }
  }

  return true
}

// 清除文件
const clearFile = () => {
  // 释放预览 URL
  if (uploadedFilePreviewUrl.value) {
    URL.revokeObjectURL(uploadedFilePreviewUrl.value)
    uploadedFilePreviewUrl.value = ''
  }
  uploadedFile.value = null
}

// 获取输入框提示文本
const getInputPlaceholder = () => {
  if (!canSendMessage.value) {
    return '请先选择一个模型'
  }
  if (selectedPluginKey.value === 'pdf_parser') {
    return '粘贴 PDF 文件，然后输入您的问题...'
  }
  if (selectedPluginKey.value === 'image_recognition') {
    return '粘贴图片，然后输入您的问题...'
  }
  return '输入您的问题...'
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
  const userMsg: Message = {
    role: 'user',
    content: userMessage,
  }

  // 如果有上传的文件，添加文件信息和预览
  if (uploadedFile.value) {
    userMsg.file = {
      name: uploadedFile.value.name,
      size: uploadedFile.value.size,
      type: uploadedFile.value.type,
    }
    // 保存图片预览 URL
    if (uploadedFilePreviewUrl.value) {
      userMsg.filePreviewUrl = uploadedFilePreviewUrl.value
    }
  }

  // 如果启用了插件，记录插件信息
  if (enableWebSearch.value) {
    userMsg.pluginKey = 'web_search'
  } else if (selectedPluginKey.value) {
    userMsg.pluginKey = selectedPluginKey.value
  }

  messages.value.push(userMsg)

  loading.value = true
  streamingContent.value = ''
  streamingThinking.value = ''

  // 滚动到底部
  nextTick(() => scrollToBottom())

  try {
    // 如果有上传文件，使用文件上传接口
    if (uploadedFile.value) {
      await sendMessageWithFile()
    } else {
      await sendNormalMessage()
    }

    // 成功后清空输入框和文件
    inputMessage.value = ''
    clearFile()
    // 清空插件选择（但保留联网搜索开关状态）
    selectedPluginKey.value = ''
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

// 发送普通消息（不带文件）
const sendNormalMessage = async () => {
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

  // 如果启用了联网搜索，添加插件信息
  if (enableWebSearch.value) {
    chatRequest.plugin_key = 'web_search'
  }

  // 调用流式API
  await streamChat(chatRequest)
}

// 发送带文件的消息
const sendMessageWithFile = async () => {
  if (!uploadedFile.value) {
    return
  }

  // 构建 FormData
  const formData = new FormData()
  formData.append('file', uploadedFile.value)

  // 添加消息列表（JSON 字符串）
  const messagesData = messages.value.map(m => ({
    role: m.role,
    content: m.content,
  }))
  formData.append('messages', JSON.stringify(messagesData))

  // 添加流式参数
  formData.append('stream', 'true')

  // 添加路由策略
  if (routingStrategy.value) {
    formData.append('routing_strategy', routingStrategy.value)
  }

  // 固定模型策略时，传递选择的模型
  if (routingStrategy.value === 'fixed' && selectedModel.value && selectedModel.value.modelKey) {
    formData.append('model', selectedModel.value.modelKey)
    if (selectedModel.value.supportReasoning === 1) {
      formData.append('enable_reasoning', String(enableReasoning.value))
    }
  }

  // 传递选择的插件
  if (selectedPluginKey.value) {
    formData.append('plugin_key', selectedPluginKey.value)
  }

  // 调用文件上传接口
  await streamChatWithFile(formData)
}

// 流式调用（带文件上传）
const streamChatWithFile = async (formData: FormData) => {
  const url = `/api/internal/chat/completions/upload`

  const response = await fetch(url, {
    method: 'POST',
    credentials: 'include', // 携带 session cookie
    body: formData,
  })

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`)
  }

  await processStream(response)
}

// 流式调用（普通消息）
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

  await processStream(response)
}

// 处理流式响应（通用方法）
const processStream = async (response: Response) => {
  const reader = response.body?.getReader()
  const decoder = new TextDecoder()

  if (!reader) {
    throw new Error('无法获取响应流')
  }

  let buffer = ''
  let assistantContent = ''
  let thinkingContent = ''

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
          const data = line.substring(5).trim()

          if (!data) {
            continue
          }

          try {
            // 解析 JSON 格式的 StreamResponse
            const streamResponse: API.StreamResponse = JSON.parse(data)

            // 检查是否结束（finishReason 为 "stop"）
            if (streamResponse.choices && streamResponse.choices.length > 0) {
              const choice = streamResponse.choices[0]

              // 如果有 finishReason，表示流结束
              if (choice.finishReason === 'stop') {
                continue
              }

              const delta = choice.delta

              // 处理深度思考内容
              if (delta?.reasoningContent) {
                thinkingContent += delta.reasoningContent
                streamingThinking.value = thinkingContent
              }

              // 处理普通文本内容
              if (delta?.content) {
                assistantContent += delta.content
                streamingContent.value = assistantContent
              }
            }
          } catch (e) {
            console.error('解析 SSE 数据失败:', data, e)
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
        const data = line.substring(5).trim()
        if (data) {
          try {
            const streamResponse: API.StreamResponse = JSON.parse(data)
            if (streamResponse.choices && streamResponse.choices.length > 0) {
              const choice = streamResponse.choices[0]
              // 如果有 finishReason，表示流结束
              if (choice.finishReason === 'stop') {
                continue
              }
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
          } catch (e) {
            console.error('解析剩余 SSE 数据失败:', data, e)
          }
        }
      }
    }
  }

  // 完成后添加到消息列表
  const msg: Message = {
    role: 'assistant',
    content: assistantContent,
  }

  // 如果有思考内容，添加到消息中
  if (thinkingContent) {
    msg.thinking = thinkingContent
    msg.answer = assistantContent
    msg.thinkingExpanded = false  // 默认收起
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
#chatPage {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px);
  background: linear-gradient(180deg, #f0f2f5 0%, #fafafa 100%);
  overflow: hidden;
}

.model-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
  padding: 12px 0;
}

.model-card {
  padding: 18px;
  border: 2px solid #f0f0f0;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
  position: relative;
}

.model-card:hover {
  border-color: #1890ff;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.15);
  transform: translateY(-2px);
}

.model-card.active {
  border-color: #1890ff;
  background: linear-gradient(135deg, #e6f7ff 0%, #f0f9ff 100%);
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.2);
}

.model-card.active::after {
  content: '✓';
  position: absolute;
  top: 8px;
  left: 8px;
  width: 22px;
  height: 22px;
  background: #1890ff;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

.model-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
  gap: 8px;
}

.model-card.active .model-header {
  padding-left: 28px;
}

.model-name {
  font-size: 17px;
  font-weight: 700;
  color: #262626;
  letter-spacing: -0.3px;
  flex: 1;
}

.model-card.active .model-name {
  color: #1890ff;
}

.model-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  font-size: 12px;
  color: #8c8c8c;
}

.model-provider {
  color: #1890ff;
  font-weight: 500;
  background: #f0f5ff;
  padding: 3px 10px;
  border-radius: 6px;
  font-size: 11px;
}

.model-price {
  color: #52c41a;
  font-weight: 600;
  font-size: 12px;
}

.model-desc {
  font-size: 13px;
  color: #595959;
  line-height: 1.5;
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
  max-width: 1000px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 24px;
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
  padding: 80px 20px;
  color: #8c8c8c;
}

.empty-state h2 {
  font-size: 28px;
  margin-bottom: 12px;
  color: #262626;
  font-weight: 600;
  background: linear-gradient(135deg, #1890ff 0%, #722ed1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.empty-state p {
  font-size: 15px;
  color: #8c8c8c;
}

.message-wrapper {
  margin-bottom: 24px;
  animation: messageSlideIn 0.3s ease-out;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.message-wrapper.user {
  align-items: flex-end;
}

.message-wrapper.assistant {
  align-items: flex-start;
}

/* 角色标签 */
.role-label {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 0;
  margin-bottom: 2px;
}

.assistant-label {
  color: #8c8c8c;
  padding-left: 4px;
}

.user-label {
  color: #1890ff;
  padding-right: 4px;
}

@keyframes messageSlideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message {
  padding: 14px 18px;
  border-radius: 16px;
  max-width: 100%;
  transition: all 0.3s;
  line-height: 1.7;
  position: relative;
}

.message.user {
  background: linear-gradient(135deg, #1890ff 0%, #40a9ff 100%);
  color: white;
  width: fit-content;
  max-width: 65%;
  box-shadow: 0 3px 14px rgba(24, 144, 255, 0.28);
  word-break: break-word;
  border-bottom-right-radius: 4px;
}

.message.user:hover {
  box-shadow: 0 5px 20px rgba(24, 144, 255, 0.38);
  transform: translateY(-1px);
}

.message.assistant {
  background: white;
  border: 1px solid #f0f0f0;
  max-width: 90%;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  border-bottom-left-radius: 4px;
}

.message.assistant:hover {
  box-shadow: 0 3px 16px rgba(0, 0, 0, 0.08);
  border-color: #e8e8e8;
}

.message-content {
  line-height: 1.7;
  word-wrap: break-word;
  font-size: 14px;
}

.message.user .message-content {
  color: rgba(255, 255, 255, 0.98);
  font-weight: 400;
}

.loading-message {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #8c8c8c;
  font-size: 14px;
}

.message-content :deep(pre) {
  background: #282c34;
  padding: 16px;
  border-radius: 10px;
  overflow-x: auto;
  margin: 12px 0;
  border: 1px solid #e0e0e0;
}

.message-content :deep(code) {
  background: #f0f2f5;
  padding: 3px 8px;
  border-radius: 6px;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
  font-size: 13px;
  color: #d63384;
  border: 1px solid #e8e8e8;
}

.message-content :deep(pre code) {
  background: transparent;
  padding: 0;
  border: none;
  color: #abb2bf;
  font-size: 13px;
}

.message-content :deep(blockquote) {
  border-left: 4px solid #1890ff;
  padding-left: 16px;
  margin: 12px 0;
  color: #595959;
  background: #f0f5ff;
  padding: 12px 16px;
  border-radius: 8px;
}

.message-content :deep(ul),
.message-content :deep(ol) {
  padding-left: 24px;
  margin: 8px 0;
}

.message-content :deep(li) {
  margin: 6px 0;
}

.message-content :deep(a) {
  color: #1890ff;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.3s;
}

.message-content :deep(a:hover) {
  border-bottom-color: #1890ff;
}

/* 思考内容样式 */
.thinking-section {
  margin-bottom: 18px;
  border-radius: 12px;
  background: linear-gradient(135deg, #faf5ff 0%, #f0f0f0 100%);
  overflow: hidden;
  border: 1px solid #d3adf7;
}

.thinking-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  cursor: pointer;
  user-select: none;
  transition: all 0.3s;
}

.thinking-toggle:hover {
  background: rgba(114, 46, 209, 0.05);
}

.thinking-title {
  font-size: 13px;
  color: #722ed1;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.thinking-title::before {
  content: '💭';
  font-size: 14px;
}

.thinking-arrow {
  font-size: 11px;
  color: #722ed1;
  transition: transform 0.3s;
}

.thinking-arrow.expanded {
  transform: rotate(180deg);
}

.thinking-detail {
  padding: 16px 18px;
  background: rgba(250, 245, 255, 0.5);
  border-top: 1px solid #efdbff;
  color: #595959;
  font-size: 14px;
  line-height: 1.7;
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
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
  border: 1px solid #f0f0f0;
}

.options-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e8e8e8;
  flex-wrap: wrap;
}

.model-select-button {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 1px solid #e8e8e8;
  border-radius: 10px;
  padding: 8px 14px;
  height: auto;
  font-size: 13px;
  transition: all 0.3s;
  background: #fafafa;
}

.model-select-button:hover {
  border-color: #1890ff;
  color: #1890ff;
  background: #f0f5ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
  transform: translateY(-1px);
}

.current-model {
  font-weight: 600;
  color: #262626;
}

.dropdown-icon {
  font-size: 11px;
  opacity: 0.5;
}

.reasoning-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  background: linear-gradient(135deg, #f9f0ff 0%, #efdbff 100%);
  border: 1px solid #d3adf7;
  border-radius: 10px;
  transition: all 0.3s;
}

.reasoning-switch:hover {
  background: linear-gradient(135deg, #f9f0ff 0%, #d3adf7 100%);
  box-shadow: 0 2px 8px rgba(114, 46, 209, 0.15);
}

.reasoning-icon {
  font-size: 15px;
  color: #722ed1;
}

.reasoning-label {
  font-size: 13px;
  color: #722ed1;
  font-weight: 500;
  margin-right: 4px;
}

/* 路由策略选择器样式 */
.strategy-select-button {
  display: flex;
  align-items: center;
  gap: 8px;
  border: 1px solid #b5f5ec;
  border-radius: 10px;
  padding: 8px 14px;
  height: auto;
  font-size: 13px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e6f7ff 100%);
  transition: all 0.3s;
}

.strategy-select-button:hover {
  border-color: #1890ff;
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
  transform: translateY(-1px);
}

.strategy-label {
  font-weight: 600;
  color: #1890ff;
}

.strategy-option {
  padding: 8px 12px;
  min-width: 260px;
  border-radius: 8px;
  transition: all 0.3s;
}

.strategy-option:hover {
  background: #f5f5f5;
}

.strategy-option-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.strategy-option-label {
  font-size: 14px;
  font-weight: 600;
  color: #262626;
}

.strategy-option-desc {
  font-size: 12px;
  color: #8c8c8c;
  line-height: 1.5;
}

.strategy-active {
  background: linear-gradient(135deg, #e6f7ff 0%, #bae7ff 100%);
  border: 1px solid #91d5ff;
}

.strategy-active .strategy-option-label {
  color: #1890ff;
}

/* 插件开关 */
.plugin-switch {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  background: linear-gradient(135deg, #f0f5ff 0%, #e8f4ff 100%);
  border: 1px solid #d6e4ff;
  border-radius: 10px;
  transition: all 0.3s;
}

.plugin-switch:hover {
  background: linear-gradient(135deg, #e6f4ff 0%, #d6e4ff 100%);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

.plugin-icon {
  font-size: 15px;
  color: #1890ff;
}

.plugin-label {
  font-size: 13px;
  color: #1890ff;
  font-weight: 500;
  margin-right: 4px;
}

/* 插件按钮 */
.plugin-button {
  display: flex;
  align-items: center;
  gap: 6px;
  border-radius: 10px;
  font-size: 13px;
  border: 1px solid #e8e8e8;
  transition: all 0.3s;
  padding: 8px 14px;
  background: white;
  color: #595959;
}

.plugin-button:hover {
  border-color: #1890ff;
  color: #1890ff;
  background: #f0f5ff;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
  transform: translateY(-1px);
}

/* Primary 状态的插件按钮 */
.plugin-button.ant-btn-primary {
  background: linear-gradient(135deg, #1890ff 0%, #096dd9 100%) !important;
  border-color: #1890ff !important;
  color: white !important;
  font-weight: 500;
}

.plugin-button.ant-btn-primary:hover {
  background: linear-gradient(135deg, #40a9ff 0%, #1890ff 100%) !important;
  border-color: #40a9ff !important;
  color: white !important;
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
}

/* 文件信息栏 */
.file-info-bar {
  padding: 14px 16px;
  background: linear-gradient(135deg, #e6fffb 0%, #e6f7ff 100%);
  border: 1px solid #b5f5ec;
  border-radius: 12px;
  margin-bottom: 14px;
  transition: all 0.3s;
}

.file-info-bar:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.file-preview-container {
  display: flex;
  align-items: center;
  gap: 14px;
}

.file-preview-image {
  width: 80px;
  height: 80px;
  object-fit: cover;
  border-radius: 8px;
  border: 2px solid #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.file-info {
  flex: 1;
}

.file-name {
  font-weight: 600;
  color: #262626;
  font-size: 14px;
}

.file-size {
  font-size: 12px;
  color: #8c8c8c;
  margin-left: 8px;
}

.delete-icon {
  font-size: 16px;
  color: #ff7875;
  cursor: pointer;
  transition: all 0.3s;
  padding: 4px;
}

.delete-icon:hover {
  color: #ff4d4f;
  background: rgba(255, 77, 79, 0.1);
  border-radius: 4px;
  transform: scale(1.15);
}

/* 消息中的文件信息 */
.message-file {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.25);
  border-radius: 10px;
  margin-bottom: 10px;
  font-size: 13px;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.message.user .message-file {
  background: rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 1);
  border-color: rgba(255, 255, 255, 0.25);
}

.message.user .message-file .file-name {
  color: rgba(255, 255, 255, 1);
  font-weight: 500;
}

.message.user .message-file .file-size {
  color: rgba(255, 255, 255, 0.85);
}

/* 图片预览 */
.image-preview {
  margin-bottom: 12px;
  max-width: 240px;
}

.image-thumbnail {
  width: 100%;
  max-width: 240px;
  max-height: 240px;
  border-radius: 12px;
  border: 2px solid rgba(255, 255, 255, 0.4);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  display: block;
  object-fit: cover;
  transition: all 0.3s;
}

.image-thumbnail:hover {
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.25);
  transform: scale(1.02);
}

.image-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.9);
}

.image-name {
  font-weight: 500;
}

.image-size {
  opacity: 0.8;
}

.input-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.chat-input {
  border-radius: 12px;
  font-size: 15px;
  border: 1px solid #e8e8e8;
  transition: all 0.3s;
}

.chat-input:hover {
  border-color: #d9d9d9;
}

.chat-input:focus {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.1);
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.input-actions :deep(.ant-btn-primary) {
  height: 40px;
  padding: 0 32px;
  font-size: 15px;
  font-weight: 500;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.2);
}

.input-actions :deep(.ant-btn-primary:hover) {
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.3);
  transform: translateY(-1px);
}

/* 滚动条样式 */
.messages-container::-webkit-scrollbar {
  width: 8px;
}

.messages-container::-webkit-scrollbar-track {
  background: transparent;
  margin: 4px;
}

.messages-container::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, #d9d9d9 0%, #bfbfbf 100%);
  border-radius: 10px;
  border: 2px solid transparent;
  background-clip: padding-box;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, #bfbfbf 0%, #8c8c8c 100%);
  background-clip: padding-box;
}
</style>
