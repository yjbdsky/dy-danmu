package validate

type CommonMessageQuery struct {
	Search        string   `json:"search" binding:"omitempty"`
	MessageType   []string `json:"message_type" binding:"omitempty"`
	UserIDs       []uint64 `json:"user_ids" binding:"omitempty"`
	RoomDisplayId string   `json:"room_display_id" binding:"required"`
	Begin         int64    `json:"begin" binding:"omitempty,min=1"`
	End           int64    `json:"end" binding:"omitempty,min=1"`
	PageRequest
}
