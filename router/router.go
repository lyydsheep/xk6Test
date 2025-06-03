package router

import (
	"email/logic/handler"
	"github.com/gin-gonic/gin"
	"os"
)

func RegisterRouter(event *handler.EventHandler, user *handler.UserHandler, limiter gin.HandlerFunc) *gin.Engine {
	switch os.Getenv("ENV") {
	case "prod":
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	server := r.Group("/api/v1")
	UserRouter(server.Group("/user"), limiter, user)
	EventRouter(server.Group("/events"), event)
	return r
}
