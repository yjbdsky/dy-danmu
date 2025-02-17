<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import { useAuthStore } from '../stores/auth'
import { getLiveConfList, deleteLiveConf } from '../api/live-conf'
import type { LiveConf } from '../types/models/live_conf'
import AddLiveConf from '../components/AddLiveConf.vue'
import EditLiveConf from '../components/EditLiveConf.vue'

const router = useRouter()
const authStore = useAuthStore()
const isAdmin = computed(() => authStore.auth?.role === 'admin')
const liveConfList = ref<LiveConf[]>([])
const showAddDialog = ref(false)
const showEditDialog = ref(false)
const currentConf = ref<LiveConf | null>(null)

// 获取直播配置列表
async function fetchList() {
  try {
    const res = await getLiveConfList()
    if (res.data.data?.list && Array.isArray(res.data.data.list)) {
      liveConfList.value = res.data.data.list
    } else {
      console.error('[LiveConf] Invalid data format:', res.data)
    }
  } catch (error) {
    console.error('[LiveConf] Fetch error:', error)
  }
}

// 删除确认
async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm('确定删除该配置?', '提示', {
      type: 'warning'
    })
    await deleteLiveConf(id)
    ElMessage.success('删除成功')
    await fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('[LiveConf] Delete error:', error)
    }
  }
}

// 跳转到详情页
function goToDetail(conf: LiveConf) {
  router.push({
    path: `/live-conf/${conf.id}`
  })
}

function handleAddSuccess() {
  fetchList()
}

function handleEdit(conf: LiveConf) {
  currentConf.value = conf
  showEditDialog.value = true
}

function handleEditSuccess() {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<template>
  <div class="p-3 sm:p-4">
    <!-- 卡片列表 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3 max-w-7xl mx-auto">
      <el-card
        v-for="item in liveConfList"
        :key="item?.id"
        shadow="hover"
        :body-style="{ padding: '12px' }"
        class="cursor-pointer transition-all duration-200 hover:-translate-y-1 hover:shadow-md min-w-0"
        @click="goToDetail(item)"
      >
        <!-- 标题和状态 -->
        <div class="flex justify-between items-start mb-2 min-w-0">
          <h3 class="text-base font-medium truncate mr-2 text-gray-800 flex-1 min-w-0">{{ item.name }}</h3>
          <el-tag 
            :type="item.enable ? 'success' : 'info'" 
            size="small"
            class="flex-shrink-0"
          >
            {{ item.enable ? '已启用' : '未启用' }}
          </el-tag>
        </div>

        <!-- 房间ID -->
        <div class="text-gray-600 text-xs mb-1.5 truncate">
          房间ID: {{ item.room_display_id }}
        </div>

        <!-- URL -->
        <div class="text-gray-500 text-xs truncate mb-3" :title="item.url">
          URL: {{ item.url }}
        </div>

        <!-- 管理按钮 -->
        <div 
          v-if="isAdmin" 
          class="flex justify-end space-x-1.5 pt-2 border-t border-gray-100" 
          @click.stop
        >
          <el-button 
            type="primary" 
            :icon="Edit"
            size="small"
            class="!px-2.5 !h-7"
            @click="handleEdit(item)"
          >
            编辑
          </el-button>
          <el-button
            type="danger"
            :icon="Delete" 
            size="small"
            class="!px-2.5 !h-7"
            @click="handleDelete(item.id)"
          >
            删除
          </el-button>
        </div>
      </el-card>
    </div>

    <AddLiveConf
      v-model="showAddDialog"
      @success="handleAddSuccess"
    />

    <EditLiveConf
      v-model="showEditDialog"
      :live-conf="currentConf"
      @success="handleEditSuccess"
    />

    <el-button
      v-if="isAdmin"
      type="primary"
      :icon="Plus"
      circle
      class="fixed right-4 bottom-4 shadow-lg !w-12 !h-12"
      @click="showAddDialog = true"
    />
  </div>
</template> 