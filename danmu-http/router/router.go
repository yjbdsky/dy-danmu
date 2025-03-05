package router

import (
	"danmu-http/internal/handler"
	"danmu-http/internal/service"
	"danmu-http/middleware"
	"github.com/gin-gonic/gin"
)

var (
	authHandler          *handler.AuthHandler
	liveConfHandler      *handler.LiveConfHandler
	giftMessageHandler   *handler.GiftMessageHandler
	commonMessageHandler *handler.CommonMessageHandler
	userHandler          *handler.UserHandler
)

func Init() {
	authHandler = handler.NewAuthHandler(service.NewAuthService())
	liveConfHandler = handler.NewLiveConfHandler(service.NewLiveConfService())
	giftMessageHandler = handler.NewGiftMessageHandler(service.NewGiftMessageService())
	commonMessageHandler = handler.NewCommonMessageHandler(service.NewCommonMessageService())
	userHandler = handler.NewUserHandler(service.NewUserService())

}

func SetupRouter() *gin.Engine {
	r := gin.New()
	//r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	InitResourceServer(r)

	html := NewHtmlHandler()
	r.GET("/", html.Index)
	r.NoRoute(html.RedirectIndex)

	//InitResource(r)
	api := r.Group("/api")
	{
		// 不需要认证的路由
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.JWT())
		{
			// Auth 相关路由
			auth := authenticated.Group("/auth")
			{
				// 管理员权限
				adminAuth := auth.Group("")
				adminAuth.Use(middleware.AdminRequired())
				{
					adminAuth.POST("/register", authHandler.Register)
					adminAuth.POST("/reset-password/:id", authHandler.ResetPassword)
					adminAuth.DELETE("/:id", authHandler.Delete)
				}

				// 所有认证用户
				auth.PUT("/self", authHandler.UpdateSelf)
				auth.GET("/self", authHandler.GetSelf)
				auth.GET("/list", authHandler.ListAll)
			}

			// LiveConf 相关路由
			liveConf := authenticated.Group("/live-conf")
			{
				// 管理员权限
				adminLiveConf := liveConf.Group("")
				adminLiveConf.Use(middleware.AdminRequired())
				{
					adminLiveConf.POST("", liveConfHandler.Create)
					adminLiveConf.PUT("", liveConfHandler.Update)
					adminLiveConf.DELETE("/:id", liveConfHandler.Delete)
				}

				// 所有认证用户
				liveConf.GET("/:id", liveConfHandler.Get)
				liveConf.GET("", liveConfHandler.List)
			}

			// GiftMessage 相关路由
			giftMessage := authenticated.Group("/gift-message")
			{
				giftMessage.POST("/ranking", giftMessageHandler.ListGiftRanking)
				giftMessage.GET("/to-user", giftMessageHandler.ListToUser)
				giftMessage.POST("", giftMessageHandler.ListGiftMessagePageWithCondition)
			}

			// CommonMessage 相关路由
			commonMessage := authenticated.Group("/common-message")
			{
				commonMessage.POST("", commonMessageHandler.ListPageableWithCondition)
			}

			// User 相关路由
			user := authenticated.Group("/user")
			{
				user.GET("", userHandler.ListAllUsers)
				user.GET("/search", userHandler.SearchUser)
			}
		}
	}

	return r
}
