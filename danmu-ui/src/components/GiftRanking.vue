<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getGiftRanking } from '../api/gift_message'
import type { UserGift } from '../types/models/gift_message'
import TimeRangePicker from './TimeRangePicker.vue'
import ToUserSelect from './ToUserSelect.vue'
import dayjs from 'dayjs'

const props = defineProps<{
  roomDisplayId: string
}>()

// 初始化为今日的开始和结束时间
const timeRange = ref<[number, number]>([
  dayjs().startOf('day').valueOf(),
  dayjs().endOf('day').valueOf()
])
const selectedToUserIds = ref<number[]>([])
const giftList = ref<UserGift[]>([])
const loading = ref(false)

async function fetchGiftRanking() {
  loading.value = true
  try {
    const res = await getGiftRanking({
      room_display_id: props.roomDisplayId,
      to_user_ids: selectedToUserIds.value.map(Number),
      begin: timeRange.value[0],
      end: timeRange.value[1]
    })
    giftList.value = res.data.data
  } catch (error) {
    console.error('[GiftRanking] Fetch error:', error)
  } finally {
    loading.value = false
  }
}

// 组件挂载时自动查询
onMounted(fetchGiftRanking)

function formatTime(timestamp: number) {
  return dayjs(timestamp).format('YYYY-MM-DD HH:mm:ss')
}
</script>

<template>
  <div class="space-y-4">
    <!-- 筛选条件和查询按钮 -->
    <div class="flex flex-wrap items-end gap-4">
      <TimeRangePicker v-model="timeRange" />
      <ToUserSelect 
        v-model="selectedToUserIds"
        :room-display-id="roomDisplayId"
      />
      <el-button 
        type="primary"
        :loading="loading"
        @click="fetchGiftRanking"
      >
        查询
      </el-button>
    </div>

    <!-- 礼物排行列表 -->
    <div class="overflow-x-auto">
      <el-table 
        :data="giftList" 
        class="w-full"
      >
        <el-table-column type="expand" width="20">
          <template #default="{ row }">
            <div class="p-4 bg-gray-50">
              <el-table :data="row.gift_list" class="w-full">
                <el-table-column label="礼物信息" min-width="200">
                  <template #default="{ row: gift }">
                    <div class="flex items-center gap-2">
                      <el-image
                        :src="gift.image"
                        class="w-8 h-8 object-cover rounded"
                        fit="cover"
                      />
                      <span>{{ gift.message }}</span>
                    </div>
                  </template>
                </el-table-column>
                
                <el-table-column label="价值" min-width="100" align="right">
                  <template #default="{ row: gift }">
                    <div class="flex flex-col items-end">
                      <span class="text-orange-500">
                        {{ gift.diamond_count * gift.combo_count }}
                      </span>
                      <span class="text-gray-400 text-xs">
                        {{ gift.diamond_count }}×{{ gift.combo_count }}
                      </span>
                    </div>
                  </template>
                </el-table-column>

                <el-table-column label="时间" min-width="160">
                  <template #default="{ row: gift }">
                    {{ formatTime(gift.timestamp) }}
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </template>
        </el-table-column>

        <!-- 排名 -->
        <el-table-column width="22" align="center">
          <template #default="{ $index }">
            <template v-if="$index < 50">
              <span 
                class="text-[11px] font-medium inline-block leading-none"
                :class="{
                  'text-yellow-500': $index === 0,
                  'text-gray-500': $index === 1,
                  'text-orange-500': $index === 2,
                  'text-blue-500': $index > 2
                }"
              >
                {{ $index + 1 }}
              </span>
            </template>
          </template>
        </el-table-column>

        <el-table-column label="送礼用户" min-width="90">
          <template #default="{ row }">
            <div class="flex flex-col">
              <span class="text-[11px] text-green-600 break-all">
                {{ row.user_name }}
              </span>
              <el-tag size="small" type="info" class="!text-xs">
                {{ row.user_display_id }}
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="接收用户" min-width="90">
          <template #default="{ row }">
            <div class="flex flex-col">
              <span class="text-[11px] text-blue-600 break-all">
                {{ row.to_user_name }}
              </span>
              <el-tag 
                v-if="row.to_user_display_id" 
                size="small" 
                type="info" 
                class="!text-xs"
              >
                {{ row.to_user_display_id }}
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="金额" min-width="68" align="right">
          <template #default="{ row }">
            <span class="text-orange-500 font-medium text-[11px] tabular-nums">{{ row.total }}</span>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<style scoped>
.el-table :deep(.el-table__expanded-cell) {
  padding: 0;
}

@media (max-width: 640px) {
  .el-table {
    font-size: 14px;
  }
}

.el-table tr {
  cursor: pointer;
}

.el-table tr:hover > td {
  background-color: var(--el-table-row-hover-bg-color);
}
</style> 