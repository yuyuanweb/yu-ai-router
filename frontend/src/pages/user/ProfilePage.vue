<template>
  <div id="profilePage" class="morandi-dashboard">
    <a-row :gutter="24">
      <!-- 配额卡片 -->
      <a-col :xs="24" :md="12">
        <div class="stat-card">
          <div class="card-header">
            <h3>配额信息</h3>
            <div class="header-icon quota-icon">✦</div>
          </div>
          <a-spin :spinning="loading">
            <div class="stat-content">
              <a-statistic
                title="Token 配额"
                :value="quotaInfo.tokenQuota === -1 ? '无限制' : quotaInfo.tokenQuota"
                class="morandi-stat primary"
              />
              <a-progress
                v-if="quotaInfo.tokenQuota !== -1"
                :percent="usagePercent"
                :status="usagePercent > 90 ? 'exception' : 'normal'"
                stroke-color="#B8A9B8"
                trail-color="#EBE5DF"
                class="morandi-progress"
              />
              <div class="quota-details">
                <div class="detail-item">
                  <span class="label">已使用</span>
                  <span class="value">{{ (quotaInfo.usedTokens || 0).toLocaleString() }} Tokens</span>
                </div>
                <div class="detail-item">
                  <span class="label">剩余配额</span>
                  <span class="value">{{ quotaInfo.remainingQuota === -1 ? '无限制' : (quotaInfo.remainingQuota || 0).toLocaleString() }}</span>
                </div>
              </div>
            </div>
          </a-spin>
        </div>
      </a-col>

      <!-- Token统计卡片 -->
      <a-col :xs="24" :md="12">
        <div class="stat-card">
          <div class="card-header">
            <h3>Token 消耗</h3>
            <div class="header-icon stats-icon">❍</div>
          </div>
          <a-spin :spinning="loading">
            <div class="stat-content">
              <a-statistic
                title="累计消耗"
                :value="summaryStats.totalTokens || 0"
                suffix="Tokens"
                class="morandi-stat accent"
              />
              <div class="mini-stats">
                <div class="mini-item">
                  <div class="mini-label">总请求数</div>
                  <div class="mini-value">{{ summaryStats.totalRequests || 0 }}</div>
                </div>
                <div class="mini-item success">
                  <div class="mini-label">成功请求数</div>
                  <div class="mini-value">{{ summaryStats.successRequests || 0 }}</div>
                </div>
              </div>
            </div>
          </a-spin>
        </div>
      </a-col>
    </a-row>

    <!-- 每日消耗趋势图 -->
    <div class="chart-card">
      <div class="card-header">
        <h3>近 7 日消耗趋势</h3>
        <span class="subtitle">Token 与请求次数</span>
      </div>
      <a-spin :spinning="chartLoading">
        <div ref="chartRef" class="chart-container"></div>
      </a-spin>
    </div>
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

const usagePercent = computed(() => {
  if (!quotaInfo.value.tokenQuota || quotaInfo.value.tokenQuota === -1 || !quotaInfo.value.usedTokens) {
    return 0
  }
  return Math.round((quotaInfo.value.usedTokens / quotaInfo.value.tokenQuota) * 100)
})

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
      chartLoading.value = false
      await nextTick()
      renderChart()
    } else {
      chartLoading.value = false
    }
  } catch (error) {
    chartLoading.value = false
  }
}

