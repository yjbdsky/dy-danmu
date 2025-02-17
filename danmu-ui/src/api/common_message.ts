import api from "../plugins/axios"
import type { CommonMessageRequest, CommonMessage } from "../types/models/common_message"
import type { ApiResponse, Page } from "../types/response"

export function getCommonMessage(data: CommonMessageRequest) {
    return api.post<ApiResponse<Page<CommonMessage>>>('/api/common-message', data )
}

