package service

import (
	"context"
	"danmu-http/internal/model"
	"danmu-http/internal/validate"
	"danmu-http/logger"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
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
	start := time.Now()
	defer func() {
		logger.Info().
			Dur("total_duration", time.Since(start)).
			Str("room_id", req.RoomDisplayId).
			Int64("begin", req.Begin).
			Int64("end", req.End).
			Msg("ListGiftRanking completed")
	}()

	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// 获取总记录数
	total, err := model.GetGiftMessagesCount(req.ToUserIds, req.RoomDisplayId, req.Begin, req.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// 设置分页参数
	const pageSize = 2000
	pageCount := (total + pageSize - 1) / pageSize

	// 创建通道接收每个goroutine的处理结果
	resultChan := make(chan map[string]*UserGift, pageCount)
	errChan := make(chan error, pageCount)
	var wg sync.WaitGroup

	// 启动多个goroutine并行处理每一页数据
	for page := int64(1); page <= pageCount; page++ {
		wg.Add(1)
		go func(pageNum int64) {
			defer wg.Done()

			// 获取当前页数据
			messages, err := model.GetGiftMessagesByToUserIdTimestampRoomIdWithPage(
				req.ToUserIds,
				req.RoomDisplayId,
				req.Begin,
				req.End,
				int(pageNum),
				pageSize,
			)
			if err != nil {
				errChan <- fmt.Errorf("failed to get page %d: %w", pageNum, err)
				return
			}

			// 处理当前页数据
			userGiftMap := make(map[string]*UserGift)
			comboMap := make(map[string]*struct {
				currentCombo int64
				gift         *Gift
			})

			// 处理消息
			for _, msg := range messages {
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

			// 处理最后的连击
			for comboKey, combo := range comboMap {
				if combo.gift != nil {
					parts := strings.Split(comboKey, "_")
					userID, _ := strconv.ParseUint(parts[0], 10, 64)
					toUserID, _ := strconv.ParseUint(parts[1], 10, 64)
					userKey := fmt.Sprintf("%d_%d", userID, toUserID)

					userGift := userGiftMap[userKey]
					userGift.GiftList = append(userGift.GiftList, combo.gift)
					userGift.Total += combo.gift.DiamondCount * combo.gift.ComboCount
				}
			}

			resultChan <- userGiftMap
		}(page)
	}

	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// 合并所有页的结果
	finalMap := make(map[string]*UserGift)
	for pageResult := range resultChan {
		for key, gift := range pageResult {
			if existing, ok := finalMap[key]; ok {
				// 合并礼物列表和总值
				existing.GiftList = append(existing.GiftList, gift.GiftList...)
				existing.Total += gift.Total
			} else {
				finalMap[key] = gift
			}
		}
	}

	// 转换为切片并排序
	result := make([]*UserGift, 0, len(finalMap))
	for _, v := range finalMap {
		result = append(result, v)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})

	logger.Info().
		Int("total_records", int(total)).
		Int("page_count", int(pageCount)).
		Int("result_count", len(result)).
		Dur("total_duration", time.Since(start)).
		Msg("Parallel processing completed")

	return result, nil
}

func (s *giftMessageService) ListToUser(ctx context.Context, roomDisplayId string) ([]*model.ToUser, error) {
	return model.GetToUsersByRoomDisplayId(roomDisplayId)
}

func (s *giftMessageService) ListGiftMessagePageWithCondition(ctx context.Context, req *validate.GiftMessageQuery) ([]*model.GiftMessage, int64, error) {
	giftMessages, total, err := model.GetGiftMessageWithConditionPage(req)
	if err != nil {
		logger.Error().
			Err(err).
			Interface("query", req).
			Msg("failed to get gift messages")
		return nil, 0, fmt.Errorf("failed to get gift messages: %w", err)
	}
	return giftMessages, total, nil
}
