package validate

type LiveConfAddRequest struct {
	RoomDisplayID string `json:"room_display_id" binding:"required"`
	URL           string `json:"url" binding:"required,url"`
	Name          string `json:"name" binding:"required"`
	Enable        bool   `json:"enable" binding:"required"`
}

type LiveConfUpdateRequest struct {
	ID            int64  `json:"id" binding:"required"`
	RoomDisplayID string `json:"room_display_id" binding:"required"`
	URL           string `json:"url" binding:"required,url"`
	Name          string `json:"name" binding:"required"`
	Enable        bool   `json:"enable" binding:"omitempty"`
}
