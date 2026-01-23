<template>
  <div id="profilePage">
    <!-- 余额卡片 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="24">
        <a-card title="账户余额" :loading="balanceLoading">
          <a-row :gutter="16">
            <a-col :span="8">
              <a-statistic
                title="当前余额"
                :value="balanceInfo.balance || 0"
                prefix="¥"
                :precision="2"
                :value-style="{
                  color: (balanceInfo.balance || 0) > 10 ? '#3f8600' : '#cf1322',
                  fontSize: '32px',
                  fontWeight: 'bold'
                }"
              />
              <a-button
                type="primary"
                size="large"
                @click="showRechargeModal"
                style="margin-top: 16px; width: 100%"
              >
                立即充值
              </a-button>
            </a-col>
            <a-col :span="8">
              <a-statistic
                title="累计充值"
                :value="balanceInfo.totalRecharge || 0"
                prefix="¥"
                :precision="2"
                :value-style="{ color: '#1890ff' }"
              />
            </a-col>
            <a-col :span="8">
              <a-statistic
                title="累计消费"
                :value="balanceInfo.totalSpending || 0"
                prefix="¥"
                :precision="2"
                :value-style="{ color: '#faad14' }"
              />
            </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>

    <a-row :gutter="16">
      <!-- 配额卡片 -->
      <a-col :span="8">
        <a-card title="配额信息" :loading="loading">
          <a-statistic
            title="Token配额"
            :value="quotaInfo.tokenQuota === -1 ? '无限制' : quotaInfo.tokenQuota"
            :value-style="{ color: '#3f8600' }"
          />
          <a-divider />
          <a-statistic
            title="已使用"
            :value="quotaInfo.usedTokens || 0"
            suffix="Tokens"
          />
          <a-divider />
          <a-statistic
            title="剩余配额"
            :value="
              quotaInfo.remainingQuota === -1 ? '无限制' : quotaInfo.remainingQuota
            "
            :value-style="{ color: '#cf1322' }"
          />
          <a-progress
            v-if="quotaInfo.tokenQuota !== -1"
            :percent="usagePercent"
            :status="usagePercent > 90 ? 'exception' : 'normal'"
            style="margin-top: 16px"
          />
        </a-card>
      </a-col>

      <!-- Token统计卡片 -->
      <a-col :span="8">
        <a-card title="Token消耗" :loading="loading">
          <a-statistic
            title="累计消耗"
            :value="summaryStats.totalTokens || 0"
            suffix="Tokens"
            :value-style="{ color: '#1890ff' }"
          />
          <a-divider />
          <a-statistic
            title="总请求数"
            :value="summaryStats.totalRequests || 0"
          />
          <a-divider />
          <a-statistic
            title="成功请求数"
            :value="summaryStats.successRequests || 0"
            :value-style="{ color: '#52c41a' }"
          />
        </a-card>
      </a-col>

      <!-- 费用卡片 -->
      <a-col :span="8">
        <a-card title="费用统计" :loading="loading">
          <a-statistic
            title="累计消费"
            :value="summaryStats.totalCost || 0"
            prefix="¥"
            :precision="2"
            :value-style="{ color: '#faad14' }"
          />
          <a-divider />
          <a-statistic
            title="今日消费"
            :value="summaryStats.todayCost || 0"
            prefix="¥"
            :precision="2"
          />
        </a-card>
      </a-col>
    </a-row>

    <!-- 每日消耗趋势图 -->
    <a-card title="每日消耗趋势（最近7天）" style="margin-top: 16px" :loading="chartLoading">
      <div ref="chartRef" style="width: 100%; height: 400px"></div>
    </a-card>

    <!-- 充值记录和消费账单 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <!-- 充值记录 -->
      <a-col :span="12">
        <a-card title="充值记录" :loading="rechargeLoading">
          <a-table
            :columns="rechargeColumns"
            :data-source="rechargeRecords"
            :pagination="rechargePagination"
            @change="handleRechargeTableChange"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'amount'">
                <span style="color: #52c41a; font-weight: bold">
                  ¥{{ record.amount?.toFixed(4) }}
                </span>
              </template>
              <template v-else-if="column.dataIndex === 'status'">
                <a-tag v-if="record.status === 'success'" color="green">成功</a-tag>
                <a-tag v-else-if="record.status === 'pending'" color="orange">处理中</a-tag>
                <a-tag v-else-if="record.status === 'failed'" color="red">失败</a-tag>
                <a-tag v-else color="gray">{{ record.status }}</a-tag>
              </template>
              <template v-else-if="column.dataIndex === 'paymentMethod'">
                <a-tag v-if="record.paymentMethod === 'stripe'" color="blue">Stripe</a-tag>
                <a-tag v-else-if="record.paymentMethod === 'alipay'" color="cyan">支付宝</a-tag>
                <a-tag v-else>{{ record.paymentMethod }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>

      <!-- 消费账单 -->
      <a-col :span="12">
        <a-card title="消费账单" :loading="billingLoading">
          <a-table
            :columns="billingColumns"
            :data-source="billingRecords"
            :pagination="billingPagination"
            @change="handleBillingTableChange"
            size="small"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'amount'">
                <!-- 充值显示正数（绿色），消费显示负数（红色） -->
                <span
                  v-if="record.billingType === 'recharge'"
                  style="color: #52c41a; font-weight: bold"
                >
                  +¥{{ record.amount?.toFixed(4) }}
                </span>
                <span
                  v-else
                  style="color: #cf1322; font-weight: bold"
                >
                  -¥{{ record.amount?.toFixed(4) }}
                </span>
              </template>
              <template v-else-if="column.dataIndex === 'billingType'">
                <a-tag v-if="record.billingType === 'api_call'" color="orange">
                  API调用
                </a-tag>
                <a-tag v-else-if="record.billingType === 'recharge'" color="green">
                  充值
                </a-tag>
                <a-tag v-else>{{ record.billingType }}</a-tag>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 充值弹窗 -->
    <a-modal
      v-model:open="rechargeModalVisible"
      title="账户充值"
      :footer="null"
      width="500px"
    >
      <a-form :model="rechargeForm" layout="vertical" @finish="handleRecharge">
        <a-form-item label="充值金额（元）" name="amount" :rules="[{ required: true, message: '请输入充值金额' }]">
          <a-input-number
            v-model:value="rechargeForm.amount"
            :min="1"
            :max="10000"
            :step="10"
            :precision="2"
            style="width: 100%"
            placeholder="输入充值金额（1-10000元）"
          >
            <template #addonBefore>¥</template>
          </a-input-number>
        </a-form-item>

        <a-form-item label="快捷金额">
          <a-space>
            <a-button @click="rechargeForm.amount = 10">10元</a-button>
            <a-button @click="rechargeForm.amount = 50">50元</a-button>
            <a-button @click="rechargeForm.amount = 100">100元</a-button>
            <a-button @click="rechargeForm.amount = 500">500元</a-button>
          </a-space>
        </a-form-item>

        <a-form-item label="支付方式">
          <a-radio-group v-model:value="rechargeForm.paymentMethod">
            <a-radio value="stripe">
              <span style="font-weight: bold">Stripe</span>
              <span style="color: #999; margin-left: 8px">（支持国际信用卡）</span>
            </a-radio>
          </a-radio-group>
        </a-form-item>

        <a-alert
          message="测试环境提示"
          description="当前为 Stripe 沙箱环境，请使用测试卡号：4242 4242 4242 4242"
          type="info"
          show-icon
          style="margin-bottom: 16px"
        />

        <a-form-item>
          <a-space>
            <a-button type="primary" html-type="submit" :loading="recharging">
              立即支付
            </a-button>
            <a-button @click="rechargeModalVisible = false">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref, computed, nextTick } from 'vue'
