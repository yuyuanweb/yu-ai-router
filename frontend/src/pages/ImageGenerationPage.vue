<template>
  <div id="imageGenerationPage">
    <a-row :gutter="16">
      <!-- 生成区域 -->
      <a-col :span="12">
        <a-card title="AI 文生图" :loading="generating">
          <a-form layout="vertical">
            <a-form-item label="描述你想要的图片">
              <a-textarea
                v-model:value="formData.prompt"
                placeholder="例如：一只可爱的猫咪在花园里玩耍，阳光明媚，卡通风格"
                :rows="4"
                :maxlength="500"
                show-count
              />
            </a-form-item>

            <a-form-item label="选择模型">
              <a-radio-group v-model:value="formData.model">
                <a-radio value="qwen-image-plus">
                  <span style="font-weight: bold">通义万相</span>
                  <span style="color: #999; margin-left: 8px">(0.08元/张)</span>
                </a-radio>
                <a-radio value="cogview-3-plus" disabled>
                  <span>智谱CogView</span>
                  <span style="color: #999; margin-left: 8px">(0.10元/张)</span>
                </a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item label="图片尺寸">
              <a-select v-model:value="formData.size">
                <a-select-option value="720*720">720x720</a-select-option>
                <a-select-option value="1024*1024">1024x1024（推荐）</a-select-option>
                <a-select-option value="1280*720">1280x720</a-select-option>
              </a-select>
            </a-form-item>

            <a-alert
              v-if="estimatedCost > 0"
              :message="`预估费用：¥${estimatedCost.toFixed(4)}`"
              type="info"
              show-icon
              style="margin-bottom: 16px"
            />

            <a-form-item>
              <a-button
                type="primary"
                size="large"
                :loading="generating"
                @click="handleGenerate"
                block
              >
                {{ generating ? '生成中...' : '开始生成' }}
              </a-button>
            </a-form-item>
          </a-form>

          <!-- 生成结果 -->
          <div v-if="generatedImages.length > 0" style="margin-top: 24px">
            <a-divider>生成结果</a-divider>
            <a-row :gutter="16">
              <a-col
                v-for="(image, index) in generatedImages"
                :key="index"
                :span="24"
              >
                <a-card hoverable>
                  <a-image :src="image.url" :alt="`生成图片 ${index + 1}`" />
                  <template #actions>
                    <a :href="image.url" target="_blank" title="在新标签页打开">
                      <EyeOutlined />
                    </a>
                    <a :href="image.url" download title="下载图片">
                      <DownloadOutlined />
                    </a>
                  </template>
                </a-card>
              </a-col>
            </a-row>
          </div>
        </a-card>
      </a-col>

      <!-- 历史记录 -->
      <a-col :span="12">
        <a-card title="生成历史" :loading="loading">
          <a-list
            :data-source="historyRecords"
            :pagination="pagination"
            @change="handlePageChange"
          >
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space>
                      <span>{{ item.prompt }}</span>
                      <a-tag v-if="item.status === 'success'" color="green">成功</a-tag>
                      <a-tag v-else color="red">失败</a-tag>
                    </a-space>
                  </template>
                  <template #description>
                    <div>
                      <div>模型：{{ item.modelKey }}</div>
                      <div>尺寸：{{ item.size }}</div>
                      <div>费用：¥{{ item.cost?.toFixed(4) }}</div>
                      <div>时间：{{ formatTime(item.createTime) }}</div>
                    </div>
                  </template>
                  <template #avatar v-if="item.imageUrl && item.status === 'success'">
                    <a-image
                      :src="item.imageUrl"
                      :width="80"
                      :height="80"
                      :preview="true"
                      style="border-radius: 4px; object-fit: cover"
                    />
                  </template>
                </a-list-item-meta>
                <template #actions v-if="item.imageUrl && item.status === 'success'">
                  <a @click="viewImage(item.imageUrl)" title="查看">
                    <EyeOutlined />
                  </a>
                  <a :href="item.imageUrl" download title="下载">
                    <DownloadOutlined />
                  </a>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { generateImage, getMyRecords } from '@/api/imageController'
import { message } from 'ant-design-vue'
import { EyeOutlined, DownloadOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'

const generating = ref(false)
const loading = ref(false)

const formData = ref({
  prompt: '',
  model: 'qwen-image-plus',
  size: '1024*1024',
})

const generatedImages = ref<API.ImageData[]>([])
const historyRecords = ref<API.ImageGenerationRecord[]>([])

const pagination = ref({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: false,
})

// 预估费用（固定生成1张）
const estimatedCost = computed(() => {
  const prices: Record<string, number> = {
    'qwen-image-plus': 0.08,
    'cogview-3-plus': 0.1,
  }
  return prices[formData.value.model] || 0
})

// 格式化时间
const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

// 生成图片
const handleGenerate = async () => {
  if (!formData.value.prompt || formData.value.prompt.trim() === '') {
    message.error('请输入图片描述')
    return
  }

  generating.value = true
  try {
    const res = await generateImage({
      prompt: formData.value.prompt,
      model: formData.value.model,
      n: 1,  // 固定生成1张
      size: formData.value.size,
    })

    if (res.data && res.data.data) {
      generatedImages.value = res.data.data
      message.success('图片生成成功！')
      
      // 刷新历史记录
      loadHistory()
    } else {
      message.error('图片生成失败')
    }
  } catch (error: any) {
    message.error('图片生成失败：' + (error.message || '未知错误'))
  } finally {
    generating.value = false
  }
}

// 加载历史记录
const loadHistory = async () => {
  loading.value = true
  try {
    const res = await getMyRecords({
      pageNum: pagination.value.current,
      pageSize: pagination.value.pageSize,
    })

    if (res.data.code === 0 && res.data.data) {
      historyRecords.value = res.data.data.records || []
      pagination.value.total = Number(res.data.data.totalRow) || 0
    } else {
      message.error('获取历史记录失败')
    }
  } catch (error) {
    message.error('获取历史记录失败')
  } finally {
    loading.value = false
  }
}

// 查看图片
const viewImage = (url: string) => {
  window.open(url, '_blank')
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.value.current = page
  loadHistory()
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
#imageGenerationPage {
  padding: 24px;
  background: #f0f2f5;
  min-height: calc(100vh - 64px);
}
</style>
