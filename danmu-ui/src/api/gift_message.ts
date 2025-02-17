import api from "../plugins/axios"
import type { GiftMessageRankingRequest, GiftMessageRequest, UserGift, GiftMessage } from "../types/models/gift_message"
import type { User } from "../types/models/user"
import type { ApiResponse, Page } from "../types/response"

export function getGiftRanking(data: GiftMessageRankingRequest) {
    return api.post<ApiResponse<UserGift[]>>('/api/gift-message/ranking',  data )
}

export function getToUsers(roomDisplayId: string) {
    return api.get<ApiResponse<User[]>>('/api/gift-message/to-user', {
        params: { room_display_id: roomDisplayId }
    })
}
  
export function getGiftMessage(data: GiftMessageRequest) {
    return api.post<ApiResponse<Page<GiftMessage>>>('/api/gift-message',  data )
}

