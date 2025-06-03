package main

import (
	"context"
	"email/common/log"
	"email/config"
	"email/dal/dao"
	"email/logic/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func init() {
	config.InitConfig()
	log.InitLogger()
	//dao.InitGormLogger()
	dao.InitDB()
}

type App struct {
	Task           *service.EmailTaskService
	TaskSub        *service.EmailTaskSubService
	EmailEvent     *service.EmailEventService
	server         *gin.Engine
	FlowController *service.AliFlowControlService
	EmailResult    *service.EmailResultService
}

func (app *App) Run() {
	ctx := context.Background()
	go app.Task.ProductTaskSub(ctx)
	// 停止发送
	go app.TaskSub.Run(ctx)
	go app.EmailEvent.Start(ctx)
	go app.FlowController.Run(ctx)
	go app.EmailResult.Run(ctx)
	app.server.Run(fmt.Sprintf("0.0.0.0:%s", config.App.Port))
	select {}
}

func NewApp(task *service.EmailTaskService, taskSub *service.EmailTaskSubService, emailEvent *service.EmailEventService,
	server *gin.Engine, flowController *service.AliFlowControlService, emailResult *service.EmailResultService) *App {
	return &App{
		Task:           task,
		TaskSub:        taskSub,
		EmailEvent:     emailEvent,
		server:         server,
		FlowController: flowController,
		EmailResult:    emailResult,
	}
}

func main() {
	app := InitializeApp()
	app.Run()
}
