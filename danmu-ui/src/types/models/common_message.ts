//GET /api/common-message
export interface CommonMessageRequest {
  room_display_id: string;
  search: string;
  user_ids: number[];
  begin: number;
  end: number;
  order_by: string;
  order_direction: string;
  page: number;
  page_size: number;
}

export interface CommonMessage {
  id: number;
  user_id: number;
  user_name: string;
  user_display_id: string;
  room_display_id: string;
  content: string;
  timestamp: number;
}