// 极致炫酷的莫兰迪图表配置
const renderChart = async () => {
  if (!chartRef.value || dailyStats.value.length === 0) return
  await nextTick()
  if (chartInstance) chartInstance.dispose()
  chartInstance = echarts.init(chartRef.value)

  const dates = dailyStats.value.map((item) => item.date)
  const tokens = dailyStats.value.map((item) => Number(item.totalTokens) || 0)
  const requests = dailyStats.value.map((item) => Number(item.requestCount) || 0)

  const option = {
    color: ['#B8A9B8', '#A9B4C2'], // 雾霾紫、雾蓝
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(250, 247, 244, 0.9)',
      borderColor: '#EBE5DF',
      textStyle: { color: '#4A4744' },
      extraCssText: 'backdrop-filter: blur(10px); border-radius: 12px; box-shadow: 0 4px 20px rgba(184, 169, 184, 0.1);',
      axisPointer: {
        type: 'line',
        lineStyle: {
          color: '#B8A9B8',
          type: 'dashed'
        }
      }
    },
    legend: {
      data: ['Token消耗', '请求次数'],
      top: 0,
      right: 20,
      icon: 'circle',
      textStyle: { color: '#868280' }
    },
    grid: { left: '1%', right: '2%', bottom: '3%', top: '15%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: { show: false }, // 去掉轴线
      axisTick: { show: false }, // 去掉刻度
      axisLabel: { color: '#868280', margin: 16 }
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          type: 'dashed',
          color: '#EBE5DF' // 极浅的分割线
        }
      },
      axisLabel: { color: '#868280' }
    },
    series: [
      {
        name: 'Token消耗',
        type: 'line',
        smooth: 0.4, // 平滑曲线
        symbol: 'circle',
        symbolSize: 8,
        showSymbol: false, // 默认不显示点，hover才显示
        lineStyle: { width: 3 },
        areaStyle: { // 炫酷的渐变填充
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(184, 169, 184, 0.4)' },
            { offset: 1, color: 'rgba(184, 169, 184, 0.0)' }
          ])
        },
        data: tokens
      },
      {
        name: '请求次数',
        type: 'line',
        smooth: 0.4,
        symbol: 'circle',
        symbolSize: 8,
        showSymbol: false,
        lineStyle: { width: 2, type: 'dashed' }, // 虚线区分
        data: requests
      }
    ]
  }
  chartInstance.setOption(option)
}

const handleResize = () => { if (chartInstance) chartInstance.resize() }

onMounted(() => {
  loadSummaryStats()
  loadDailyStats()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (chartInstance) { chartInstance.dispose(); chartInstance = null }
})
</script>

<style scoped>
.morandi-dashboard {
  padding: 24px;
  background: var(--bg-page);
  min-height: calc(100vh - 72px);
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 统一卡片外观 */
.stat-card, .chart-card {
  background: var(--bg-card);
  border-radius: var(--radius-lg);
  padding: 28px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-soft);
  height: 100%;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-main);
  letter-spacing: 0.5px;
}

.subtitle {
  font-size: 12px;
  color: var(--text-secondary);
}

.header-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
}
.quota-icon { background: rgba(184, 169, 184, 0.15); color: var(--color-primary); }
.stats-icon { background: rgba(169, 180, 194, 0.15); color: var(--color-secondary); }

/* 数字样式覆盖 */
:deep(.morandi-stat .ant-statistic-title) {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 8px;
}
:deep(.morandi-stat .ant-statistic-content) {
  font-family: 'DIN Alternate', 'Helvetica Neue', sans-serif; /* 如果有数字字体更好看 */
}
:deep(.morandi-stat.primary .ant-statistic-content) {
  color: var(--color-primary);
}
:deep(.morandi-stat.accent .ant-statistic-content) {
  color: var(--color-secondary);
}

.morandi-progress {
  margin: 20px 0;
}

.quota-details {
  display: flex;
  justify-content: space-between;
  padding-top: 20px;
  border-top: 1px dashed var(--border-color);
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.detail-item .label { font-size: 12px; color: var(--text-secondary); }
.detail-item .value { font-size: 16px; font-weight: 600; color: var(--text-main); }

/* 迷你统计框 */
.mini-stats {
  display: flex;
  gap: 16px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px dashed var(--border-color);
}

.mini-item {
  flex: 1;
  background: rgba(245, 242, 239, 0.5);
  padding: 16px;
  border-radius: 12px;
}
.mini-item.success { background: rgba(168, 191, 168, 0.1); }
.mini-label { font-size: 12px; color: var(--text-secondary); margin-bottom: 6px; }
.mini-value { font-size: 20px; font-weight: 600; color: var(--text-main); }
.mini-item.success .mini-value { color: #6C9A7C; }

/* 图表容器 */
.chart-container {
  width: 100%;
  height: 380px;
}
</style>
