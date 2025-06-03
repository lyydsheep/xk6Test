package service

import (
	"context"
	"email/common/enum"
	"email/common/log"
	"email/dal/model"
	"email/logic/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand/v2"
	"strings"
	"time"
)

type EmailTaskService struct {
	EmailTaskSubRepo            repository.EmailTaskSubRepo
	TemplateRepo                repository.EmailTemplateRepo
	UsrUserRepo                 repository.UsrUserRepo
	EmailSysUserRepo            repository.EmailSysUserRepo
	EmailTaskRepo               repository.EmailTaskRepo
	EmailUnsubscribeUsrUserRepo repository.EmailUnsubscribeUsrUserRepo
	TaskConfigRepo              repository.TaskConfigRepo
}

func (svc *EmailTaskService) ProductTaskSub(ctx context.Context) {
	log.New(ctx).Info("start to product taskSub")
	for {
		// 1.定时轮询数据库
		// 并发操作，需要加锁
		randTime := rand.Float32()
		time.Sleep(5*time.Second + time.Duration(randTime)*time.Second)
		task, err := svc.EmailTaskRepo.GetEmailTask(ctx)
		if err != nil {
			log.New(ctx).Error("Fail to read task", "error", err)
			continue
		}
		if task == nil {
			log.New(ctx).Debug("no task")
			time.Sleep(5*time.Second + time.Duration(randTime)*time.Second)
			continue
		}

		// 一边生成一边存储 ---> 消费端消费
		// 批量写
		// 注意超时控制

		/* 根据邮件任务生成子任务 */

		// 1.查询标签对应的用户
		usrUsers, err := svc.UsrUserRepo.ReadTagUser(ctx, task.Cid, task.UserTag)
		if err != nil {
			log.New(ctx).Error("Fail to fetch tag user", "taskID", task.ID, "tag", task.UserTag, "error", err)
			continue
		}

		log.New(ctx).Info("setSuccess to fetch tag user", "taskID", task.ID, "tag", task.UserTag, "userNum", len(usrUsers))

		// 2.生成子任务

		priority, err := svc.TemplateRepo.ReadPriority(ctx, task.TemplateID)
		if err != nil {
			log.New(ctx).Error("Fail to read task priority", "taskID", task.ID, "error", err)
			dingTalk.notifyError("error: 读取邮件任务配置出错", err)
			continue
		}
		log.New(ctx).Info("read task priority", "taskID", task.ID, "priority", priority)

		for i := range usrUsers {
			// 查看用户是否在退订列表
			_, err = svc.EmailUnsubscribeUsrUserRepo.Read(ctx, task.Cid, usrUsers[i].UserID)
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			taskSub := &model.EmlTaskSub{
				Cid:         task.Cid,
				TaskID:      task.ID,
				FromEmail:   task.FromEmail,
				ToUserID:    usrUsers[i].UserID,
				ToEmail:     usrUsers[i].Email,
				TemplateID:  task.TemplateID,
				Status:      enum.EmailTaskSubStatusUnsent,
				Version:     0,
				Type:        task.Type,
				Retry:       0,
				Priority: priority,
				EmailDomain: strings.Split(usrUsers[i].Email, "@")[1],
			}

			taskSub, err = svc.EmailTaskSubRepo.Create(ctx, taskSub)
			if err != nil {
				// 创建子任务失败，直接进入子任务终态
				err = svc.EmailTaskSubRepo.HandleFinal(ctx, 0, enum.EmailTaskSubStatusFailure, task.ID, fmt.Sprintf("create emailSubTask Fail. ToUserID: %s. err: %s", taskSub.ToUserID, err.Error()), 0, 1, -1, 0)
				log.New(ctx).Error("Fail to create taskSub", "taskID", task.ID, "toEmail", usrUsers[i].Email, "error", err)
				continue
			}
			log.New(ctx).Debug("setSuccess to create taskSub", "taskID", task.ID, "toEmail", usrUsers[i].Email, "taskSubID", taskSub.ID)
		}

		// 3.更新邮件任务状态
		if err = svc.EmailTaskRepo.UpdateEmailTask(ctx, task.ID, enum.EmailTaskStatusFinished, 0, 0); err != nil {
			log.New(ctx).Error("Fail to update task status", "taskID", task.ID, "error", err)
		}
	}
}

func NewEmailTaskService(TaskSubRepo repository.EmailTaskSubRepo,
	TemplateRepo repository.EmailTemplateRepo, TaskConfigRepo repository.TaskConfigRepo,
	UsrUserRepo repository.UsrUserRepo, EmailTaskRepo repository.EmailTaskRepo, EmailSysUserRepo repository.EmailSysUserRepo, EmailUnsubscribeUsrUserRepo repository.EmailUnsubscribeUsrUserRepo) *EmailTaskService {
	return &EmailTaskService{
		TemplateRepo:                TemplateRepo,
		EmailTaskSubRepo:            TaskSubRepo,
		UsrUserRepo:                 UsrUserRepo,
		EmailTaskRepo:               EmailTaskRepo,
		EmailSysUserRepo:            EmailSysUserRepo,
		EmailUnsubscribeUsrUserRepo: EmailUnsubscribeUsrUserRepo,
		TaskConfigRepo:              TaskConfigRepo,
	}
}
