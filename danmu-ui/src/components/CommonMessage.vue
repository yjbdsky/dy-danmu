<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getCommonMessage } from '../api/common_message'
import type { CommonMessage } from '../types/models/common_message'
import TimeRangePicker from './TimeRangePicker.vue'
import UserMultiSelect from './UserMultiSelect.vue'
import dayjs from 'dayjs'

const props = defineProps<{
  roomDisplayId: string
}>()

// 查询参数
const search = ref('')
const userIds = ref<number[]>([])
const timeRange = ref<[number, number]>([
  dayjs().startOf('day').valueOf(),
  dayjs().endOf('day').valueOf()
])
const currentPage = ref(1)
const pageSize = ref(100)
const orderBy = ref('')
const orderDirection = ref('')

const loading = ref(false)
const messageList = ref<CommonMessage[]>([])
const total = ref(0)

// 初始化排序为时间倒序
orderBy.value = 'timestamp'
orderDirection.value = 'desc'

async function fetchCommonMessage() {
  loading.value = true
  try {
    const res = await getCommonMessage({
      search: search.value,
      user_ids: userIds.value,
      room_display_id: props.roomDisplayId,
      begin: timeRange.value[0],
      end: timeRange.value[1],
      page: currentPage.value,
      page_size: pageSize.value,
      order_by: orderBy.value,
      order_direction: orderDirection.value
    })
    messageList.value = res.data.data.list
    total.value = res.data.data.total
  } catch (error) {
    console.error('[CommonMessage] Fetch error:', error)
  } finally {
    loading.value = false
  }
}

function handleSort({ prop, order }: { prop: string; order: string }) {
  orderBy.value = prop
  orderDirection.value = order === 'ascending' ? 'asc' : 'desc'
  fetchCommonMessage()
}

onMounted(fetchCommonMessage)
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
          class="!w-full sm:!max-w-[565px]"
          clearable
          @keyup.enter="fetchCommonMessage"
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
          <el-button 
            type="primary"
            :loading="loading"
            class="!w-[240px]"
            @click="fetchCommonMessage"
          >
            查询
          </el-button>
        </div>
      </div>
    </div>

    <!-- 消息列表 -->
    <div class="overflow-x-auto">
      <el-table 
        :data="messageList" 
        class="w-full"
        size="small"
        @sort-change="handleSort"
      >
        <el-table-column 
          label="用户" 
          prop="user_name"
          sortable="custom"
          min-width="80"
        >
          <template #default="{ row }">
            <div class="flex flex-col">
              <span class="text-[11px] text-green-600 break-all">
                {{ row.user_name }}
              </span>
              <span class="text-[10px] text-gray-500">
                {{ row.user_display_id }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column 
          label="消息" 
          prop="message"
          min-width="200"
        >
          <template #default="{ row }">
            <span class="text-[11px] whitespace-normal break-words">{{ row.content }}</span>
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

      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :pager-count="5"
          size="small"
          layout="prev, pager, next"
          @current-change="fetchCommonMessage"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.el-table tr {
  cursor: pointer;
}

.el-table tr:hover > td {
  background-color: var(--el-table-row-hover-bg-color);
}
</style> 