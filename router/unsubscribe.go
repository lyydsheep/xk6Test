package router

import (
	"email/logic/handler"
	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.RouterGroup, limiter gin.HandlerFunc, user *handler.UserHandler) {
	// 添加邮件退订路由
	// 限流
	// TODO 限流器需要依赖注入
	router.POST("/unsubscribe", limiter, user.UnsubscribeHandler)
}
