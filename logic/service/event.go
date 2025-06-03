package service

import (
	"context"
	"email/common/enum"
	"email/common/log"
	"email/common/request"
	"email/dal/model"
	"email/logic/repository"
	"errors"
	"gorm.io/gorm"
	"time"
)

type EventService struct {
	EmailSysUserRepo  repository.EmailSysUserRepo
	EmailTaskRepo     repository.EmailTaskRepo
	EmailTaskSubRepo  repository.EmailTaskSubRepo
	EmailTemplateRepo repository.EmailTemplateRepo
}

func NewEventService(emailSysUserRepo repository.EmailSysUserRepo, emailTaskRepo repository.EmailTaskRepo, emailTaskSubRepo repository.EmailTaskSubRepo, emailTemplateRepo repository.EmailTemplateRepo) *EventService {
	return &EventService{
		EmailSysUserRepo:  emailSysUserRepo,
		EmailTaskRepo:     emailTaskRepo,
		EmailTaskSubRepo:  emailTaskSubRepo,
		EmailTemplateRepo: emailTemplateRepo,
	}
}

func (svc *EventService) CalculateRate(ctx context.Context, req request.CreateRTEventReq) (bool, error) {
	// 校验服务是否繁忙
	// 根据 sysUser 的最大速率进行计算timeout 时间内发送邮件数量
	// 和剩余的邮件数量进行比较
	sysUser, err := svc.EmailSysUserRepo.ReadByType(ctx, req.Cid, enum.SysUserRealTime)
	if err != nil {
		log.New(ctx).Error("Failed by reading sysUser", "error", err)
		return false, err
	}
	log.New(ctx).Debug("get sysUser", "sysUser", sysUser)
	interval, err := enum.SpeedToTime(sysUser.DefaultSpeed)
	if err != nil {
		log.New(ctx).Error("invalid speed format", "speed", sysUser.DefaultSpeed, "err", err)
		return false, err
	}
	log.New(ctx).Debug("get interval", "interval", interval, "fromAddress", sysUser.FromAddress)
	// 根据 sysUser 的最大速率进行计算timeout 时间内发送邮件数量
	maxCount := int64(time.Minute * time.Duration(req.Timeout) / interval)
	existCount, err := svc.EmailTaskSubRepo.ReadInitCount(ctx, sysUser.FromAddress)
	if err != nil {
		log.New(ctx).Error("Failed by reading init count", "error", err)
		return false, err
	}
	log.New(ctx).Info("Maximum number of dispatches and existing dispatches were successfully captured", "maxCount", maxCount, "existCount", existCount)
	if existCount >= maxCount {
		log.New(ctx).Warn("The maximum number of dispatches has been reached", "maxCount", maxCount, "existCount", existCount)
		return false, nil
	}
	return true, nil
}

func (svc *EventService) CreateTasks(ctx context.Context, req request.CreateRTEventReq) (int64, error) {
	// 创建父任务和子任务
	sysUser, err := svc.EmailSysUserRepo.ReadByType(ctx, req.Cid, enum.SysUserRealTime)
	if err != nil {
		log.New(ctx).Error("Failed by reading sysUser", "error", err)
		return 0, err
	}
	log.New(ctx).Debug("Read sysUser for Create Task", "sysUser", sysUser)
	tx := svc.EmailTaskRepo.(*repository.EmailTaskRepoV1).Query.Begin()
	defer tx.Rollback()
	task := model.EmlTask{
		Cid:       req.Cid,
		FromEmail: sysUser.FromAddress,
		Status:    enum.EmailTaskStatusFinished,
	}
	if err = svc.EmailTaskRepo.CreateWithTx(ctx, &task, tx); err != nil {
		log.New(ctx).Error("Failed by creating task", "error", err)
		return 0, err
	}
	taskSub := model.EmlTaskSub{
		Cid:        req.Cid,
		TaskID:     task.ID,
		FromEmail:  sysUser.FromAddress,
		ToUserID:   req.UsrUserId,
		ToEmail:    req.ToEmail,
		TemplateID: req.EmlTemplateId,
		Status:     enum.EmailTaskSubStatusUnsent,
		Version:    0,
		Type:       enum.TaskConfigRealTime,
		Retry:      0,
		// 即时发送，优先级最高
		Priority: 0,
		Data:     req.TemplateVariables,
	}
	if err = svc.EmailTaskSubRepo.CreateWithTx(ctx, &taskSub, tx); err != nil {
		log.New(ctx).Error("Failed by creating taskSub", "error", err)
		return 0, err
	}
	tx.Commit()
	return taskSub.ID, nil
}

func (svc *EventService) VerifyRealTime(ctx context.Context, cid int64, sysUserId int64, templateId int64) error {
	sysUser, err := svc.EmailSysUserRepo.ReadByType(ctx, cid, enum.SysUserRealTime)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.New(ctx).Error("Real time mail user does not exist", "ID from request", sysUserId, "err", err)
			return errors.New("sysUser does not exist")
		}
		log.New(ctx).Error("Failed by reading sysUser", "error", err)
		return err
	}
	if sysUser.ID != sysUserId {
		log.New(ctx).Error("sysUser.ID is wrong", "ID from database", sysUser.ID, "ID from request", sysUserId)
		return errors.New("sysUser.ID is wrong")
	}
	log.New(ctx).Info("sysUserId is correct")
	template, err := svc.EmailTemplateRepo.Read(ctx, cid, templateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.New(ctx).Error("template does not exist", "err", err)
			return err
		}
		log.New(ctx).Error("Failed by reading template")
		return err
	}
	if template.Category != enum.TemplateCategoryRealTime {
		log.New(ctx).Error("template category is wrong", "template.Category", template.Category)
		return errors.New("template category is not real time")
	}

	log.New(ctx).Info("sysUserId is correct")
	return nil
}