import { getMySummaryStats, getMyDailyStats } from '@/api/statsController'
import { getMyBalance, getMyBillingRecords } from '@/api/balanceController'
import { createStripeRecharge, getMyRechargeRecords } from '@/api/rechargeController'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'

const loading = ref(false)
const chartLoading = ref(false)
const balanceLoading = ref(false)
const rechargeLoading = ref(false)
const billingLoading = ref(false)
const recharging = ref(false)

const quotaInfo = ref<any>({})
const summaryStats = ref<API.UserSummaryStatsVO>({})
const dailyStats = ref<any[]>([])
const balanceInfo = ref<API.BalanceVO>({})
const rechargeRecords = ref<API.RechargeRecord[]>([])
const billingRecords = ref<API.BillingRecord[]>([])
const chartRef = ref()
let chartInstance: echarts.ECharts | null = null

// 充值弹窗
const rechargeModalVisible = ref(false)
const rechargeForm = ref({
  amount: 10,
  paymentMethod: 'stripe',
})

// 分页配置
const rechargePagination = ref({
  current: 1,
  pageSize: 5,
  total: 0,
  showSizeChanger: false,
})

const billingPagination = ref({
  current: 1,
  pageSize: 5,
  total: 0,
  showSizeChanger: false,
})

// 计算配额使用百分比
const usagePercent = computed(() => {
  if (
    !quotaInfo.value.tokenQuota ||
    quotaInfo.value.tokenQuota === -1 ||
    !quotaInfo.value.usedTokens
  ) {
    return 0
  }
  return Math.round((quotaInfo.value.usedTokens / quotaInfo.value.tokenQuota) * 100)
})

