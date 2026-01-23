<template>
  <div id="userManagePage">
    <!-- 搜索表单 -->
    <a-form layout="inline" :model="searchParams" @finish="doSearch">
      <a-form-item label="账号">
        <a-input v-model:value="searchParams.userAccount" placeholder="输入账号" />
      </a-form-item>
      <a-form-item label="用户名">
        <a-input v-model:value="searchParams.userName" placeholder="输入用户名" />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit">搜索</a-button>
      </a-form-item>
    </a-form>
    <a-divider />
    <!-- 表格 -->
    <a-table
      :columns="columns"
      :data-source="data"
      :pagination="pagination"
      @change="doTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'userAvatar'">
          <a-image :src="record.userAvatar" :width="60" />
        </template>
        <template v-else-if="column.dataIndex === 'userRole'">
          <div v-if="record.userRole === 'admin'">
            <a-tag color="green">管理员</a-tag>
          </div>
          <div v-else>
            <a-tag color="blue">普通用户</a-tag>
          </div>
        </template>
        <template v-else-if="column.dataIndex === 'userStatus'">
          <a-tag v-if="record.userStatus === 'active'" color="green">正常</a-tag>
          <a-tag v-else color="red">禁用</a-tag>
        </template>
        <template v-else-if="column.dataIndex === 'tokenQuota'">
          <span v-if="record.tokenQuota === -1">无限制</span>
          <span v-else>{{ record.tokenQuota?.toLocaleString() }}</span>
        </template>
        <template v-else-if="column.dataIndex === 'usedTokens'">
          {{ (record.usedTokens || 0).toLocaleString() }}
        </template>
        <template v-else-if="column.dataIndex === 'createTime'">
          {{ dayjs(record.createTime).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
        <template v-else-if="column.key === 'action'">
          <a-space>
            <a-button type="link" size="small" @click="showAnalysis(record)">分析</a-button>
            <a-button type="link" size="small" @click="showQuotaModal(record)"
              >配额</a-button
            >
            <a-button
              v-if="record.userStatus === 'active'"
              type="link"
              size="small"
              danger
              @click="doDisable(record.id)"
              >禁用</a-button
            >
            <a-button
              v-else
              type="link"
              size="small"
              @click="doEnable(record.id)"
              >启用</a-button
            >
            <a-button danger size="small" @click="doDelete(record.id)">删除</a-button>
          </a-space>
        </template>
      </template>
    </a-table>

    <!-- 配额管理弹窗 -->
    <a-modal
      v-model:open="quotaModalVisible"
      title="配额管理"
      @ok="handleQuotaSubmit"
      @cancel="quotaModalVisible = false"
    >
      <a-form :model="quotaForm" layout="vertical">
        <a-form-item label="用户">
          <span>{{ currentUser?.userName }} ({{ currentUser?.userAccount }})</span>
        </a-form-item>
        <a-form-item label="当前已使用">
          <span>{{ (currentUser?.usedTokens || 0).toLocaleString() }} Tokens</span>
        </a-form-item>
        <a-form-item label="Token配额">
          <a-input-number
            v-model:value="quotaForm.tokenQuota"
            :min="-1"
            :step="10000"
            style="width: 100%"
            placeholder="输入配额（-1表示无限制）"
          />
          <div style="margin-top: 8px; color: #999; font-size: 12px">
            提示：-1 表示无限制，0 表示不能使用
          </div>
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button @click="quotaModalVisible = false">取消</a-button>
        <a-button type="primary" danger @click="handleResetQuota">重置已使用</a-button>
        <a-button type="primary" @click="handleQuotaSubmit">保存</a-button>
      </template>
    </a-modal>

    <!-- 用户分析弹窗 -->
    <a-modal
      v-model:open="analysisModalVisible"
      title="用户使用分析"
      :footer="null"
      width="700px"
    >
      <a-descriptions bordered :column="2" v-if="analysisData">
        <a-descriptions-item label="用户账号">{{
          analysisData.userAccount
        }}</a-descriptions-item>
        <a-descriptions-item label="用户名">{{
          analysisData.userName
        }}</a-descriptions-item>
        <a-descriptions-item label="状态">
          <a-tag v-if="analysisData.userStatus === 'active'" color="green">正常</a-tag>
          <a-tag v-else color="red">禁用</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="角色">
          <a-tag v-if="analysisData.userRole === 'admin'" color="green">管理员</a-tag>
          <a-tag v-else color="blue">普通用户</a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="Token配额">
          <span v-if="analysisData.tokenQuota === -1">无限制</span>
          <span v-else>{{ analysisData.tokenQuota?.toLocaleString() }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="已使用">
          {{ (analysisData.usedTokens || 0).toLocaleString() }} Tokens
        </a-descriptions-item>
        <a-descriptions-item label="剩余配额">
          <span v-if="analysisData.remainingQuota === -1">无限制</span>
          <span v-else>{{ analysisData.remainingQuota?.toLocaleString() }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="总请求数">
          {{ analysisData.totalRequests }}
        </a-descriptions-item>
        <a-descriptions-item label="成功请求数">
          {{ analysisData.successRequests }}
        </a-descriptions-item>
        <a-descriptions-item label="累计Token">
          {{ (analysisData.totalTokens || 0).toLocaleString() }} Tokens
        </a-descriptions-item>
        <a-descriptions-item label="累计费用">
          ¥{{ (analysisData.totalCost || 0).toFixed(6) }}
        </a-descriptions-item>
        <a-descriptions-item label="今日费用">
          ¥{{ (analysisData.todayCost || 0).toFixed(6) }}
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>
<script lang="ts" setup>
import { computed, onMounted, reactive, ref } from 'vue'
import {
  deleteUser,
  listUserVoByPage,
  disableUser,
  enableUser,
  setUserQuota,
  resetUserQuota,
  getUserAnalysis,
} from '@/api/userController.ts'
import { message, Modal } from 'ant-design-vue'
import dayjs from 'dayjs'

const columns = [
  {
    title: 'id',
    dataIndex: 'id',
    width: 80,
  },
  {
    title: '账号',
    dataIndex: 'userAccount',
    width: 120,
  },
  {
    title: '用户名',
    dataIndex: 'userName',
    width: 120,
  },
  {
    title: '头像',
    dataIndex: 'userAvatar',
    width: 150,
  },
  {
    title: '简介',
    dataIndex: 'userProfile',
    width: 150,
  },
  {
    title: '用户角色',
    dataIndex: 'userRole',
    width: 100,
  },
  {
    title: '状态',
    dataIndex: 'userStatus',
    width: 100,
  },
  {
    title: 'Token配额',
    dataIndex: 'tokenQuota',
    width: 120,
  },
  {
    title: '已使用',
    dataIndex: 'usedTokens',
    width: 100,
  },
  {
    title: '创建时间',
    dataIndex: 'createTime',
    width: 180,
  },
  {
    title: '操作',
    key: 'action',
    width: 280,
    fixed: 'right',
  },
]

// 展示的数据
const data = ref<API.UserVO[]>([])
const total = ref(0)

// 搜索条件
const searchParams = reactive<API.UserQueryRequest>({
  pageNum: 1,
  pageSize: 10,
})

// 获取数据
const fetchData = async () => {
  const res = await listUserVoByPage({
    ...searchParams,
  })
  if (res.data.data) {
    data.value = res.data.data.records ?? []
    total.value = res.data.data.totalRow ?? 0
  } else {
    message.error('获取数据失败，' + res.data.message)
  }
}

// 分页参数
const pagination = computed(() => {
  return {
    current: searchParams.pageNum ?? 1,
    pageSize: searchParams.pageSize ?? 10,
    total: total.value,
    showSizeChanger: true,
    showTotal: (total: number) => `共 ${total} 条`,
  }
})

// 表格分页变化时的操作
const doTableChange = (page: { current: number; pageSize: number }) => {
  searchParams.pageNum = page.current
  searchParams.pageSize = page.pageSize
  fetchData()
}

// 搜索数据
const doSearch = () => {
  // 重置页码
  searchParams.pageNum = 1
  fetchData()
}

// 删除数据
const doDelete = async (id: number) => {
  if (!id) {
    return
  }
  const res = await deleteUser({ id })
  if (res.data.code === 0) {
    message.success('删除成功')
    // 刷新数据
    fetchData()
  } else {
    message.error('删除失败')
  }
}

// 禁用用户
const doDisable = async (userId: number) => {
  Modal.confirm({
    title: '确认禁用',
    content: '确定要禁用该用户吗？禁用后用户将无法登录和使用服务。',
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      const res = await disableUser({ userId })
      if (res.data.code === 0) {
        message.success('禁用成功')
        fetchData()
      } else {
        message.error('禁用失败：' + res.data.message)
      }
    },
  })
}

// 启用用户
const doEnable = async (userId: number) => {
  const res = await enableUser({ userId })
  if (res.data.code === 0) {
    message.success('启用成功')
    fetchData()
  } else {
    message.error('启用失败：' + res.data.message)
  }
}

// 配额管理弹窗
const quotaModalVisible = ref(false)
const currentUser = ref<API.UserVO>()
const quotaForm = reactive({
  tokenQuota: 0,
})

const showQuotaModal = (record: API.UserVO) => {
  currentUser.value = record
  quotaForm.tokenQuota = record.tokenQuota || 0
  quotaModalVisible.value = true
}

const handleQuotaSubmit = async () => {
  if (!currentUser.value) {
    return
  }

  const res = await setUserQuota({
    userId: currentUser.value.id,
    tokenQuota: quotaForm.tokenQuota,
  })

  if (res.data.code === 0) {
    message.success('配额设置成功')
    quotaModalVisible.value = false
    fetchData()
  } else {
    message.error('配额设置失败：' + res.data.message)
  }
}

const handleResetQuota = async () => {
  if (!currentUser.value) {
    return
  }

  Modal.confirm({
    title: '确认重置',
    content: '确定要重置该用户的已使用配额吗？',
    okText: '确认',
    cancelText: '取消',
    onOk: async () => {
      const res = await resetUserQuota({ userId: currentUser.value!.id! })
      if (res.data.code === 0) {
        message.success('重置成功')
        quotaModalVisible.value = false
        fetchData()
      } else {
        message.error('重置失败：' + res.data.message)
      }
    },
  })
}

// 用户分析弹窗
const analysisModalVisible = ref(false)
const analysisData = ref<API.UserAnalysisVO>()

const showAnalysis = async (record: API.UserVO) => {
  const res = await getUserAnalysis({ userId: record.id! })
  if (res.data.code === 0) {
    analysisData.value = res.data.data
    analysisModalVisible.value = true
  } else {
    message.error('获取分析数据失败：' + res.data.message)
  }
}

// 页面加载时请求一次
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
#userManagePage {
  padding: 24px;
  background: white;
  margin-top: 16px;
}
</style>
