import api from "../plugins/axios"
import type { LiveConf, CreateLiveConfRequest, UpdateLiveConfRequest } from "../types/models/live_conf"
import type { ApiResponse } from "../types/response"

interface LiveConfListResponse {
  list: LiveConf[]
}

export function getLiveConfList() {
  return api.get<ApiResponse<LiveConfListResponse>>('/api/live-conf')
}

export function getLiveConf(id: string) {
  return api.get<ApiResponse<LiveConf>>(`/api/live-conf/${id}`)
}

export function createLiveConf(data: CreateLiveConfRequest) {
  return api.post<ApiResponse<null>>('/api/live-conf', data)
}

export function updateLiveConf(data: UpdateLiveConfRequest) {
  return api.put<ApiResponse<null>>(`/api/live-conf`, data)
}

export function deleteLiveConf(id: string) {
  return api.delete<ApiResponse<null>>(`/api/live-conf/${id}`)
}
