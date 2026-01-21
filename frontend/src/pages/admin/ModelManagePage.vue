<template>
  <div id="modelManagePage">
    <h2>模型管理</h2>

    <!-- 搜索表单 -->
    <a-form layout="inline" :model="searchParams" @finish="doSearch">
      <a-form-item label="模型名称">
        <a-input v-model:value="searchParams.modelName" placeholder="输入模型名称" />
      </a-form-item>
      <a-form-item label="模型类型">
        <a-select v-model:value="searchParams.modelType" placeholder="选择类型" style="width: 150px" allow-clear>
          <a-select-option value="chat">对话模型</a-select-option>
          <a-select-option value="embedding">向量模型</a-select-option>
          <a-select-option value="image">图像模型</a-select-option>
        </a-select>
      </a-form-item>
      <a-form-item label="状态">
        <a-select v-model:value="searchParams.status" placeholder="选择状态" style="width: 120px" allow-clear>
          <a-select-option value="active">启用</a-select-option>
          <a-select-option value="inactive">禁用</a-select-option>
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
      添加模型
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
        <template v-if="column.dataIndex === 'modelType'">
          <a-tag v-if="record.modelType === 'chat'" color="blue">对话模型</a-tag>
          <a-tag v-else-if="record.modelType === 'embedding'" color="green">向量模型</a-tag>
          <a-tag v-else-if="record.modelType === 'image'" color="purple">图像模型</a-tag>
          <a-tag v-else>{{ record.modelType }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'status'">
          <a-tag v-if="record.status === 'active'" color="success">启用</a-tag>
          <a-tag v-else-if="record.status === 'inactive'" color="default">禁用</a-tag>
          <a-tag v-else color="warning">{{ record.status }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'healthStatus'">
          <a-tag v-if="record.healthStatus === 'healthy'" color="success">健康</a-tag>
          <a-tag v-else-if="record.healthStatus === 'unhealthy'" color="default">不健康</a-tag>
          <a-tag v-else-if="record.healthStatus === 'degraded'" color="warning">降级</a-tag>
          <a-tag v-else color="error">{{ record.healthStatus }}</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'inputPrice'">
          ¥{{ record.inputPrice }}/千Token
        </template>
        <template v-else-if="column.dataIndex === 'outputPrice'">
          ¥{{ record.outputPrice }}/千Token
        </template>
        <template v-else-if="column.dataIndex === 'avgLatency'">
          {{ record.avgLatency }}/ms
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
      :title="isEdit ? '编辑模型' : '添加模型'"
      @ok="handleSubmit"
      @cancel="handleCancel"
      width="600px"
    >
      <a-form
        :model="formData"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="提供者" required>
          <a-select v-model:value="formData.providerId" placeholder="选择提供者" :disabled="isEdit">
            <a-select-option v-for="provider in providers" :key="provider.id" :value="provider.id">
              {{ provider.displayName }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="模型标识" required>
          <a-input v-model:value="formData.modelKey" placeholder="如: qwen-plus" :disabled="isEdit" />
        </a-form-item>
        <a-form-item label="模型名称" required>
          <a-input v-model:value="formData.modelName" placeholder="如: Qwen Plus" />
        </a-form-item>
        <a-form-item label="模型类型" required>
          <a-select v-model:value="formData.modelType" placeholder="选择类型" :disabled="isEdit">
            <a-select-option value="chat">对话模型</a-select-option>
            <a-select-option value="embedding">向量模型</a-select-option>
            <a-select-option value="image">图像模型</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="formData.description" placeholder="模型描述" :rows="3" />
        </a-form-item>
        <a-form-item label="上下文长度">
          <a-input-number v-model:value="formData.contextLength" :min="1024" :max="200000" style="width: 100%" />
        </a-form-item>
        <a-form-item label="输入价格(元/千Token)">
          <a-input-number v-model:value="formData.inputPrice" :min="0" :step="0.001" :precision="6" style="width: 100%" />
        </a-form-item>
        <a-form-item label="输出价格(元/千Token)">
          <a-input-number v-model:value="formData.outputPrice" :min="0" :step="0.001" :precision="6" style="width: 100%" />
        </a-form-item>
        <a-form-item label="优先级">
          <a-input-number v-model:value="formData.priority" :min="0" :max="1000" style="width: 100%" />
        </a-form-item>
        <a-form-item label="超时时间(毫秒)">
          <a-input-number v-model:value="formData.defaultTimeout" :min="1000" :max="300000" style="width: 100%" />
        </a-form-item>
        <a-form-item label="状态" v-if="isEdit">
          <a-select v-model:value="formData.status" placeholder="选择状态">
            <a-select-option value="active">启用</a-select-option>
            <a-select-option value="inactive">禁用</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { addModel, deleteModel, listModelVoByPage, updateModel } from '@/api/modelController'
import { listProviderVo } from '@/api/modelProviderController'
import { message } from 'ant-design-vue'

// 表格列定义
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: 80,
  },
  {
    title: '模型标识',
    dataIndex: 'modelKey',
    width: 150,
  },
  {
    title: '模型名称',
    dataIndex: 'modelName',
    width: 150,
  },
  {
    title: '提供者',
    dataIndex: 'providerDisplayName',
    width: 120,
  },
  {
    title: '类型',
    dataIndex: 'modelType',
    width: 100,
  },
  {
    title: '输入价格',
    dataIndex: 'inputPrice',
    width: 130,
  },
  {
    title: '输出价格',
    dataIndex: 'outputPrice',
    width: 130,
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
    title: '健康状态',
    dataIndex: 'healthStatus',
    width: 80,
  },
  {
    title: '平均延迟',
    dataIndex: 'avgLatency',
    width: 80,
  },
  {
    title: '成功率',
    dataIndex: 'successRate',
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
const providers = ref<any[]>([])

// 搜索参数
const searchParams = reactive({
  modelName: '',
  modelType: undefined,
  status: undefined,
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
  providerId: undefined,
  modelKey: '',
  modelName: '',
  modelType: 'chat',
  description: '',
  contextLength: 4096,
  inputPrice: 0,
  outputPrice: 0,
  priority: 100,
  defaultTimeout: 60000,
  status: 'active',
})

// 加载数据
const loadData = async () => {
  loading.value = true
  try {
    const res = await listModelVoByPage(searchParams)
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

// 加载提供者列表
const loadProviders = async () => {
  try {
    const res = await listProviderVo()
    if (res.data.data) {
      providers.value = res.data.data || []
    }
  } catch (error) {
    console.error('加载提供者列表失败', error)
  }
}

// 搜索
const doSearch = () => {
  searchParams.pageNum = 1
  loadData()
}

// 重置搜索
const resetSearch = () => {
  searchParams.modelName = ''
  searchParams.modelType = undefined
  searchParams.status = undefined
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
    providerId: record.providerId,
    modelKey: record.modelKey,
    modelName: record.modelName,
    modelType: record.modelType,
    description: record.description,
    contextLength: record.contextLength,
    inputPrice: record.inputPrice,
    outputPrice: record.outputPrice,
    priority: record.priority,
    defaultTimeout: record.defaultTimeout,
    status: record.status,
  })
  modalVisible.value = true
}

// 重置表单
const resetFormData = () => {
  Object.assign(formData, {
    id: undefined,
    providerId: undefined,
    modelKey: '',
    modelName: '',
    modelType: 'chat',
    description: '',
    contextLength: 4096,
    inputPrice: 0,
    outputPrice: 0,
    priority: 100,
    defaultTimeout: 60000,
    status: 'active',
  })
}

// 提交
const handleSubmit = async () => {
  try {
    const res = isEdit.value
      ? await updateModel(formData)
      : await addModel(formData)

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
    const res = await deleteModel({ id })
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
  loadProviders()
})
</script>

<style scoped>
#modelManagePage {
  padding: 20px;
}
</style>
