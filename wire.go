//go:build wireinject

package main

import (
	"email/dal/cache"
	"email/dal/dao"
	"email/logic/handler"
	"email/logic/repository"
	"email/logic/service"
	"email/middleware"
	"email/router"
	"github.com/google/wire"
)

const (
	UnSubscribeRate = 1
)

func InitializeApp() *App {
	wire.Build(NewApp,
		router.RegisterRouter, service.NewFlowControlService, service.NewEmailResultService, handler.NewEventHandler,
		service.NewEmailTaskSubService, service.NewEmailTaskService, service.NewEmailEventService, handler.NewUserHandler, middleware.RateLimiter, service.NewEventService,
		repository.NewEmailTaskRepoV1, repository.NewEmailTemplateRepoV1, repository.NewEmailSysUserRepoV1, repository.NewEmailUnsubscribeUsrUserRepoV1,
		repository.NewOpsEmlNumberSuccessRepoV1,
		repository.NewUsrUserRepoV1, repository.NewTaskSubRepoV1, repository.NewEmailRedemptionCodeV1, repository.NewEmailRedemptionTaskSubV1,
		repository.NewTaskConfigRepoV1, repository.NewEmailResultRepoV1,
		repository.NewEmailDomainCredibilityRepoV1,
		dao.DB, cache.NewMemCache, wire.Value(UnSubscribeRate))

	return nil
}
