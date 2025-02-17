<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getLiveConf } from '../../api/live-conf'
import type { LiveConf } from '../../types/models/live_conf'
import { CopyDocument } from '@element-plus/icons-vue'
import GiftRanking from '../../components/GiftRanking.vue'
import GiftMessage from '../../components/GiftMessage.vue'
import CommonMessage from '../../components/CommonMessage.vue'

const route = useRoute()
const activeTab = ref('gift-rank')
const liveConf = ref<LiveConf | null>(null)
const loading = ref(true)

function truncateUrl(url: string) {
  return url.length > 30 ? url.substring(0, 27) + '...' : url
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text)
    .then(() => {
      ElMessage.success('已复制到剪贴板')
    })
    .catch(() => {
      ElMessage.error('复制失败')
    })
}

async function fetchLiveConf() {
  loading.value = true
  try {
    const id = route.params.id as string
    const res = await getLiveConf(id)
    liveConf.value = res.data.data
  } catch (error) {
    console.error('[LiveConf] Detail fetch error:', error)
    ElMessage.error('获取直播配置详情失败')
  } finally {
    loading.value = false
  }
}

onMounted(fetchLiveConf)
</script>

<template>
  <div class="p-3 sm:p-4 max-w-7xl mx-auto">
    <!-- LiveConf 信息 -->
    <div v-if="liveConf" class="bg-white rounded-lg shadow-sm p-4 sm:p-6 mb-4 sm:mb-6">
      <h2 class="text-lg sm:text-xl font-medium mb-3 sm:mb-4">{{ liveConf.name }}</h2>
      <div class="space-y-2 text-gray-600 text-sm sm:text-base">
        <div class="flex items-center">
          <span class="w-20 sm:w-24">房间ID:</span>
          <span>{{ liveConf.room_display_id }}</span>
        </div>
        <div class="flex items-center">
          <span class="w-20 sm:w-24">URL:</span>
          <span 
            class="cursor-pointer" 
            :title="liveConf?.url"
            @click="() => liveConf && copyToClipboard(liveConf.url)"
          >
            {{ liveConf ? truncateUrl(liveConf.url) : '' }}
            <el-icon class="ml-1 text-gray-400"><CopyDocument /></el-icon>
          </span>
        </div>
        <div class="flex items-center">
          <span class="w-20 sm:w-24">状态:</span>
          <el-tag :type="liveConf.enable ? 'success' : 'info'" size="small">
            {{ liveConf.enable ? '已启用' : '未启用' }}
          </el-tag>
        </div>
      </div>
    </div>

    <!-- 标签页 -->
    <el-tabs v-model="activeTab" class="bg-white rounded-lg shadow-sm p-3 sm:p-4">
      <el-tab-pane label="礼物排行" name="gift-rank">
        <div v-if="loading" class="h-96 flex items-center justify-center">
          <el-skeleton :rows="5" animated />
        </div>
        <GiftRanking 
          v-else-if="liveConf?.room_display_id"
          :room-display-id="liveConf.room_display_id"
        />
        <div v-else class="h-96 flex items-center justify-center text-gray-400">
          无法加载房间信息
        </div>
      </el-tab-pane>
      
      <el-tab-pane label="礼物消息" name="gift-message">
        <div v-if="loading" class="h-96 flex items-center justify-center">
          <el-skeleton :rows="5" animated />
        </div>
        <GiftMessage
          v-else-if="liveConf?.room_display_id"
          :room-display-id="liveConf.room_display_id"
        />
        <div v-else class="h-96 flex items-center justify-center text-gray-400">
          无法加载房间信息
        </div>
      </el-tab-pane>
      
      <el-tab-pane label="普通消息" name="common-message">
        <div v-if="loading" class="h-96 flex items-center justify-center">
          <el-skeleton :rows="5" animated />
        </div>
        <CommonMessage
          v-else-if="liveConf?.room_display_id"
          :room-display-id="liveConf.room_display_id"
        />
        <div v-else class="h-96 flex items-center justify-center text-gray-400">
          无法加载房间信息
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>