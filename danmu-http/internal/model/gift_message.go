package model

import "danmu-http/internal/validate"

const TableNameGiftMessage = "gift_messages"

// GiftMessage mapped from table <gift_messages>
type GiftMessage struct {
	ID              int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID          uint64 `gorm:"column:user_id;not null" json:"user_id"`                       // User ID who sent the gift
	UserName        string `gorm:"column:user_name;not null" json:"user_name"`                   // User Name
	UserDisplayId   string `gorm:"column:user_display_id;not null" json:"user_display_id"`       // User Display ID
	ToUserID        uint64 `gorm:"column:to_user_id;not null" json:"to_user_id"`                 // To User ID
	ToUserName      string `gorm:"column:to_user_name;not null" json:"to_user_name"`             // To User Name
	ToUserDisplayId string `gorm:"column:to_user_display_id;not null" json:"to_user_display_id"` // To User Display ID
	GiftName        string `gorm:"column:gift_name;not null" json:"gift_name"`                   // Gift ID (could be a foreign key)
	GiftID          int64  `gorm:"column:gift_id;not null" json:"gift_id"`
	RoomID          int64  `gorm:"column:room_id;not null" json:"room_id"`                 // Room ID
	RoomDisplayId   string `gorm:"column:room_display_id;not null" json:"room_display_id"` // Room Display ID
	RoomName        string `gorm:"column:room_name;not null" json:"room_name"`             // Room Name
	Message         string `gorm:"column:message;not null" json:"message"`                 // The gift message
	Timestamp       int64  `gorm:"column:timestamp;not null" json:"timestamp"`
	DiamondCount    uint32 `gorm:"column:diamond_count;not null" json:"diamond_count"`
	Image           string `gorm:"column:image_url" json:"image_url"`
	ComboCount      string `gorm:"column:combo_count" json:"combo_count"`
}

type ToUser struct {
	ToUserID        uint64 `gorm:"column:to_user_id" json:"user_id"`
	ToUserName      string `gorm:"column:to_user_name" json:"user_name"`
	ToUserDisplayId string `gorm:"column:to_user_display_id" json:"display_id"`
}

func (*GiftMessage) TableName() string {
	return TableNameGiftMessage
}

func GetGiftMessagesByToUserIdTimestampRoomId(toUserIds []uint64, roomDisplayId string, begin int64, end int64) ([]*GiftMessage, error) {
	db := DB.Model(&GiftMessage{})
	if len(toUserIds) > 0 {
		db = db.Where("to_user_id IN (?)", toUserIds)
	}
	if begin != 0 {
		db = db.Where("timestamp >= ?", begin)
	}
	if end != 0 {
		db = db.Where("timestamp <= ?", end)
	}
	if roomDisplayId != "" {
		db = db.Where("room_display_id = ?", roomDisplayId)
	}
	db = db.Order("timestamp ASC")
	var giftMessages []*GiftMessage
	err := db.Find(&giftMessages).Error
	if err != nil {
		return nil, err
	}
	return giftMessages, nil
}

func GetToUsersByRoomDisplayId(roomDisplayId string) ([]*ToUser, error) {
	var toUsers []*ToUser
	err := DB.Model(&GiftMessage{}).
		Where("room_display_id = ? AND to_user_id != ?", roomDisplayId, 0).
		Select("DISTINCT ON (to_user_id) to_user_id, to_user_name, to_user_display_id").
		Find(&toUsers).Error
	if err != nil {
		return nil, err
	}
	return toUsers, nil
}

func GetGiftMessageWithConditionPage(req *validate.GiftMessageQuery) ([]*GiftMessage, int64, error) {
	var giftMessages []*GiftMessage
	var total int64

	db := DB.Model(&GiftMessage{})

	// Build query conditions
	if len(req.UserIDs) > 0 {
		db = db.Where("user_id IN (?)", req.UserIDs)
	}
	if len(req.ToUserIds) > 0 {
		db = db.Where("to_user_id IN (?)", req.ToUserIds)
	}
	if req.Begin != 0 {
		db = db.Where("timestamp >= ?", req.Begin)
	}
	if req.End != 0 {
		db = db.Where("timestamp <= ?", req.End)
	}
	if req.DiamondCount != 0 {
		db = db.Where("diamond_count >= ?", req.DiamondCount)
	}
	if req.RoomDisplayId != "" {
		db = db.Where("room_display_id = ?", req.RoomDisplayId)
	}
	if req.Search != "" {
		db = db.Where("message LIKE ?", "%"+req.Search+"%")
	}

	if req.OrderBy == "" {
		req.OrderBy = "timestamp"
		req.OrderDirection = "desc"
	}
	// Get total count before pagination
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply ordering and pagination
	err = db.Order(req.OrderBy + " " + req.OrderDirection).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&giftMessages).Error
	if err != nil {
		return nil, 0, err
	}

	return giftMessages, total, nil
}

// 添加分页查询方法
func GetGiftMessagesByToUserIdTimestampRoomIdWithPage(toUserIds []uint64, roomDisplayId string, begin, end int64, page, pageSize int) ([]*GiftMessage, error) {
	db := DB.Model(&GiftMessage{})
	if roomDisplayId != "" {
		db = db.Where("room_display_id = ?", roomDisplayId)
	}

	if len(toUserIds) > 0 {
		db = db.Where("to_user_id IN (?)", toUserIds)
	}
	if begin != 0 && end != 0 {
		db = db.Where("timestamp BETWEEN ? AND ?", begin, end)
	}
	db = db.Order("timestamp ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	var messages []*GiftMessage
	if err := db.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// 添加获取总数的方法
func GetGiftMessagesCount(toUserIds []uint64, roomDisplayId string, begin, end int64) (int64, error) {
	var count int64
	db := DB.Model(&GiftMessage{})
	if roomDisplayId != "" {
		db = db.Where("room_display_id = ?", roomDisplayId)
	}

	if len(toUserIds) > 0 {
		db = db.Where("to_user_id IN (?)", toUserIds)
	}
	if begin != 0 && end != 0 {
		db = db.Where("timestamp BETWEEN ? AND ?", begin, end)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
