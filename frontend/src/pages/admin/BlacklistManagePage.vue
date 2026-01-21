<template>
  <div id="blacklistManagePage">
    <!-- 添加 IP 表单 -->
    <a-card title="添加 IP 到黑名单" :bordered="false" style="margin-bottom: 16px">
      <a-form layout="inline" :model="addForm" @finish="doAdd">
        <a-form-item label="IP 地址" :rules="[{ required: true, message: '请输入 IP 地址' }]">
          <a-input
            v-model:value="addForm.ip"
            placeholder="输入 IP 地址，如 192.168.1.1"
            style="width: 280px"
          />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="addLoading">
            <template #icon><PlusOutlined /></template>
            添加到黑名单
          </a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 黑名单列表 -->
    <a-card title="黑名单列表" :bordered="false">
      <template #extra>
        <a-button @click="fetchData" :loading="loading">
          <template #icon><ReloadOutlined /></template>
          刷新
        </a-button>
      </template>

      <a-spin :spinning="loading">
        <a-empty v-if="blacklist.length === 0" description="暂无黑名单 IP" />
        <a-list v-else :data-source="blacklist" :grid="{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 4, xxl: 6 }">
          <template #renderItem="{ item }">
            <a-list-item>
              <a-card hoverable size="small">
                <template #actions>
                  <a-popconfirm
                    title="确定要移除该 IP 吗？"
                    ok-text="确定"
                    cancel-text="取消"
                    @confirm="doRemove(item)"
                  >
                    <a-button type="link" danger size="small">
                      <template #icon><DeleteOutlined /></template>
                      移除
                    </a-button>
                  </a-popconfirm>
                </template>
                <a-card-meta>
                  <template #title>
                    <StopOutlined style="color: #ff4d4f; margin-right: 8px" />
                    {{ item }}
                  </template>
                </a-card-meta>
              </a-card>
            </a-list-item>
          </template>
        </a-list>
      </a-spin>

      <div style="margin-top: 16px; color: #666; font-size: 12px">
        <InfoCircleOutlined style="margin-right: 4px" />
        共 {{ blacklist.length }} 个 IP 被拉黑。黑名单中的 IP 将被禁止访问所有 API 接口。
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { getBlacklist, addToBlacklist, removeFromBlacklist } from '@/api/blacklistController'
import { message } from 'ant-design-vue'
import {
  PlusOutlined,
  ReloadOutlined,
  DeleteOutlined,
  StopOutlined,
  InfoCircleOutlined,
} from '@ant-design/icons-vue'

// 黑名单列表
const blacklist = ref<string[]>([])
const loading = ref(false)
const addLoading = ref(false)

// 添加表单
const addForm = reactive({
  ip: '',
})

// 获取黑名单数据
const fetchData = async () => {
  loading.value = true
  try {
    const res = await getBlacklist()
    if (res.data.code === 0 && res.data.data) {
      blacklist.value = res.data.data
    } else {
      message.error('获取黑名单失败：' + res.data.message)
    }
  } catch (e) {
    message.error('获取黑名单失败')
  } finally {
    loading.value = false
  }
}

// 添加 IP 到黑名单
const doAdd = async () => {
  if (!addForm.ip) {
    message.warning('请输入 IP 地址')
    return
  }

  // 简单的 IP 格式验证
  const ipPattern = /^(\d{1,3}\.){3}\d{1,3}$/
  if (!ipPattern.test(addForm.ip)) {
    message.warning('请输入正确的 IP 地址格式')
    return
  }

  addLoading.value = true
  try {
    const res = await addToBlacklist({ ip: addForm.ip })
    if (res.data.code === 0) {
      message.success('添加成功')
      addForm.ip = ''
      fetchData()
    } else {
      message.error('添加失败：' + res.data.message)
    }
  } catch (e) {
    message.error('添加失败')
  } finally {
    addLoading.value = false
  }
}

// 从黑名单移除 IP
const doRemove = async (ip: string) => {
  try {
    const res = await removeFromBlacklist({ ip })
    if (res.data.code === 0) {
      message.success('移除成功')
      fetchData()
    } else {
      message.error('移除失败：' + res.data.message)
    }
  } catch (e) {
    message.error('移除失败')
  }
}

// 页面加载时获取数据
onMounted(() => {
  fetchData()
})
</script>

<style scoped>
#blacklistManagePage {
  padding: 24px;
}

:deep(.ant-card-head-title) {
  font-weight: 600;
}

:deep(.ant-card-meta-title) {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 14px;
}
</style>