// 充值记录表格列
const rechargeColumns = [
  {
    title: '金额',
    dataIndex: 'amount',
    width: 120,
  },
  {
    title: '支付方式',
    dataIndex: 'paymentMethod',
    width: 110,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 90,
  },
  {
    title: '时间',
    dataIndex: 'createTime',
    width: 180,
    customRender: ({ text }: any) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
  },
]

// 消费账单表格列
const billingColumns = [
  {
    title: '金额',
    dataIndex: 'amount',
    width: 130,
  },
  {
    title: '类型',
    dataIndex: 'billingType',
    width: 100,
  },
  {
    title: '说明',
    dataIndex: 'description',
    ellipsis: true,
  },
  {
    title: '时间',
    dataIndex: 'createTime',
    width: 180,
    customRender: ({ text }: any) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
  },
]

// 加载综合统计数据
const loadSummaryStats = async () => {
  loading.value = true
  try {
    const res = await getMySummaryStats()
    if (res.data.code === 0 && res.data.data) {
      summaryStats.value = res.data.data
      quotaInfo.value = {
        tokenQuota: res.data.data.tokenQuota,
        usedTokens: res.data.data.usedTokens,
        remainingQuota: res.data.data.remainingQuota,
      }
    } else {
      message.error('获取统计数据失败：' + res.data.message)
    }
  } catch (error) {
    message.error('获取统计数据失败')
  } finally {
    loading.value = false
  }
}

// 加载每日统计数据
const loadDailyStats = async () => {
  chartLoading.value = true
  try {
    const endDate = new Date()
    const startDate = new Date()
    startDate.setDate(startDate.getDate() - 6)

    const res = await getMyDailyStats({
      startDate: startDate.toISOString().split('T')[0],
      endDate: endDate.toISOString().split('T')[0],
    })

    if (res.data.code === 0 && res.data.data) {
      dailyStats.value = res.data.data
      // 等待 loading 状态更新后再渲染图表
      chartLoading.value = false
      await nextTick()
      renderChart()
    } else {
      message.error('获取每日统计失败：' + res.data.message)
      chartLoading.value = false
    }
  } catch (error) {
    message.error('获取每日统计失败')
    chartLoading.value = false
  }
}

// 渲染图表
const renderChart = async () => {
  if (!chartRef.value || dailyStats.value.length === 0) {
    return
  }

  await nextTick()

  // 如果已存在图表实例，先销毁
  if (chartInstance) {
    chartInstance.dispose()
  }

  // 初始化图表
  chartInstance = echarts.init(chartRef.value)

  // 提取并转换数据为数字类型
  const dates = dailyStats.value.map((item) => item.date)
  const tokens = dailyStats.value.map((item) => Number(item.totalTokens) || 0)
  const costs = dailyStats.value.map((item) => Number(item.totalCost) || 0)
  const requests = dailyStats.value.map((item) => Number(item.requestCount) || 0)

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
      formatter: (params: any) => {
        let result = `${params[0].axisValue}<br/>`
        params.forEach((param: any) => {
          const value = param.seriesName === '费用（元）'
            ? `¥${param.value.toFixed(2)}`
            : param.value
          result += `${param.marker}${param.seriesName}: ${value}<br/>`
        })
        return result
      },
    },
    legend: {
      data: ['Token消耗', '费用（元）', '请求次数'],
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
    },
    yAxis: [
      {
        type: 'value',
        name: 'Tokens / 请求次数',
        position: 'left',
      },
      {
        type: 'value',
        name: '费用（元）',
        position: 'right',
        axisLabel: {
          formatter: '¥{value}',
        },
      },
    ],
    series: [
      {
        name: 'Token消耗',
        type: 'line',
        data: tokens,
        smooth: true,
        itemStyle: {
          color: '#1890ff',
        },
        lineStyle: {
          color: '#1890ff',
        },
      },
      {
        name: '费用（元）',
        type: 'line',
        yAxisIndex: 1,
        data: costs,
        smooth: true,
        itemStyle: {
          color: '#faad14',
        },
        lineStyle: {
          color: '#faad14',
        },
      },
      {
        name: '请求次数',
        type: 'line',
        data: requests,
        smooth: true,
        itemStyle: {
          color: '#52c41a',
        },
        lineStyle: {
          color: '#52c41a',
        },
      },
    ],
  }

  chartInstance.setOption(option)
}

