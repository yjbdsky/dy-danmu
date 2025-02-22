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
	// 添加超时控制
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 获取总记录数
	total, err := model.GetGiftMessagesCount(req.ToUserIds, req.RoomDisplayId, req.Begin, req.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// 设置分页和worker参数
	const (
		pageSize   = 5000
		maxWorkers = 5  // 控制并发数量
		bufferSize = 10 // 结果通道缓冲
	)
	pageCount := (total + pageSize - 1) / pageSize

	// 创建任务和结果通道
	type pageTask struct {
		pageNum int64
		ctx     context.Context
	}
	tasks := make(chan pageTask, pageCount)
	results := make(chan map[string]*UserGift, bufferSize)
	errors := make(chan error, 1)

	// 启动worker pool
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				// 检查上下文是否取消
				if task.ctx.Err() != nil {
					return
				}

				userGiftMap, err := s.processPage(req, task.pageNum, pageSize)
				if err != nil {
					select {
					case errors <- err:
					default:
					}
					return
				}

				select {
				case results <- userGiftMap:
				case <-task.ctx.Done():
					return
				}
			}
		}()
	}

	// 发送任务
	go func() {
		for page := int64(1); page <= pageCount; page++ {
			select {
			case tasks <- pageTask{pageNum: page, ctx: ctx}:
			case <-ctx.Done():
				return
			}
		}
		close(tasks)
	}()

	// 等待所有worker完成并关闭结果通道
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// 检查是否有错误
	select {
	case err := <-errors:
		if err != nil {
			return nil, err
		}
	default:
	}

	// 合并结果并排序
	return s.mergeAndSortResults(ctx, results)
}

// processPage 处理单个分页的数据
func (s *giftMessageService) processPage(req *validate.ListGiftRankingRequest, pageNum int64, pageSize int) (map[string]*UserGift, error) {
	messages, err := model.GetGiftMessagesByToUserIdTimestampRoomIdWithPage(
		req.ToUserIds,
		req.RoomDisplayId,
		req.Begin,
		req.End,
		int(pageNum),
		pageSize,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get page %d: %w", pageNum, err)
	}

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
				GiftList:        make([]*Gift, 0, 20), // 预分配容量
			}
			userGiftMap[userKey] = userGift
		}

		s.processGiftMessage(msg, userGift, comboMap)
	}

	// 处理剩余的连击
	s.processRemainingCombos(comboMap, userGiftMap)

	return userGiftMap, nil
}

// processGiftMessage 处理单个礼物消息
func (s *giftMessageService) processGiftMessage(msg *model.GiftMessage, userGift *UserGift, comboMap map[string]*struct {
	currentCombo int64
	gift         *Gift
}) {
	comboCount, _ := strconv.ParseInt(msg.ComboCount, 10, 64)
	if comboCount == 0 {
		comboCount = 1
	}

	gift := &Gift{
		GiftID:       msg.GiftID,
		GiftName:     msg.GiftName,
		DiamondCount: int64(msg.DiamondCount),
		ComboCount:   comboCount,
		Image:        msg.Image,
		Message:      msg.Message,
		Timestamp:    msg.Timestamp,
	}

	comboKey := fmt.Sprintf("%d_%d_%d", msg.UserID, msg.ToUserID, msg.GiftID)
	if combo, exists := comboMap[comboKey]; exists {
		if comboCount > combo.currentCombo {
			combo.currentCombo = comboCount
			combo.gift = gift
		} else {
			if combo.gift != nil {
				userGift.GiftList = append(userGift.GiftList, combo.gift)
				userGift.Total += combo.gift.DiamondCount * combo.gift.ComboCount
			}
			comboMap[comboKey] = &struct {
				currentCombo int64
				gift         *Gift
			}{
				currentCombo: comboCount,
				gift:         gift,
			}
		}
	} else {
		comboMap[comboKey] = &struct {
			currentCombo int64
			gift         *Gift
		}{
			currentCombo: comboCount,
			gift:         gift,
		}
	}
}

// processRemainingCombos 处理剩余的连击
func (s *giftMessageService) processRemainingCombos(comboMap map[string]*struct {
	currentCombo int64
	gift         *Gift
}, userGiftMap map[string]*UserGift) {
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
}

// mergeAndSortResults 合并并排序结果
func (s *giftMessageService) mergeAndSortResults(ctx context.Context, results chan map[string]*UserGift) ([]*UserGift, error) {
	finalMap := make(map[string]*UserGift)

	// 合并结果
	for pageResult := range results {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			for key, gift := range pageResult {
				if existing, ok := finalMap[key]; ok {
					existing.GiftList = append(existing.GiftList, gift.GiftList...)
					existing.Total += gift.Total
				} else {
					finalMap[key] = gift
				}
			}
		}
	}

	// 转换为切片
	result := make([]*UserGift, 0, len(finalMap))
	for _, v := range finalMap {
		result = append(result, v)
	}

	// 排序
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
		logger.Error().
			Err(err).
			Interface("query", req).
			Msg("failed to get gift messages")
		return nil, 0, fmt.Errorf("failed to get gift messages: %w", err)
	}
	return giftMessages, total, nil
}
