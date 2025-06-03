package handler

import (
	"email/common/log"
	"email/common/request"
	"email/logic/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var (
	validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())
)

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{
		svc: svc,
	}
}

func (h *EventHandler) CreateRTEvent(ctx *gin.Context) {
	// 校验参数
	var req request.CreateRTEventReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.New(ctx).Error("failed by binding json", "error", err)
		parameterErr(ctx, "failed by binding json "+err.Error())
		return
	}
	if err := validate.Struct(req); err != nil {
		log.New(ctx).Error("failed by checking parameter format", "error", err)
		parameterErr(ctx, "failed by checking parameter format "+err.Error())
		return
	}
	// 校验模板 ID 和 Cid 是否正确
	// 模板类型必须是 real_time
	// 用户类型必须是 real_time
	log.New(ctx).Info("start to verify sysUserID and templateID")
	if err := h.svc.VerifyRealTime(ctx, req.Cid, req.EmlSysUserId, req.EmlTemplateId); err != nil {
		log.New(ctx).Error("failed by checking sysUser and template", "error", err)
		parameterErr(ctx, fmt.Sprintf("Incorrect user ID: %d or template ID %d. err: %s", req.EmlSysUserId, req.EmlTemplateId, err))
		return
	}
	if req.UsrUserId == "" {
		req.UsrUserId = req.ToEmail
		log.New(ctx).Info("no usrUserId, use toEmail as usrUserId")
	}

	// 校验服务是否繁忙
	// 根据 sysUser 的最大速率进行计算timeout 时间内发送邮件数量
	// 和剩余的邮件数量进行比较
	log.New(ctx).Info("start to calculate rate", "req", req)
	ok, err := h.svc.CalculateRate(ctx, req)
	if err != nil || !ok {
		log.New(ctx).Error("The service is busy and unable to send emails within the time limit", "timeout", req.Timeout, "error", err)
		serverBusy(ctx, "The service is busy and unable to send emails within the time limit "+err.Error())
		return
	}
	log.New(ctx).Info("calculate rate successfully. start to create task", "req", req)
	// 创建父任务和子任务
	var taskSubID int64
	if taskSubID, err = h.svc.CreateTasks(ctx, req); err != nil {
		log.New(ctx).Error("Create task error", "error", err)
		serverErr(ctx, "Create task error "+err.Error())
		return
	}
	log.New(ctx).Info("create tasks rate successfully")
	// 返回结果
	success(ctx, taskSubID)
	return
}

func parameterErr(ctx *gin.Context, description string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Status:  StatusError,
		Message: "Task creation failed",
		Data: Data{
			TaskSubId:      0,
			Classification: ParameterErr,
			Description:    description,
		},
	})

}

func serverBusy(ctx *gin.Context, description string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Status:  StatusError,
		Message: "Task creation failed",
		Data: Data{
			TaskSubId:      0,
			Classification: ServerBusy,
			Description:    description,
		},
	})
}

func success(ctx *gin.Context, taskSubID int64) {
	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Status:  StatusSuccess,
		Message: "Created Successfully",
		Data: Data{
			TaskSubId:      taskSubID,
			Classification: Success,
			Description:    "success",
		},
	})
}

func serverErr(ctx *gin.Context, description string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    200,
		Status:  StatusError,
		Message: "Task creation failed",
		Data: Data{
			TaskSubId:      0,
			Classification: ServerErr,
			Description:    description,
		},
	})
}