// 响应式调整
const handleResize = () => {
  if (chartInstance) {
    chartInstance.resize()
  }
}

// 加载余额信息
const loadBalance = async () => {
  balanceLoading.value = true
  try {
    const res = await getMyBalance()
    if (res.data.code === 0 && res.data.data) {
      balanceInfo.value = res.data.data
    } else {
      message.error('获取余额失败：' + res.data.message)
    }
  } catch (error) {
    message.error('获取余额失败')
  } finally {
    balanceLoading.value = false
  }
}

// 加载充值记录
const loadRechargeRecords = async () => {
  rechargeLoading.value = true
  try {
    const res = await getMyRechargeRecords({
      pageNum: rechargePagination.value.current,
      pageSize: rechargePagination.value.pageSize,
    })
    if (res.data.code === 0 && res.data.data) {
      rechargeRecords.value = res.data.data.records || []
      rechargePagination.value.total = Number(res.data.data.totalRow) || 0
    } else {
      message.error('获取充值记录失败')
    }
  } catch (error) {
    message.error('获取充值记录失败')
  } finally {
    rechargeLoading.value = false
  }
}

// 加载消费账单
const loadBillingRecords = async () => {
  billingLoading.value = true
  try {
    const res = await getMyBillingRecords({
      pageNum: billingPagination.value.current,
      pageSize: billingPagination.value.pageSize,
    })
    if (res.data.code === 0 && res.data.data) {
      billingRecords.value = res.data.data.records || []
      billingPagination.value.total = Number(res.data.data.totalRow) || 0
    } else {
      message.error('获取消费账单失败')
    }
  } catch (error) {
    message.error('获取消费账单失败')
  } finally {
    billingLoading.value = false
  }
}

// 显示充值弹窗
const showRechargeModal = () => {
  rechargeModalVisible.value = true
}

// 处理充值
const handleRecharge = async () => {
  if (!rechargeForm.value.amount || rechargeForm.value.amount <= 0) {
    message.error('请输入有效的充值金额')
    return
  }

  recharging.value = true
  try {
    const res = await createStripeRecharge({
      amount: rechargeForm.value.amount,
    })

    if (res.data.code === 0 && res.data.data?.checkoutUrl) {
      // 跳转到 Stripe 支付页面
      message.success('正在跳转到支付页面...')
      window.location.href = res.data.data.checkoutUrl
    } else {
      message.error('创建充值订单失败：' + res.data.message)
    }
  } catch (error: any) {
    message.error('创建充值订单失败：' + (error.message || '未知错误'))
  } finally {
    recharging.value = false
  }
}

// 充值记录表格分页变化
const handleRechargeTableChange = (pagination: any) => {
  rechargePagination.value.current = pagination.current
  loadRechargeRecords()
}

// 消费账单表格分页变化
const handleBillingTableChange = (pagination: any) => {
  billingPagination.value.current = pagination.current
  loadBillingRecords()
}

// 组件挂载时添加监听
onMounted(() => {
  loadBalance()
  loadSummaryStats()
  loadDailyStats()
  loadRechargeRecords()
  loadBillingRecords()
  window.addEventListener('resize', handleResize)
})

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
})
</script>

<style scoped>
#profilePage {
  padding: 24px;
  background: #f0f2f5;
  min-height: calc(100vh - 64px);
}
</style>
