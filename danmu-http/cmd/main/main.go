package main

import (
	"danmu-http/internal/model"
	"danmu-http/logger"
	"danmu-http/router"
	"danmu-http/setting"
	"fmt"
	"log"
)

func main() {
	setting.Init()
	logger.Init()
	model.Init()
	router.Init()
	defer model.Close()

	r := router.SetupRouter()
	if err := r.Run(fmt.Sprintf(":%s", setting.AppSetting.Port)); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
