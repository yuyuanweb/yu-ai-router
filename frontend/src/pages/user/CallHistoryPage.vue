<template>
  <div id="callHistoryPage" class="morandi-page">
    <!-- 搜索表单区域 -->
    <div class="search-card">
      <a-form layout="inline" :model="searchParams" @finish="doSearch" class="morandi-form">
        <div class="form-fields">
          <a-form-item label="模型">
            <a-input v-model:value="searchParams.requestModel" placeholder="输入模型标识" allow-clear />
          </a-form-item>
          <a-form-item label="状态">
            <a-select v-model:value="searchParams.status" placeholder="选择状态" allow-clear>
              <a-select-option value="success">成功</a-select-option>
              <a-select-option value="failed">失败</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="来源">
            <a-select v-model:value="searchParams.source" placeholder="选择来源" allow-clear>
              <a-select-option value="web">网页</a-select-option>
              <a-select-option value="api">API</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="时间">
            <a-space>
              <a-date-picker v-model:value="startDate" placeholder="开始日期" @change="handleDateChange" />
              <span class="date-separator">~</span>
              <a-date-picker v-model:value="endDate" placeholder="结束日期" @change="handleDateChange" />
            </a-space>
          </a-form-item>
        </div>
        <div class="form-actions">
          <a-button type="primary" html-type="submit" class="capsule-btn primary-btn">搜索</a-button>
          <a-button @click="resetSearch" class="capsule-btn reset-btn">重置</a-button>
        </div>
      </a-form>
    </div>

    <!-- 数据表格区域 -->
    <div class="table-card">
      <a-table
        :columns="columns"
        :data-source="data"
        :pagination="pagination"
        :loading="loading"
        @change="doTableChange"
        class="morandi-table"
        :row-key="(record: any) => record.traceId"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'status'">
            <span class="status-dot" :class="record.status"></span>
            <span :class="record.status === 'success' ? 'text-success' : 'text-error'">
              {{ record.status === 'success' ? '成功' : '失败' }}
            </span>
          </template>
          <template v-else-if="column.dataIndex === 'source'">
            <a-tag class="soft-tag" :color="record.source === 'web' ? 'blue' : 'purple'">
              {{ record.source === 'web' ? '网页' : 'API' }}
            </a-tag>
          </template>
          <template v-else-if="column.dataIndex === 'cost'">
            <span class="cost-text">¥{{ (record.cost || 0).toFixed(6) }}</span>
          </template>
          <template v-else-if="column.dataIndex === 'duration'">
            <span class="duration-text">{{ record.duration }} ms</span>
          </template>
          <template v-else-if="column.dataIndex === 'routingStrategy'">
            <a-tag class="soft-tag" :color="record.routingStrategy === 'auto' ? 'cyan' : 'default'">
              {{ record.routingStrategy }}
            </a-tag>
          </template>
          <template v-else-if="column.dataIndex === 'createTime'">
            <span class="time-text">{{ dayjs(record.createTime).format('YYYY-MM-DD HH:mm:ss') }}</span>
          </template>
          <template v-else-if="column.key === 'action'">
            <a-button type="link" class="action-link" @click="showDetail(record)">详情</a-button>
          </template>
        </template>
      </a-table>
    </div>

    <!-- 详情弹窗 -->
    <a-modal
      v-model:open="detailModalVisible"
      title="调用详情"
      :footer="null"
      width="800px"
      class="morandi-modal"
    >
      <a-descriptions bordered :column="2" v-if="currentRecord" class="morandi-descriptions">
        <!-- 详情内容保持原样 -->
        <a-descriptions-item label="TraceId" :span="2">{{ currentRecord.traceId }}</a-descriptions-item>
        <a-descriptions-item label="模型">{{ currentRecord.requestModel }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag v-if="currentRecord.status === 'success'" color="success">成功</a-tag>
          <a-tag v-else color="error">失败</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="来源">{{ currentRecord.source }}</a-descriptions-item>
        <a-descriptions-item label="路由策略">{{ currentRecord.routingStrategy }}</a-descriptions-item>
        <a-descriptions-item label="输入Token">{{ currentRecord.promptTokens }}</a-descriptions-item>
        <a-descriptions-item label="输出Token">{{ currentRecord.completionTokens }}</a-descriptions-item>
        <a-descriptions-item label="总Token">{{ currentRecord.totalTokens }}</a-descriptions-item>
        <a-descriptions-item label="费用">¥{{ (currentRecord.cost || 0).toFixed(6) }}</a-descriptions-item>
        <a-descriptions-item label="耗时">{{ currentRecord.duration }} ms</a-descriptions-item>
        <a-descriptions-item label="Fallback">{{ currentRecord.isFallback ? '是' : '否' }}</a-descriptions-item>
        <a-descriptions-item label="客户端IP" :span="2">{{ currentRecord.clientIp }}</a-descriptions-item>
        <a-descriptions-item label="User-Agent" :span="2">{{ currentRecord.userAgent }}</a-descriptions-item>
        <a-descriptions-item label="创建时间" :span="2">{{ dayjs(currentRecord.createTime).format('YYYY-MM-DD HH:mm:ss') }}</a-descriptions-item>
        <a-descriptions-item v-if="currentRecord.errorMessage" label="错误信息" :span="2">
          <a-typography-text type="danger">{{ currentRecord.errorMessage }}</a-typography-text>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { pageMyHistory } from '@/api/statsController'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'

const columns = [
  { title: 'TraceId', dataIndex: 'traceId', width: 150, ellipsis: true },
  { title: '模型', dataIndex: 'requestModel', width: 120 },
  { title: '状态', dataIndex: 'status', width: 100 },
  { title: '来源', dataIndex: 'source', width: 100 },
  { title: 'Token', dataIndex: 'totalTokens', width: 100 },
  { title: '费用', dataIndex: 'cost', width: 120 },
  { title: '耗时', dataIndex: 'duration', width: 100 },
  { title: '路由策略', dataIndex: 'routingStrategy', width: 110 },
  { title: '创建时间', dataIndex: 'createTime', width: 180 },
  { title: '操作', key: 'action', width: 80, fixed: 'right' },
]

const data = ref<API.RequestLog[]>([])
const total = ref(0)
const loading = ref(false)
const detailModalVisible = ref(false)
const currentRecord = ref<API.RequestLog>()
const startDate = ref()
const endDate = ref()

const searchParams = reactive<API.RequestLogQueryRequest>({
  pageNum: 1,
  pageSize: 10,
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await pageMyHistory({ ...searchParams })
    if (res.data.code === 0 && res.data.data) {
      data.value = res.data.data.records ?? []
      total.value = res.data.data.totalRow ?? 0
    } else {
      message.error('获取数据失败：' + res.data.message)
    }
  } catch (error) {
    message.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

const pagination = computed(() => ({
  current: searchParams.pageNum ?? 1,
  pageSize: searchParams.pageSize ?? 10,
  total: total.value,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
}))

const doTableChange = (page: { current: number; pageSize: number }) => {
  searchParams.pageNum = page.current
  searchParams.pageSize = page.pageSize
  fetchData()
}

const doSearch = () => {
  searchParams.pageNum = 1
  fetchData()
}

const resetSearch = () => {
  searchParams.requestModel = undefined
  searchParams.status = undefined
  searchParams.source = undefined
  searchParams.startDate = undefined
  searchParams.endDate = undefined
  startDate.value = undefined
  endDate.value = undefined
  searchParams.pageNum = 1
  fetchData()
}

const handleDateChange = () => {
  if (startDate.value) searchParams.startDate = dayjs(startDate.value).format('YYYY-MM-DD')
  if (endDate.value) searchParams.endDate = dayjs(endDate.value).format('YYYY-MM-DD')
}

const showDetail = (record: API.RequestLog) => {
  currentRecord.value = record
  detailModalVisible.value = true
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.morandi-page {
  padding: 24px;
  background: var(--bg-page);
  min-height: calc(100vh - 72px);
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 搜索卡片 */
.search-card {
  background: var(--bg-card);
  padding: 24px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-soft);
}

.morandi-form {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.form-fields {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  flex: 1;
}

/* 覆盖表单基础样式 */
:deep(.ant-form-item) {
  margin-bottom: 0;
}

:deep(.ant-input), :deep(.ant-select-selector), :deep(.ant-picker) {
  border-radius: 8px !important;
  border-color: var(--border-color) !important;
  background: rgba(250, 247, 244, 0.6) !important;
}
:deep(.ant-input:hover), :deep(.ant-select-selector:hover), :deep(.ant-picker:hover) {
  border-color: var(--color-primary) !important;
}

.date-separator {
  color: var(--text-secondary);
}

.form-actions {
  display: flex;
  gap: 12px;
}

/* 统一胶囊按钮 */
.capsule-btn {
  border-radius: 20px;
  padding: 0 24px;
  height: 36px;
  transition: all 0.3s ease;
}
.primary-btn {
  background: var(--color-primary);
  border: none;
  box-shadow: 0 4px 12px rgba(184, 169, 184, 0.2);
}
.primary-btn:hover {
  background: #C5B8C5;
  transform: translateY(-1px);
}
.reset-btn {
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  background: transparent;
}
.reset-btn:hover {
  color: var(--color-primary);
  border-color: var(--color-primary);
}

/* 表格卡片 */
.table-card {
  background: var(--bg-card);
  padding: 24px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-soft);
  flex: 1;
  overflow: hidden;
}

/* 深度美化表格 */
:deep(.morandi-table .ant-table) {
  background: transparent;
}
:deep(.morandi-table .ant-table-thead > tr > th) {
  background: transparent !important;
  color: var(--text-main) !important;
  font-weight: 600;
  border-bottom: 1px dashed var(--border-color) !important;
}
:deep(.morandi-table .ant-table-tbody > tr > td) {
  border-bottom: 1px dashed var(--border-color) !important;
  color: var(--text-secondary);
  transition: all 0.3s;
}
:deep(.morandi-table .ant-table-tbody > tr:hover > td) {
  background: rgba(184, 169, 184, 0.05) !important;
}

/* 状态小圆点 */
.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 6px;
}
.status-dot.success { background-color: #85C79B; }
.status-dot.failed { background-color: #D89696; }
.text-success { color: #6C9A7C; font-weight: 500; }
.text-error { color: #C48383; font-weight: 500; }

.cost-text { color: var(--text-main); font-weight: 500; }
.duration-text { color: var(--text-secondary); font-family: monospace; }
.time-text { color: var(--text-secondary); font-size: 13px; }

/* 柔和标签 */
.soft-tag {
  border-radius: 6px;
  border: none;
  background: rgba(184, 169, 184, 0.1);
  color: var(--text-main);
}

.action-link {
  color: var(--color-primary);
  font-weight: 500;
  transition: all 0.3s;
}
.action-link:hover {
  color: #8B7896;
}

/* 弹窗 */
:deep(.morandi-modal .ant-modal-content) {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  box-shadow: 0 12px 40px rgba(0,0,0,0.08);
}
:deep(.morandi-descriptions .ant-descriptions-item-label) {
  background: rgba(245, 242, 239, 0.5) !important;
  color: var(--text-secondary) !important;
}
:deep(.morandi-descriptions .ant-descriptions-view) {
  border-radius: 12px;
  border-color: var(--border-color);
}
:deep(.morandi-descriptions .ant-descriptions-item) {
  border-color: var(--border-color) !important;
}
</style>
