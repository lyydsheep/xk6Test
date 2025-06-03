package router

import (
	"email/logic/handler"
	"github.com/gin-gonic/gin"
)

func EventRouter(router *gin.RouterGroup, event *handler.EventHandler) {
	router.POST("/rt", event.CreateRTEvent)
}
