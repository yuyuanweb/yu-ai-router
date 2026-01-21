<template>
  <div id="providerManagePage">
    <h2>模型提供者管理</h2>
    
    <!-- 搜索表单 -->
    <a-form layout="inline" :model="searchParams" @finish="doSearch">
      <a-form-item label="提供者名称">
        <a-input v-model:value="searchParams.displayName" placeholder="输入提供者名称" />
      </a-form-item>
      <a-form-item label="健康状态">
        <a-select v-model:value="searchParams.healthStatus" placeholder="选择健康状态" style="width: 150px" allow-clear>
          <a-select-option value="healthy">健康</a-select-option>
          <a-select-option value="degraded">降级</a-select-option>
          <a-select-option value="unhealthy">不健康</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit">搜索</a-button>
        <a-button style="margin-left: 10px" @click="resetSearch">重置</a-button>
      </a-form-item>
    </a-form>
    
    <a-divider />
    
    <!-- 添加按钮 -->
    <a-button type="primary" @click="showAddModal" style="margin-bottom: 16px">
      添加提供者
    </a-button>
    
    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data-source="data"
      :pagination="pagination"
      :loading="loading"
      @change="doTableChange"
      row-key="id"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'healthStatus'">
          <a-tag v-if="record.healthStatus === 'healthy'" color="success">健康</a-tag>
          <a-tag v-else-if="record.healthStatus === 'degraded'" color="warning">降级</a-tag>
          <a-tag v-else-if="record.healthStatus === 'unhealthy'" color="error">不健康</a-tag>
          <a-tag v-else color="default">未知</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'status'">
          <a-tag v-if="record.status === 'active'" color="success">启用</a-tag>
          <a-tag v-else-if="record.status === 'inactive'" color="default">禁用</a-tag>
          <a-tag v-else color="warning">{{ record.status }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'avgLatency'">
          {{ record.avgLatency }}ms
        </template>
        <template v-else-if="column.dataIndex === 'successRate'">
          {{ record.successRate }}%
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="showEditModal(record)">编辑</a-button>
            <a-button type="link" size="small" danger @click="doDelete(record.id)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>
    
    <!-- 添加/编辑模态框 -->
    <a-modal
      v-model:open="modalVisible"
      :title="isEdit ? '编辑提供者' : '添加提供者'"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="600px"
    >
      <a-form
        :model="formData"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="提供者标识" required>
          <a-input v-model:value="formData.providerName" placeholder="如: qwen" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="显示名称" required>
          <a-input v-model:value="formData.displayName" placeholder="如: 通义千问" />
        </a-form-item>
        <a-form-item label="API基础URL" required>
          <a-input v-model:value="formData.baseUrl" placeholder="如: https://dashscope.aliyuncs.com/compatible-mode/v1" />
        </a-form-item>
        <a-form-item label="API密钥" required>
          <a-input-password v-model:value="formData.apiKey" placeholder="输入API密钥" />
        </a-form-item>
        <a-form-item label="优先级">
          <a-input-number v-model:value="formData.priority" :min="0" :max="1000" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态" v-if="isEdit">
          <a-select v-model:value="formData.status" placeholder="选择状态">
            <a-select-option value="active">启用</a-select-option>
            <a-select-option value="inactive">禁用</a-select-option>
            <a-select-option value="maintenance">维护中</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { addProvider, deleteProvider, listProviderVoByPage, updateProvider } from '@/api/modelProviderController'
import { message } from 'ant-design-vue'

// 表格列定义
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: 80,
  },
  {
    title: '提供者标识',
    dataIndex: 'providerName',
    width: 120,
  },
  {
    title: '显示名称',
    dataIndex: 'displayName',
    width: 150,
  },
  {
    title: '健康状态',
    dataIndex: 'healthStatus',
    width: 100,
  },
  {
    title: '平均延迟',
    dataIndex: 'avgLatency',
    width: 100,
  },
  {
    title: '成功率',
    dataIndex: 'successRate',
    width: 100,
  },
  {
    title: '优先级',
    dataIndex: 'priority',
    width: 80,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
  },
  {
    title: '操作',
    key: 'action',
    width: 150,
    fixed: 'right',
  },
]

const data = ref<any[]>([])
const loading = ref(false)

// 搜索参数
const searchParams = reactive({
  displayName: '',
  healthStatus: undefined,
  pageNum: 1,
  pageSize: 10,
})

// 分页配置
const pagination = computed(() => ({
  current: searchParams.pageNum,
  pageSize: searchParams.pageSize,
  total: total.value,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}))

const total = ref(0)

// 模态框相关
const modalVisible = ref(false)
const isEdit = ref(false)
const formData = reactive({
  id: undefined,
  providerName: '',
  displayName: '',
  baseUrl: '',
  apiKey: '',
  priority: 100,
  status: 'active',
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await listProviderVoByPage(searchParams)
    if (res.data.data) {
      data.value = res.data.data.records || []
      total.value = res.data.data.totalRow || 0
    } else {
      message.error('加载失败：' + res.data.message)
    }
  } catch (error) {
    message.error('加载失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const doSearch = () => {
  searchParams.pageNum = 1
  loadData()
}

// 重置搜索
const resetSearch = () => {
  searchParams.displayName = ''
  searchParams.healthStatus = undefined
  searchParams.pageNum = 1
  loadData()
}

// 表格变化
const doTableChange = (pag: any) => {
  searchParams.pageNum = pag.current
  searchParams.pageSize = pag.pageSize
  loadData()
}

// 显示添加模态框
const showAddModal = () => {
  isEdit.value = false
  resetFormData()
  modalVisible.value = true
}

// 显示编辑模态框
const showEditModal = (record: any) => {
  isEdit.value = true
  Object.assign(formData, {
    id: record.id,
    providerName: record.providerName,
    displayName: record.displayName,
    baseUrl: record.baseUrl,
    apiKey: '',
    priority: record.priority,
    status: record.status,
  })
  modalVisible.value = true
}

// 重置表单
const resetFormData = () => {
  Object.assign(formData, {
    id: undefined,
    providerName: '',
    displayName: '',
    baseUrl: '',
    apiKey: '',
    priority: 100,
    status: 'active',
  })
}

// 提交
const handleSubmit = async () => {
  try {
    const res = isEdit.value
      ? await updateProvider(formData)
      : await addProvider(formData)
    
    if (res.data.code === 0) {
      message.success(isEdit.value ? '更新成功' : '添加成功')
      modalVisible.value = false
      loadData()
    } else {
      message.error(res.data.message || '操作失败')
    }
  } catch (error) {
    message.error('操作失败')
  }
}

// 取消
const handleCancel = () => {
  modalVisible.value = false
  resetFormData()
}

// 删除
const doDelete = async (id: number) => {
  try {
    const res = await deleteProvider({ id })
    if (res.data.code === 0) {
      message.success('删除成功')
      loadData()
    } else {
      message.error('删除失败：' + res.data.message)
    }
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
#providerManagePage {
  padding: 20px;
}
</style>
