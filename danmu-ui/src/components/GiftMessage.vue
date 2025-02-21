<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { getGiftMessage } from '../api/gift_message'
import type { GiftMessage } from '../types/models/gift_message'
import TimeRangePicker from './TimeRangePicker.vue'
import ToUserSelect from './ToUserSelect.vue'
import UserMultiSelect from './UserMultiSelect.vue'
import dayjs from 'dayjs'

const props = defineProps<{
  roomDisplayId: string
}>()

// 查询参数
const search = ref('')
const userIds = ref<number[]>([])
const toUserIds = ref<number[]>([])
const timeRange = ref<[number, number]>([
  dayjs().startOf('day').valueOf(),
  dayjs().endOf('day').valueOf()
])
const diamondCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(100)
const orderBy = ref('')
const orderDirection = ref('')

const loading = ref(false)
const giftList = ref<GiftMessage[]>([])
const total = ref(0)

const tableContainer = ref<HTMLElement | null>(null)

// 初始化排序为时间倒序
orderBy.value = 'timestamp'
orderDirection.value = 'desc'

async function fetchGiftMessage(resetPage = true) {
  if (resetPage) {
    currentPage.value = 1
  }
  loading.value = true
  try {
    const res = await getGiftMessage({
      search: search.value,
      user_ids: userIds.value,
      to_user_ids: toUserIds.value,
      room_display_id: props.roomDisplayId,
      begin: timeRange.value[0],
      end: timeRange.value[1],
      diamond_count: diamondCount.value,
      page: currentPage.value,
      page_size: pageSize.value,
      order_by: orderBy.value,
      order_direction: orderDirection.value
    })
    giftList.value = res.data.data.list
    total.value = res.data.data.total
  } catch (error) {
    console.error('[GiftMessage] Fetch error:', error)
  } finally {
    loading.value = false
  }
}

function handleSort({ prop, order }: { prop: string; order: string }) {
  orderBy.value = prop
  orderDirection.value = order === 'ascending' ? 'asc' : 'desc'
  fetchGiftMessage()
}

function formatTime(timestamp: number) {
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

// 处理底部分页器页码改变（滚动到顶部）
function handleBottomPageChange(page: number) {
  currentPage.value = page
  fetchGiftMessage(false)
  nextTick(() => {
    tableContainer.value?.scrollIntoView({ behavior: 'smooth' })
  })
}

// 处理顶部分页器页码改变（不滚动）
function handleTopPageChange(page: number) {
  currentPage.value = page
  fetchGiftMessage(false)
}

onMounted(fetchGiftMessage)
</script>

<template>
  <div class="space-y-4">
    <!-- 筛选条件 -->
    <div class="space-y-3">
      <!-- 搜索框独占一行 -->
      <div>
        <el-input
          v-model="search"
          placeholder="搜索"
          class="!w-[300px] sm:!w-[565px]"
          clearable
          @keyup.enter="fetchGiftMessage"
        />
      </div>
      
      <!-- 其他条件分两列 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-2 gap-y-3 max-w-[640px]">
        <!-- 左列 -->
        <div class="space-y-3">
          <TimeRangePicker v-model="timeRange" />
        </div>
        
        <!-- 右列 -->
        <div class="space-y-3">
          <UserMultiSelect v-model="userIds" />
          <ToUserSelect 
            v-model="toUserIds"
            :room-display-id="roomDisplayId"
          />
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-2">
              <span class="text-gray-600 whitespace-nowrap">钻石 ≥</span>
              <el-input-number
                v-model="diamondCount"
                :min="0"
                :default-value="0"
                :controls="true"
                placeholder=""
                class="!w-28"
                :step="1"
              />
            </div>
            <el-button 
              type="primary"
              :loading="loading"
              @click="fetchGiftMessage"
            >
              查询
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- 礼物消息列表 -->
    <div 
      ref="tableContainer" 
      class="overflow-x-auto"
    >
      <!-- 顶部分页器 -->
      <div v-if="total > 0" class="flex justify-end mb-4">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :pager-count="5"
          size="small"
          layout="total, sizes, prev, pager, next"
          @current-change="handleTopPageChange"
        />
      </div>

      <el-table 
        :data="giftList" 
        class="w-full"
        size="small"
        @sort-change="handleSort"
      >
        <el-table-column type="expand" width="20">
          <template #default="scope">
            <div class="p-3 bg-gray-50">
              <div class="grid grid-cols-2 gap-3 max-w-2xl">
                <div class="space-y-2">
                  <el-tag size="small" type="success" class="w-full !justify-center">
                    {{ scope.row.user_name }}
                  </el-tag>
                  <el-tag size="small" type="info" class="w-full !justify-center">
                    id: {{ scope.row.user_display_id }}
                  </el-tag>
                </div>
                <div class="space-y-2">
                  <el-tag size="small" type="warning" class="w-full !justify-center">
                    {{ scope.row.to_user_name || '主播' }}
                  </el-tag>
                  <el-tag size="small" type="info" class="w-full !justify-center">
                    id: {{ scope.row.to_user_display_id || '主播' }}
                  </el-tag>
                </div>
                <div class="space-y-2">
                  <el-tag size="small" type="primary" class="w-full !justify-center">
                    {{ formatTime(scope.row.timestamp) }}
                  </el-tag>
                </div>
                <div class="space-y-2">
                  <el-tag size="small" type="warning" class="w-full !justify-center">
                    {{ scope.row.gift_name }} ：{{ scope.row.diamond_count }} 钻石
                  </el-tag>
                </div>
                <div class="col-span-2">
                  <el-tag size="small" type="success" class="w-full !justify-center">
                    {{ scope.row.message }}
                  </el-tag>
                </div>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column 
          label="送礼用户" 
          prop="user_name"
          sortable="custom"
          min-width="80"
          show-overflow-tooltip
        />

        <el-table-column 
          label="接收用户" 
          prop="to_user_name"
          sortable="custom"
          min-width="80"
          show-overflow-tooltip
        />

        <el-table-column 
          label="礼物" 
          prop="gift_name"
          sortable="custom"
          min-width="60"
          show-overflow-tooltip
        />

        <el-table-column 
          label="连击" 
          width="50"
          align="right"
        >
          <template #default="{ row }">
            <span class="text-orange-500 text-[11px] tabular-nums">
              {{ row.combo_count }}
            </span>
          </template>
        </el-table-column>

        <el-table-column 
          label="时间" 
          prop="timestamp"
          sortable="custom"
          min-width="70"
        >
          <template #default="{ row }">
            <div class="flex flex-col text-[11px] text-blue-500">
              <span>{{ dayjs(row.timestamp).format('YY.MM.DD') }}</span>
              <span>{{ dayjs(row.timestamp).format('HH:mm:ss') }}</span>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 修改底部分页器 -->
      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :pager-count="5"
          size="small"
          layout="total, sizes, prev, pager, next"
          @current-change="handleBottomPageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.el-table :deep(.el-table__expanded-cell) {
  padding: 0;
}

.el-table tr {
  cursor: pointer;
}

.el-table tr:hover > td {
  background-color: var(--el-table-row-hover-bg-color);
}
</style> 