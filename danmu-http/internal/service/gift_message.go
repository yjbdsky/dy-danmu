package service

import (
	"context"
	"danmu-http/internal/model"
	"danmu-http/internal/validate"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type GiftMessageService interface {
	ListGiftRanking(ctx context.Context, req *validate.ListGiftRankingRequest) ([]*UserGift, error)
	ListToUser(ctx context.Context, roomDisplayId string) ([]*model.ToUser, error)
	ListGiftMessagePageWithCondition(ctx context.Context, req *validate.GiftMessageQuery) ([]*model.GiftMessage, int64, error)
}

type giftMessageService struct {
}

func NewGiftMessageService() GiftMessageService {
	return &giftMessageService{}
}

type UserGift struct {
	UserID          uint64  `json:"user_id"`
	UserName        string  `json:"user_name"`
	UserDisplayId   string  `json:"user_display_id"`
	Total           int64   `json:"total"`
	RoomDisplayId   string  `json:"room_display_id"`
	RoomName        string  `json:"room_name"`
	ToUserID        uint64  `json:"to_user_id"`
	ToUserName      string  `json:"to_user_name"`
	ToUserDisplayId string  `json:"to_user_display_id"`
	GiftList        []*Gift `json:"gift_list"`
}

type Gift struct {
	GiftID       int64  `json:"gift_id"`
	GiftName     string `json:"gift_name"`
	DiamondCount int64  `json:"diamond_count"`
	ComboCount   int64  `json:"combo_count"`
	Image        string `json:"image"`
	Message      string `json:"message"`
	Timestamp    int64  `json:"timestamp"`
}

func (s *giftMessageService) ListGiftRanking(ctx context.Context, req *validate.ListGiftRankingRequest) ([]*UserGift, error) {

	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	giftMessages, err := model.GetGiftMessagesByToUserIdTimestampRoomId(req.ToUserIds, req.RoomDisplayId, req.Begin, req.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get gift messages: %w", err)
	}

	// 用map来统计每个用户的礼物信息
	userGiftMap := make(map[string]*UserGift)
	comboMap := make(map[string]*struct {
		currentCombo int64
		gift         *Gift
	})

	for _, msg := range giftMessages {
		// Create composite key for user-toUser pair
		userKey := fmt.Sprintf("%d_%d", msg.UserID, msg.ToUserID)

		// 获取或创建用户礼物统计
		userGift, exists := userGiftMap[userKey]
		if !exists {
			userGift = &UserGift{
				UserID:          msg.UserID,
				UserName:        msg.UserName,
				UserDisplayId:   msg.UserDisplayId,
				RoomDisplayId:   msg.RoomDisplayId,
				RoomName:        msg.RoomName,
				ToUserID:        msg.ToUserID,
				ToUserName:      msg.ToUserName,
				ToUserDisplayId: msg.ToUserDisplayId,
				GiftList:        make([]*Gift, 0),
			}
			userGiftMap[userKey] = userGift
		}

		// 转换comboCount从string到int64
		comboCount, err := strconv.ParseInt(msg.ComboCount, 10, 64)
		if err != nil {
			comboCount = 1
		}

		// 创建礼物对象
		gift := &Gift{
			GiftID:       msg.GiftID,
			GiftName:     msg.GiftName,
			DiamondCount: int64(msg.DiamondCount),
			ComboCount:   comboCount,
			Image:        msg.Image,
			Message:      msg.Message,
			Timestamp:    msg.Timestamp,
		}

		// Update combo key to include ToUserID
		comboKey := fmt.Sprintf("%d_%d_%d", msg.UserID, msg.ToUserID, msg.GiftID)
		combo, exists := comboMap[comboKey]

		if exists {
			if comboCount > combo.currentCombo {
				// 连击数增加，更新当前连击信息
				combo.currentCombo = comboCount
				combo.gift = gift
			} else {
				// 连击中断，保存之前的峰值
				if combo.gift != nil {
					userGift.GiftList = append(userGift.GiftList, combo.gift)
					userGift.Total += combo.gift.DiamondCount * combo.gift.ComboCount
				}
				// 开始新的连击序列
				comboMap[comboKey] = &struct {
					currentCombo int64
					gift         *Gift
				}{
					currentCombo: comboCount,
					gift:         gift,
				}
			}
		} else {
			// 新的连击序列
			comboMap[comboKey] = &struct {
				currentCombo int64
				gift         *Gift
			}{
				currentCombo: comboCount,
				gift:         gift,
			}
		}
	}

	// 处理最后的连击礼物
	for comboKey, combo := range comboMap {
		if combo.gift != nil {
			// Parse the composite key (format: "userID_toUserID_giftID")
			parts := strings.Split(comboKey, "_")
			userID, _ := strconv.ParseUint(parts[0], 10, 64)
			toUserID, _ := strconv.ParseUint(parts[1], 10, 64)
			userKey := fmt.Sprintf("%d_%d", userID, toUserID)

			userGift := userGiftMap[userKey]
			userGift.GiftList = append(userGift.GiftList, combo.gift)
			userGift.Total += combo.gift.DiamondCount * combo.gift.ComboCount
		}
	}

	// Convert map to slice
	result := make([]*UserGift, 0, len(userGiftMap))
	for _, v := range userGiftMap {
		result = append(result, v)
	}

	// 按照total降序排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})

	return result, nil
}

func (s *giftMessageService) ListToUser(ctx context.Context, roomDisplayId string) ([]*model.ToUser, error) {
	return model.GetToUsersByRoomDisplayId(roomDisplayId)
}

func (s *giftMessageService) ListGiftMessagePageWithCondition(ctx context.Context, req *validate.GiftMessageQuery) ([]*model.GiftMessage, int64, error) {
	giftMessages, total, err := model.GetGiftMessageWithConditionPage(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get gift messages: %w", err)
	}

	return giftMessages, total, nil
}
