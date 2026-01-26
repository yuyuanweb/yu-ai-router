<template>
  <div id="myKeysPage">
    <h2>我的密钥（BYOK）</h2>
    <p class="page-desc">
      配置您自己的模型提供者 API Key，使用 BYOK 模式时不消耗平台余额和配额，直接使用您的账号调用模型。
    </p>

    <a-card>
      <a-space style="margin-bottom: 16px">
        <a-button type="primary" @click="showAddModal">
          <template #icon><PlusOutlined /></template>
          添加密钥
        </a-button>
        <a-button @click="loadKeys">
          <template #icon><ReloadOutlined /></template>
          刷新列表
        </a-button>
      </a-space>

      <!-- 密钥列表 -->
      <a-table
        :columns="columns"
        :data-source="keys"
        :loading="loading"
        :pagination="false"
        row-key="id"
      >
        <!-- 提供者名称 -->
        <template #providerName="{ record }">
          <a-tag color="blue">{{ record.providerName }}</a-tag>
        </template>

        <!-- API Key 脱敏显示 -->
        <template #apiKey="{ record }">
          <code class="api-key-display">{{ record.apiKey }}</code>
        </template>

        <!-- 状态 -->
        <template #status="{ record }">
          <a-tag v-if="record.status === 'active'" color="success">启用</a-tag>
          <a-tag v-else color="default">禁用</a-tag>
        </template>

        <!-- 操作 -->
        <template #action="{ record }">
          <a-space>
            <a-button size="small" @click="showEditModal(record)">编辑</a-button>
            <a-popconfirm
              title="确定要删除这个密钥吗？"
              @confirm="handleDelete(record.id)"
            >
              <a-button size="small" danger>删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- 添加密钥弹窗 -->
    <a-modal
      v-model:open="addModalVisible"
      title="添加 API Key"
      @ok="handleAdd"
      @cancel="addModalVisible = false"
      width="600px"
    >
      <a-form :model="addForm" layout="vertical">
        <a-form-item label="模型提供者" required>
          <a-select v-model:value="addForm.providerId" placeholder="请选择模型提供者">
            <a-select-option v-for="provider in providers" :key="provider.id" :value="provider.id">
              {{ provider.displayName }} ({{ provider.providerName }})
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item label="API Key" required>
          <a-input-password 
            v-model:value="addForm.apiKey" 
            placeholder="请输入您自己的 API Key"
            :maxlength="512"
          />
          <div class="form-hint">
            密钥将加密存储，仅您可见。使用您自己的密钥调用模型时，不消耗平台余额。
          </div>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 编辑密钥弹窗 -->
    <a-modal
      v-model:open="editModalVisible"
      title="更新 API Key"
      @ok="handleUpdate"
      @cancel="editModalVisible = false"
      width="600px"
    >
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="模型提供者">
          <a-input :value="currentKey?.providerName" disabled />
        </a-form-item>
        
        <a-form-item label="新的 API Key">
          <a-input-password 
            v-model:value="editForm.apiKey" 
            placeholder="请输入新的 API Key（留空表示不修改）"
            :maxlength="512"
          />
        </a-form-item>
        
        <a-form-item label="状态">
          <a-radio-group v-model:value="editForm.status">
            <a-radio value="active">启用</a-radio>
            <a-radio value="inactive">禁用</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import {
  addUserProviderKey,
  listMyProviderKeys,
  updateUserProviderKey,
  deleteUserProviderKey,
} from '@/api/userProviderKeyController'
import { listProviderVo } from '@/api/modelProviderController'

// 列定义
const columns = [
  {
    title: '提供者',
    dataIndex: 'providerName',
    key: 'providerName',
    slots: { customRender: 'providerName' },
    width: 150,
  },
  {
    title: 'API Key',
    dataIndex: 'apiKey',
    key: 'apiKey',
    slots: { customRender: 'apiKey' },
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    slots: { customRender: 'status' },
    width: 100,
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    key: 'createTime',
    width: 180,
  },
  {
    title: '操作',
    key: 'action',
    slots: { customRender: 'action' },
    width: 180,
  },
]

const keys = ref<API.UserProviderKeyVO[]>([])
const providers = ref<API.ProviderVO[]>([])
const loading = ref(false)
const addModalVisible = ref(false)
const editModalVisible = ref(false)
const currentKey = ref<API.UserProviderKeyVO | null>(null)

// 添加表单
const addForm = ref({
  providerId: undefined as number | undefined,
  apiKey: '',
})

// 编辑表单
const editForm = ref({
  id: undefined as number | undefined,
  apiKey: '',
  status: 'active',
})

// 加载密钥列表
const loadKeys = async () => {
  loading.value = true
  try {
    const res = await listMyProviderKeys()
    if (res.data.code === 0 && res.data.data) {
      keys.value = res.data.data
    }
  } catch (err) {
    console.error('加载密钥列表失败', err)
    message.error('加载密钥列表失败')
  } finally {
    loading.value = false
  }
}

// 加载提供者列表
const loadProviders = async () => {
  try {
    const res = await listProviderVo()
    if (res.data.code === 0 && res.data.data) {
      providers.value = res.data.data
    }
  } catch (err) {
    console.error('加载提供者列表失败', err)
  }
}

// 显示添加弹窗
const showAddModal = () => {
  addForm.value = {
    providerId: undefined,
    apiKey: '',
  }
  addModalVisible.value = true
}

// 添加密钥
const handleAdd = async () => {
  if (!addForm.value.providerId || !addForm.value.apiKey) {
    message.error('请填写完整信息')
    return
  }

  try {
    const res = await addUserProviderKey(addForm.value)
    if (res.data.code === 0) {
      message.success('添加成功')
      addModalVisible.value = false
      loadKeys()
    }
  } catch (err) {
    const errorMsg = err instanceof Error ? err.message : '添加失败'
    message.error(errorMsg)
  }
}

// 显示编辑弹窗
const showEditModal = (record: API.UserProviderKeyVO) => {
  currentKey.value = record
  editForm.value = {
    id: record.id as number,
    apiKey: '',
    status: record.status || 'active',
  }
  editModalVisible.value = true
}

// 更新密钥
const handleUpdate = async () => {
  if (!editForm.value.id) {
    return
  }

  try {
    const res = await updateUserProviderKey(editForm.value)
    if (res.data.code === 0) {
      message.success('更新成功')
      editModalVisible.value = false
      loadKeys()
    }
  } catch (err) {
    const errorMsg = err instanceof Error ? err.message : '更新失败'
    message.error(errorMsg)
  }
}

// 删除密钥
const handleDelete = async (id: number) => {
  try {
    const res = await deleteUserProviderKey({ id })
    if (res.data.code === 0) {
      message.success('删除成功')
      loadKeys()
    }
  } catch {
    message.error('删除失败')
  }
}

onMounted(() => {
  loadKeys()
  loadProviders()
})
</script>

<style scoped>
#myKeysPage {
  padding: 24px;
}

h2 {
  margin-bottom: 8px;
}

.page-desc {
  color: #8c8c8c;
  margin-bottom: 24px;
  line-height: 1.6;
}

.api-key-display {
  font-family: 'Monaco', 'Consolas', monospace;
  background: #f5f5f5;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: #262626;
}

.form-hint {
  font-size: 12px;
  color: #8c8c8c;
  margin-top: 8px;
  line-height: 1.5;
}
</style>
