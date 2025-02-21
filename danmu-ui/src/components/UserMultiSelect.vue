<script setup lang="ts">
    import { reactive, toRefs, onMounted } from "vue";
    import { searchUser } from '../api/user'
    import type { User } from '../types/models/user'
    import { Search, Refresh } from '@element-plus/icons-vue'

    const props = defineProps<{
        width?: string
        modelValue: number[]
    }>()

    const emit = defineEmits<{
        (e: 'update:modelValue', value: number[]): void
    }>()

    const state = reactive({
        userList: [] as User[],
        total: 0,
        selectedUserIds: props.modelValue || [],
        currentPage: 1,
        pageSize: 200,
        keyword: ''
    });

    const { userList, total, selectedUserIds, currentPage, pageSize, keyword } = toRefs(state);

    async function fetchUsers(page = currentPage.value) {
        try {
            const res = await searchUser({
                keyword: keyword.value,
                page,
                page_size: pageSize.value
            })
            state.userList = res.data.data.list
            state.total = res.data.data.total
            state.currentPage = page
        } catch (error) {
            console.error('[Users] Fetch error:', error)
        }
    }

    function handleSearch() {
        fetchUsers(1)  // 搜索时重置到第一页
    }

    function handlePageChange(newPage: number) {
        fetchUsers(newPage)
    }

    function resetAndFetch() {
        // 重置所有状态到初始值
        state.keyword = ''
        state.currentPage = 1
        state.selectedUserIds = []
        emit('update:modelValue', [])  // 通知父组件清空选择
        
        // 重新获取数据
        fetchUsers(1)
    }

    onMounted(fetchUsers)
</script>

<template>
    <div>
        <el-select
            v-model="selectedUserIds"
            multiple
            collapse-tags
            collapse-tags-tooltip
            :style="{ width: width || '240px' }"
            :popper-class="`width-${width || '240px'} user-select-dropdown`"
            placeholder="请选择用户"
            :filterable="false"
            @change="(newValue: number[]) => emit('update:modelValue', newValue)"
        >
            <template #prefix>
                <el-button
                    type="primary"
                    link
                    :icon="Refresh"
                    @click.stop="resetAndFetch"
                />
            </template>

            <div class="relative max-h-[300px]">
                <!-- 固定在顶部的搜索框 -->
                <div class="sticky top-0 left-0 right-0 bg-white border-b p-2 z-10">
                    <div class="flex gap-2">
                        <el-input
                            v-model="keyword"
                            placeholder="搜索用户"
                            class="!min-w-0"
                            @keyup.enter="handleSearch"
                        >
                            <template #append>
                                <el-button @click="handleSearch">
                                    <el-icon><Search /></el-icon>
                                </el-button>
                            </template>
                        </el-input>
                    </div>
                </div>

                <!-- 选项列表 -->
                <div class="overflow-y-auto" style="height: calc(300px - 80px)">
                    <el-option
                        v-for="item in userList"
                        :key="item.id"
                        :label="`${item.user_name}`"
                        :value="item.user_id"
                    >
                        <div class="flex items-center space-x-1 px-2 py-1">
                            <el-tag 
                                type="info" 
                                size="large"
                                class="!w-[120px] !truncate !text-center !m-0"
                            >
                                {{ item.display_id }}
                            </el-tag>
                            <el-tag 
                                type="success" 
                                size="large"
                                class="!w-[140px] !truncate !text-center !m-0"
                            >
                                {{ item.user_name }}
                            </el-tag>
                        </div>
                    </el-option>
                    
                    <!-- 无数据提示 -->
                    <div v-if="userList.length === 0" class="p-4 text-center text-gray-500">
                        暂无数据
                    </div>
                </div>

                <!-- 固定在底部的分页器 -->
                <div class="sticky bottom-0 left-0 right-0 bg-white border-t">
                    <el-pagination
                        :current-page="currentPage"
                        :page-size="pageSize"
                        :total="total"
                        :pager-count="5"
                        size="small"
                        class="!py-2"
                        layout="prev, pager, next"
                        @current-change="handlePageChange"
                    />
                </div>
            </div>
        </el-select>
    </div>
</template>


<style>
.width-240px {
  min-width: 280px;  /* 增加宽度以适应两个 tag */
}

.user-select-dropdown .el-select-dropdown__wrap {
  max-height: none;
  padding: 0;
}

.user-select-dropdown .el-pagination {
  justify-content: center;
  --el-pagination-button-width: 24px;  /* 减小按钮宽度 */
}

.user-select-dropdown .el-pager li {
  min-width: 24px;  /* 减小页码按钮宽度 */
}

.user-select-dropdown .el-select-dropdown__list {
  padding: 0;
}
</style>

