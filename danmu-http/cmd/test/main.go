package main

import (
	"context"
	"danmu-http/internal/model"
	"danmu-http/internal/service"
	"danmu-http/internal/validate"
	"danmu-http/setting"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	setting.Init()
	model.Init()
	giftservice := service.NewGiftMessageService()
	request := &validate.ListGiftRankingRequest{
		ToUserIds:     []uint64{111618450414},
		RoomDisplayId: "949735304546",
		Begin:         1739683802000,
		End:           1739701802000,
	}
	result, err := giftservice.ListGiftRanking(context.Background(), request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	/*	tousers, err := giftservice.ListToUser(context.Background(), "949735304546")
		fmt.Println(tousers)*/
	file, err := ExportToExcel(result)
	if err != nil {
		// 处理错误
	}
	err = file.SaveAs("gift_ranking.xlsx")
	if err != nil {
		// 处理错误
	}
}

func ExportToExcel(data []*service.UserGift) (*excelize.File, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// 删除默认的 Sheet1
	f.DeleteSheet("Sheet1")

	// 创建主表
	mainSheet := "打赏排行榜"
	f.NewSheet(mainSheet)

	// 设置主表表头
	headers := []string{"排名", "用户名", "房间", "主播", "总金额(钻石)"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(mainSheet, cell, header)
	}

	// 设置表头样式
	style, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DCE6F1"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create style error: %w", err)
	}
	f.SetRowStyle(mainSheet, 1, 1, style)

	// 填充主表数据
	for i, gift := range data {
		row := i + 2
		f.SetCellValue(mainSheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(mainSheet, fmt.Sprintf("B%d", row), gift.UserName)
		f.SetCellValue(mainSheet, fmt.Sprintf("C%d", row), gift.RoomName)
		f.SetCellValue(mainSheet, fmt.Sprintf("D%d", row), gift.ToUserName)
		f.SetCellValue(mainSheet, fmt.Sprintf("E%d", row), gift.Total)

		// 为每个用户创建详情表
		detailSheet := fmt.Sprintf("%s的礼物详情", gift.UserName)
		f.NewSheet(detailSheet)

		// 设置详情表表头
		detailHeaders := []string{"礼物名称", "钻石单价", "数量", "总价值", "赠送时间", "留言"}
		for j, header := range detailHeaders {
			cell := fmt.Sprintf("%c1", 'A'+j)
			f.SetCellValue(detailSheet, cell, header)
		}
		f.SetRowStyle(detailSheet, 1, 1, style)

		// 填充详情表数据
		for j, giftDetail := range gift.GiftList {
			row := j + 2
			f.SetCellValue(detailSheet, fmt.Sprintf("A%d", row), giftDetail.GiftName)
			f.SetCellValue(detailSheet, fmt.Sprintf("B%d", row), giftDetail.DiamondCount)
			f.SetCellValue(detailSheet, fmt.Sprintf("C%d", row), giftDetail.ComboCount)
			f.SetCellValue(detailSheet, fmt.Sprintf("D%d", row), giftDetail.DiamondCount*giftDetail.ComboCount)
			// 转换时间戳为可读时间
			timeStr := time.Unix(giftDetail.Timestamp/1000, 0).Format("2006-01-02 15:04:05")
			f.SetCellValue(detailSheet, fmt.Sprintf("E%d", row), timeStr)
			f.SetCellValue(detailSheet, fmt.Sprintf("F%d", row), giftDetail.Message)
		}

		// 设置列宽
		f.SetColWidth(detailSheet, "A", "A", 15)
		f.SetColWidth(detailSheet, "B", "D", 12)
		f.SetColWidth(detailSheet, "E", "E", 20)
		f.SetColWidth(detailSheet, "F", "F", 30)
	}

	// 设置主表列宽
	f.SetColWidth(mainSheet, "A", "A", 8)
	f.SetColWidth(mainSheet, "B", "D", 15)
	f.SetColWidth(mainSheet, "E", "E", 15)

	// 设置第一个sheet为活动sheet
	index, err := f.GetSheetIndex(mainSheet)
	if err != nil {
		return nil, fmt.Errorf("get sheet index error: %w", err)
	}
	f.SetActiveSheet(index)

	return f, nil
}
