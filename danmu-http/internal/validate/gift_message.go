package validate

type ListGiftRankingRequest struct {
	ToUserIds     []uint64 `json:"to_user_ids" binding:"omitempty"`
	RoomDisplayId string   `json:"room_display_id" binding:"required"`
	Begin         int64    `json:"begin" binding:"required,min=1"`
	End           int64    `json:"end" binding:"required,min=1"`
}

type GiftMessageQuery struct {
	Search        string   `json:"search" binding:"omitempty"`
	UserIDs       []uint64 `json:"user_ids" binding:"omitempty"`
	ToUserIds     []uint64 `json:"to_user_ids" binding:"omitempty"`
	RoomDisplayId string   `json:"room_display_id" binding:"required"`
	Begin         int64    `json:"begin" binding:"omitempty,min=1"`
	End           int64    `json:"end" binding:"omitempty,min=1"`
	DiamondCount  int64    `json:"diamond_count" binding:"omitempty,min=0"`
	PageRequest
}
