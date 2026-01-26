<template>
  <div id="pluginManagePage">
    <h2>插件管理</h2>
    
    <a-card>
      <a-space style="margin-bottom: 16px">
        <a-button type="primary" @click="loadPlugins">
          <template #icon><ReloadOutlined /></template>
          刷新列表
        </a-button>
        <a-button @click="handleReloadAllPlugins">
          <template #icon><ThunderboltOutlined /></template>
          重新加载所有插件
        </a-button>
      </a-space>

      <a-table
        :columns="columns"
        :data-source="plugins"
        :loading="loading"
        :pagination="false"
        row-key="id"
      >
        <!-- 插件名称 -->
        <template #pluginName="{ record }">
          <a-tag :color="getPluginColor(record.pluginKey)">
            {{ record.pluginName }}
          </a-tag>
        </template>

        <!-- 插件类型 -->
        <template #pluginType="{ record }">
          <a-tag v-if="record.pluginType === 'builtin'" color="blue">内置</a-tag>
          <a-tag v-else color="green">自定义</a-tag>
        </template>

        <!-- 状态 -->
        <template #status="{ record }">
          <a-tag v-if="record.status === 'active'" color="success">启用</a-tag>
          <a-tag v-else color="default">禁用</a-tag>
        </template>

        <!-- 配置 -->
        <template #config="{ record }">
          <a-button type="link" size="small" @click="showConfigDetail(record)">
            查看配置
          </a-button>
        </template>

        <!-- 操作 -->
        <template #action="{ record }">
          <a-space>
            <a-button 
              v-if="record.status === 'active'" 
              size="small" 
              danger
              @click="handleDisable(record)"
            >
              禁用
            </a-button>
            <a-button 
              v-else 
              size="small" 
              type="primary"
              @click="handleEnable(record)"
            >
              启用
            </a-button>
            <a-button size="small" @click="handleEdit(record)">编辑</a-button>
            <a-button size="small" @click="handleReload(record)">重载</a-button>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 编辑插件弹窗 -->
    <a-modal
      v-model:open="editModalVisible"
      :title="`编辑插件 - ${currentPlugin?.pluginName}`"
      @ok="handleUpdate"
      @cancel="editModalVisible = false"
      width="600px"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="插件名称">
          <a-input v-model:value="editForm.pluginName" />
        </a-form-item>
        
        <a-form-item label="插件描述">
          <a-textarea v-model:value="editForm.description" :rows="3" />
        </a-form-item>
        
        <a-form-item label="插件配置（JSON）">
          <a-textarea v-model:value="editForm.config" :rows="6" />
        </a-form-item>
        
        <a-form-item label="优先级">
          <a-input-number v-model:value="editForm.priority" :min="0" :max="1000" style="width: 100%" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 配置详情弹窗 -->
    <a-modal
      v-model:open="configModalVisible"
      :title="`${currentPlugin?.pluginName} - 配置详情`"
      :footer="null"
      width="600px"
    >
      <pre style="background: #f5f5f5; padding: 16px; border-radius: 4px; overflow: auto">{{ formatConfig(currentPlugin?.config) }}</pre>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, ThunderboltOutlined } from '@ant-design/icons-vue'
import {
  listPlugins,
  enablePlugin,
  disablePlugin,
  updatePlugin,
  reloadPlugin,
  reloadAllPlugins,
} from '@/api/pluginController'

// 列定义
const columns = [
  {
    title: '插件标识',
    dataIndex: 'pluginKey',
    key: 'pluginKey',
    width: 150,
  },
  {
    title: '插件名称',
    dataIndex: 'pluginName',
    key: 'pluginName',
    slots: { customRender: 'pluginName' },
    width: 120,
  },
  {
    title: '类型',
    dataIndex: 'pluginType',
    key: 'pluginType',
    slots: { customRender: 'pluginType' },
    width: 100,
  },
  {
    title: '描述',
    dataIndex: 'description',
    key: 'description',
  },
  {
    title: '配置',
    key: 'config',
    slots: { customRender: 'config' },
    width: 100,
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    key: 'priority',
    width: 80,
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 80,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 250,
  },
]

