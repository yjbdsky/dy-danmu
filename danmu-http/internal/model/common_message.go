package model

import "danmu-http/internal/validate"

const TableNameCommonMessage = "common_messages"

// CommonMessage mapped from table <common_messages>
type CommonMessage struct {
	ID            uint64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MessageType   string `gorm:"column:message_type;not null" json:"message_type"`
	RoomID        int64  `gorm:"column:room_id;not null" json:"room_id"`
	RoomDisplayId string `gorm:"column:room_display_id;not null" json:"room_display_id"`
	RoomName      string `gorm:"column:room_name;not null" json:"room_name"`
	UserName      string `gorm:"column:user_name;not null" json:"user_name"`
	UserID        uint64 `gorm:"column:user_id;not null" json:"user_id"`
	UserDisplayId string `gorm:"column:user_display_id;not null" json:"user_display_id"`
	Content       string `gorm:"column:content;not null" json:"content"`
	Timestamp     int64  `gorm:"column:timestamp;not null" json:"timestamp"`
}

// TableName CommonMessage's table name
func (*CommonMessage) TableName() string {
	return TableNameCommonMessage
}

func GetCommonMessageWithConditionPage(req *validate.CommonMessageQuery) ([]*CommonMessage, int64, error) {
	var commonMessages []*CommonMessage
	var total int64

	db := DB.Model(&CommonMessage{})

	if len(req.MessageType) > 0 {
		db = db.Where("message_type IN (?)", req.MessageType)
	}

	if len(req.UserIDs) > 0 {
		db = db.Where("user_id IN (?)", req.UserIDs)
	}

	if req.RoomDisplayId != "" {
		db = db.Where("room_display_id = ?", req.RoomDisplayId)
	}

	if req.Begin != 0 {
		db = db.Where("timestamp >= ?", req.Begin)
	}

	if req.End != 0 {
		db = db.Where("timestamp <= ?", req.End)
	}

	if req.Search != "" {
		db = db.Where("content LIKE ?", "%"+req.Search+"%")
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if req.OrderBy == "" {
		req.OrderBy = "timestamp"
		req.OrderDirection = "desc"
	}

	err = db.Order(req.OrderBy + " " + req.OrderDirection).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&commonMessages).Error
	if err != nil {
		return nil, 0, err
	}

	return commonMessages, total, nil
}
