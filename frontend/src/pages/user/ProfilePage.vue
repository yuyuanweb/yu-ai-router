<template>
  <div id="profilePage">
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
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref, computed, nextTick } from 'vue'
import { getMySummaryStats, getMyDailyStats } from '@/api/statsController'
import { message } from 'ant-design-vue'
import * as echarts from 'echarts'

const loading = ref(false)
const chartLoading = ref(false)
const quotaInfo = ref<any>({})
const summaryStats = ref<API.UserSummaryStatsVO>({})
const dailyStats = ref<any[]>([])
const chartRef = ref()
let chartInstance: echarts.ECharts | null = null

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

// 组件挂载时添加监听
onMounted(() => {
  loadSummaryStats()
  loadDailyStats()
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