const plugins = ref<API.PluginConfigVO[]>([])
const loading = ref(false)
const editModalVisible = ref(false)
const configModalVisible = ref(false)
const currentPlugin = ref<API.PluginConfigVO | null>(null)

// 编辑表单
const editForm = ref<{
  id?: number
  pluginName?: string
  description?: string
  config?: string
  priority?: number
}>({
  id: undefined,
  pluginName: '',
  description: '',
  config: '',
  priority: 100,
})

// 加载插件列表
const loadPlugins = async () => {
  loading.value = true
  try {
    const res = await listPlugins()
    if (res.data.code === 0 && res.data.data) {
      plugins.value = res.data.data
    }
  } catch (err) {
    console.error('加载插件列表失败', err)
    message.error('加载插件列表失败')
  } finally {
    loading.value = false
  }
}

// 启用插件
const handleEnable = async (record: API.PluginConfigVO) => {
  if (!record.pluginKey) return
  try {
    const res = await enablePlugin({ pluginKey: record.pluginKey })
    if (res.data.code === 0) {
      message.success('插件已启用')
      loadPlugins()
    }
  } catch {
    message.error('启用失败')
  }
}

// 禁用插件
const handleDisable = async (record: API.PluginConfigVO) => {
  if (!record.pluginKey) return
  try {
    const res = await disablePlugin({ pluginKey: record.pluginKey })
    if (res.data.code === 0) {
      message.success('插件已禁用')
      loadPlugins()
    }
  } catch {
    message.error('禁用失败')
  }
}

// 编辑插件
const handleEdit = (record: API.PluginConfigVO) => {
  currentPlugin.value = record
  editForm.value = {
    id: typeof record.id === 'number' ? record.id : parseInt(String(record.id)),
    pluginName: record.pluginName || '',
    description: record.description || '',
    config: record.config || '',
    priority: record.priority || 100,
  }
  editModalVisible.value = true
}

// 更新插件
const handleUpdate = async () => {
  try {
    const res = await updatePlugin(editForm.value)
    if (res.data.code === 0) {
      message.success('更新成功')
      editModalVisible.value = false
      loadPlugins()
    }
  } catch {
    message.error('更新失败')
  }
}

// 重载插件
const handleReload = async (record: API.PluginConfigVO) => {
  if (!record.pluginKey) return
  try {
    const res = await reloadPlugin({ pluginKey: record.pluginKey })
    if (res.data.code === 0) {
      message.success('重载成功')
    }
  } catch {
    message.error('重载失败')
  }
}

// 重载所有插件
const handleReloadAllPlugins = async () => {
  try {
    const res = await reloadAllPlugins()
    if (res.data.code === 0) {
      message.success('所有插件重载成功')
    }
  } catch {
    message.error('重载失败')
  }
}

// 显示配置详情
const showConfigDetail = (record: API.PluginConfigVO) => {
  currentPlugin.value = record
  configModalVisible.value = true
}

// 格式化配置 JSON
const formatConfig = (config: string | undefined) => {
  if (!config) return '{}'
  try {
    return JSON.stringify(JSON.parse(config), null, 2)
  } catch {
    return config
  }
}

// 获取插件颜色
const getPluginColor = (pluginKey: string | undefined) => {
  switch (pluginKey) {
    case 'web_search':
      return 'blue'
    case 'pdf_parser':
      return 'orange'
    case 'image_recognition':
      return 'purple'
    default:
      return 'default'
  }
}

onMounted(() => {
  loadPlugins()
})
</script>

<style scoped>
#pluginManagePage {
  padding: 24px;
}

h2 {
  margin-bottom: 20px;
}
</style>
