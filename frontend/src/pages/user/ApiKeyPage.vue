<template>
  <div class="api-key-page">
    <a-card title="API Key 管理">
      <template #extra>
        <a-button type="primary" @click="showCreateModal">创建 API Key</a-button>
      </template>

      <!-- API Key 列表 -->
      <a-table
        :columns="columns"
        :data-source="apiKeys"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'active' ? 'green' : 'red'">
              {{ record.status === 'active' ? '有效' : '已撤销' }}
            </a-tag>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-space>
              <a-button
                size="small"
                danger
                :disabled="record.status !== 'active'"
                @click="handleRevoke(record)"
              >
                撤销
              </a-button>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 创建 API Key 弹窗 -->
    <a-modal
      v-model:open="createModalVisible"
      title="创建 API Key"
      @ok="handleCreate"
      :confirm-loading="creating"
    >
      <a-form :model="createForm" :label-col="{ span: 6 }">
        <a-form-item label="Key 名称" name="keyName">
          <a-input
            v-model:value="createForm.keyName"
            placeholder="请输入 Key 名称（可选）"
          />
        </a-form-item>
      </a-form>

      <!-- 新创建的 Key -->
      <a-alert
        v-if="newApiKey"
        type="success"
        message="API Key 创建成功！"
        style="margin-top: 16px"
      >
        <template #description>
          <div style="margin-top: 8px">
            <strong>请复制保存以下 API Key，关闭后将无法再次查看完整内容：</strong>
            <a-input
              :value="newApiKey.keyValue"
              readonly
              style="margin-top: 8px"
            >
              <template #addonAfter>
                <a-button size="small" @click="copyApiKey">复制</a-button>
              </template>
            </a-input>
          </div>
        </template>
      </a-alert>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import {
  listMyApiKeys,
  createApiKey,
  revokeApiKey,
} from '@/api/apiKeyController'

const columns = [
  {
    title: 'Key 名称',
    dataIndex: 'keyName',
    key: 'keyName',
  },
  {
    title: 'API Key',
    dataIndex: 'keyValue',
    key: 'keyValue',
  },
  {
    title: '状态',
    key: 'status',
    dataIndex: 'status',
  },
  {
    title: '已使用Token数',
    dataIndex: 'totalTokens',
    key: 'totalTokens',
  },
  {
    title: '最后使用时间',
    dataIndex: 'lastUsedTime',
    key: 'lastUsedTime',
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
  },
  {
    title: '操作',
    key: 'action',
  },
]

const apiKeys = ref<API.ApiKeyVO[]>([])
const loading = ref(false)
const createModalVisible = ref(false)
const creating = ref(false)
const createForm = ref<API.ApiKeyCreateRequest>({})
const newApiKey = ref<API.ApiKeyVO>()
const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
})

// 加载 API Keys
const loadApiKeys = async (pageNum?: number, pageSize?: number) => {
  loading.value = true
  try {
    const res = await listMyApiKeys({
      pageNum: pageNum || pagination.value.current,
      pageSize: pageSize || pagination.value.pageSize,
    })
    
    if (res.data.code === 0 && res.data.data) {
      apiKeys.value = res.data.data.records || []
      pagination.value.total = res.data.data.totalRow || 0
    }
  } catch (error: any) {
    message.error('加载 API Keys 失败：' + error.message)
  } finally {
    loading.value = false
  }
}

// 处理分页变化
const handleTableChange = (pag: any) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
  loadApiKeys(pag.current, pag.pageSize)
}

// 显示创建弹窗
const showCreateModal = () => {
  createForm.value = {}
  newApiKey.value = undefined
  createModalVisible.value = true
}

// 创建 API Key
const handleCreate = async () => {
  creating.value = true
  try {
    const res = await createApiKey(createForm.value)
    if (res.data.code === 0 && res.data.data) {
      message.success('创建成功')
      newApiKey.value = res.data.data
      // 刷新列表
      await loadApiKeys()
    } else {
      message.error(res.data.message || '创建失败')
    }
  } catch (error: any) {
    message.error('创建失败：' + error.message)
  } finally {
    creating.value = false
  }
}

// 撤销 API Key
const handleRevoke = (record: API.ApiKeyVO) => {
  Modal.confirm({
    title: '确认撤销',
    content: `确定要撤销 API Key "${record.keyName || record.keyValue}" 吗？撤销后将无法恢复。`,
    onOk: async () => {
      try {
        const res = await revokeApiKey({ id: record.id })
        if (res.data.code === 0) {
          message.success('撤销成功')
          await loadApiKeys()
        } else {
          message.error(res.data.message || '撤销失败')
        }
      } catch (error: any) {
        message.error('撤销失败：' + error.message)
      }
    },
  })
}

// 复制 API Key
const copyApiKey = () => {
  if (newApiKey.value?.keyValue) {
    navigator.clipboard.writeText(newApiKey.value.keyValue)
    message.success('已复制到剪贴板')
  }
}

onMounted(() => {
  loadApiKeys()
})
</script>

<style scoped>
.api-key-page {
  padding: 24px;
}
</style>
