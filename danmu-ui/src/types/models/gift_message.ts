//GET /api/gift-message/ranking
export interface GiftMessageRankingRequest {
  to_user_ids: number[];
  room_display_id: string;
  begin: number;
  end: number;
}

export interface UserGift {
  user_id: number;
  user_name: string;
  user_display_id: string;
  total: number;
  room_display_id: string;
  room_name: string;
  to_user_id: number;
  to_user_name: string;
  to_user_display_id: string;
  gift_list: Gift[];
}

export interface Gift {
  gift_id: number;
  gift_name: string;
  diamond_count: number;
  combo_count: number;
  image: string;
  message: string;
  timestamp: number;
}

//GET /api/gift-message
export interface GiftMessageRequest {
  search: string;
  user_ids: number[];
  to_user_ids: number[];
  room_display_id: string;
  begin: number;
  end: number;
  diamond_count: number;
  page: number;
  page_size: number;
  order_by: string;
  order_direction: string;
}

export interface GiftMessage {
  id: number;
  user_id: number;
  user_name: string;
  user_display_id: string;
  to_user_id: number;
  to_user_name: string;
  to_user_display_id: string;
  gift_id: number;
  gift_name: string;
  diamond_count: number;
  combo_count: number;
  message: string;
  timestamp: number;
}

