package service

import (
	"context"
	"email/common/enum"
	"email/common/log"
	"email/common/util"
	"email/config"
	"email/dal/dao"
	"email/dal/model"
	"email/logic/repository"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type EmailTaskSubService struct {
	EmailTaskSubRepo            repository.EmailTaskSubRepo
	EmailSysUserRepo            repository.EmailSysUserRepo
	EmailTaskRepo               repository.EmailTaskRepo
	UsrUserRepo                 repository.UsrUserRepo
	TemplateRepo                repository.EmailTemplateRepo
	EmlRedemptionCodeRepo       repository.EmailRedemptionCodeRepo
	EmailRedemptionTaskSubRepo  repository.EmailRedemptionTaskSubRepo
	EmailUnsubscribeUsrUserRepo repository.EmailUnsubscribeUsrUserRepo
	EmailDomainCredibilityRepo  repository.EmailDomainCredibilityRepo
	Query                       *dao.Query
}

func NewEmailTaskSubService(emailTaskSubRepo repository.EmailTaskSubRepo, emailSysUserRepo repository.EmailSysUserRepo,
	emailTaskRepo repository.EmailTaskRepo, usrUserRepo repository.UsrUserRepo, templateRepo repository.EmailTemplateRepo,
	EmlRedemptionCodeRepo repository.EmailRedemptionCodeRepo, emailRedemptionTaskSubRepo repository.EmailRedemptionTaskSubRepo,
	emailUnsubscribeUsrUserRepo repository.EmailUnsubscribeUsrUserRepo, emailDomainCredibilityRepo repository.EmailDomainCredibilityRepo,
	query *dao.Query) *EmailTaskSubService {
	return &EmailTaskSubService{
		EmailTaskSubRepo:            emailTaskSubRepo,
		EmailSysUserRepo:            emailSysUserRepo,
		EmailTaskRepo:               emailTaskRepo,
		UsrUserRepo:                 usrUserRepo,
		TemplateRepo:                templateRepo,
		EmlRedemptionCodeRepo:       EmlRedemptionCodeRepo,
		EmailRedemptionTaskSubRepo:  emailRedemptionTaskSubRepo,
		EmailUnsubscribeUsrUserRepo: emailUnsubscribeUsrUserRepo,
		EmailDomainCredibilityRepo:  emailDomainCredibilityRepo,
		Query:                       query,
	}
}

type EmailData struct {
	FromEmail     string
	ToEmail       string
	Subject       string
	ContentType   string
	Content       string
	AliyunTagName string
}

func (svc *EmailTaskSubService) Run(ctx context.Context) {
	m := map[string]struct{}{}
	for {
		sysUsers, err := svc.EmailSysUserRepo.ReadAll(ctx)
		if err != nil {
			log.New(ctx).Error("Fail to read all sysUsers", "error", err)
			return
		}
		// TODO 默认所有都是阿里云 api，后续得支持更灵活的配置信息
		for _, sysUser := range sysUsers {
			if sysUser.EmailServer == enum.Ali {
				key := fmt.Sprintf("%dkey%s", sysUser.Cid, sysUser.Username)
				if _, ok := m[key]; !ok {
					// 不存在
					log.New(ctx).Info("start to consume taskSub", "accountName", sysUser.Username)
					m[key] = struct{}{}
					// TODO 传入新的 ctx
					go svc.Consume(ctx, sysUser)
				}
			}
		}
		time.Sleep(time.Hour)
	}
}

func (svc *EmailTaskSubService) Consume(ctx context.Context, sysUser *model.EmlSysUser) {
	var errCount int64
	for {
		// 1.定时轮询数据库
		randTime := rand.Float32()
		lastErrCount := atomic.LoadInt64(&errCount)
		emailTaskSub, err := svc.EmailTaskSubRepo.GetBatchByFromEmail(ctx, sysUser.FromAddress)
		if err != nil {
			log.New(ctx).Error("Fail to fetch taskSub", "error", err)
			continue
		}

		if len(emailTaskSub) == 0 {
			log.New(ctx).Debug("no taskSub to consume", "sysUser.ID", sysUser.ID, "sysUser.AccountName", sysUser.Username)
			time.Sleep(1*time.Minute + time.Duration(randTime)*time.Second)
		}

		// 2.消费任务
		for i := range emailTaskSub {
			if len(emailTaskSub[i].EmailDomain) < 2 {
				emailTaskSub[i].EmailDomain = strings.Split(emailTaskSub[i].ToEmail, "@")[1]
			}
			entity, err := svc.EmailDomainCredibilityRepo.Read(ctx, emailTaskSub[i].Cid, emailTaskSub[i].EmailDomain,
				emailTaskSub[i].FromEmail)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					log.New(ctx).Error("Fail to read emailDomainCredibility", "taskSubID", emailTaskSub[i].ID, "error", err)
				}
				// TODO 集群时得加锁
				// 不存在则新创建一个
				entity, err = svc.EmailDomainCredibilityRepo.Create(ctx, &model.EmlDomainCredibility{
					Cid:             emailTaskSub[i].Cid,
					Domain:          emailTaskSub[i].EmailDomain,
					FromAddress:     emailTaskSub[i].FromEmail,
					SuccessRateHour: util.Float64Pointer(0),
					SuccessRateDay:  util.Float64Pointer(0),
					Speed:           sysUser.DefaultSpeed,
					LastSentTime:    time.Now().Add(-time.Hour),
				})
				if err != nil {
					log.New(ctx).Error("Fail to create emailDomainCredibility", "taskSubID", emailTaskSub[i].ID, "error", err)
					continue
				}
			}
			lastTime := entity.LastSentTime
			speed := entity.Speed
			Sleep(ctx, speed, lastTime)
			svc.work(emailTaskSub[i], &errCount, entity)
		}
		if atomic.LoadInt64(&errCount)-lastErrCount > 10 {
			dingTalk.notifyFatal("fatal: 短时间内错误过多。睡眠十分钟", errors.New(""))
			log.New(ctx).Error("too many errors. sleep 10 minutes")
			time.Sleep(10 * time.Minute)
		}
	}
}

func Sleep(ctx context.Context, speed string, lastTime time.Time) {
	//if _, ok := enum.MapStringToSpeed[speed]; !ok {
	//	log.New(ctx).Error("invalid speed format, sleep 16 second", "speed", speed)
	//	time.Sleep(16 * time.Second)
	//	return
	//}
	// 已经睡了多久
	gap := time.Since(lastTime)
	if gap <= 0 {
		log.New(ctx).Debug("no need to sleep")
		return
	}
	interval, err := enum.SpeedToTime(speed)
	if err != nil {
		log.New(ctx).Error("invalid speed format, sleep 16 second", "speed", speed, "err", err)
		interval = time.Second * 16
	}
	// 还需要睡多久, 如果≤0， 默认直接返回
	//sleepTime := enum.MapStringToSpeed[speed] - gap
	sleepTime := interval - gap
	log.New(ctx).Debug(fmt.Sprintf("sleep %v", sleepTime))
	// 后续的sql 查询可能会导致两个 goroutine 之间的时间间隔变短
	time.Sleep(sleepTime)
}

func (svc *EmailTaskSubService) work(emailTaskSub *model.EmlTaskSub, errCount *int64, entity *model.EmlDomainCredibility) {
	newCtx := context.Background()
	valCtx := context.WithValue(newCtx, "traceId", fmt.Sprintf("%d", emailTaskSub.ID))
	log.New(valCtx).Info("start to consume taskSub", "taskSubID", emailTaskSub.ID,
		"startTime", time.Now().Format("2006-01-02 15:04:05"), "emailDomain", emailTaskSub.EmailDomain)
	err := svc.consumeTaskSub(valCtx, emailTaskSub, entity)
	if err != nil {
		log.New(valCtx).Error("Fail to consume taskSub. ErrCount + 1", "taskSubID", emailTaskSub.ID, "error", err)
		atomic.AddInt64(errCount, 1)
		dingTalk.notifyError(fmt.Sprintf("error: 发件出错, taskSubID: %d, toEmail: %s", emailTaskSub.ID, emailTaskSub.ToEmail), err)
	}
}

func (svc *EmailTaskSubService) consumeTaskSub(ctx context.Context, taskSub *model.EmlTaskSub, entity *model.EmlDomainCredibility) error {
	// 查询 smtp 配置
	sysUser, err := svc.EmailSysUserRepo.Read(ctx, taskSub.Cid, taskSub.FromEmail)
	if err != nil {
		log.New(ctx).Error("Fail to read sysUser", "CID: ", taskSub.Cid, "fromEmail", taskSub.FromEmail, "error", err)
		if err := svc.EmailTaskSubRepo.HandleFinal(ctx, taskSub.ID, enum.EmailTaskSubStatusFailure, taskSub.TaskID, fmt.Sprintf("Read EmlSysUser: %s", err.Error()), 0, 0, -1, 0); err != nil {
			log.New(ctx).Error("Fail to update taskSub status failure", "taskSubID", taskSub.ID, "error", err)
		}
		return errors.New("Fail to read emailSysUser " + err.Error())
	}

	// 聚合数据
	data, err := svc.aggregateData(ctx, sysUser, taskSub)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.New(ctx).Warn("Fail to consume taskSub, sleep some seconds then try again", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID, "err", err)
		time.Sleep(time.Second + time.Second*time.Duration(rand.IntN(10)) + time.Second*time.Duration(rand.Float32()))
		data, err = svc.aggregateData(ctx, sysUser, taskSub)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	if err != nil {
		log.New(ctx).Error("Fail to aggregate data", "taskSubID", taskSub.ID, "error", err)
		if err := svc.EmailTaskSubRepo.HandleFinal(ctx, taskSub.ID, enum.EmailTaskSubStatusFailure, taskSub.TaskID, fmt.Sprintf("Read EmlSysUser: %s", err.Error()), 0, 0, -1, 0); err != nil {
			log.New(ctx).Error("Fail to update taskSub status failure", "taskSubID", taskSub.ID, "error", err)
		}
		return errors.New("Fail to aggregateData " + err.Error())
	}

	// 更新发件时间
	entity.LastSentTime = time.Now()
	if _, err = svc.EmailDomainCredibilityRepo.UpdateLastSentTime(ctx, entity); err != nil {
		log.New(ctx).Error("Fail to update emailDomainCredibility", "taskSubID", taskSub.ID, "error", err)
		return errors.New("Fail to update emailDomainCredibility " + err.Error())
	}
	// 执行任务
	if err = svc.executeTask(ctx, sysUser, taskSub, data); err != nil {
		log.New(ctx).Error("Fail to send email", "taskSubID", taskSub.ID, "error", err)
		return errors.New("Fail to send email " + err.Error())
	}
	// TODO 发件成功数据量非常大, 需要按照周期聚合后发布
	dingTalk.notifyDebug(fmt.Sprintf("info: 发件成功. CID: %d, taskSubID %d, toEmail: %s", taskSub.Cid, taskSub.ID, taskSub.ToEmail))
	log.New(ctx).Info("email sent successfully", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID)
	return nil
}

func (svc *EmailTaskSubService) executeTask(ctx context.Context, sysUser *model.EmlSysUser, taskSub *model.EmlTaskSub, data *EmailData) error {
	// 预先设置 终态，避免多次发送
	// TODO 运行过程中如果服务中断，会导致任务无法重试
	taskSub.Status = enum.EmailTaskSubStatusPreDone
	row, err := svc.EmailTaskSubRepo.Update(ctx, taskSub.Version, taskSub)
	if err != nil {
		log.New(ctx).Error("Fail to preconfigured taskSub status sent", "taskSubID", taskSub.ID, "error", err)
		return err
	}
	if row == 0 {
		// 说明已经被其他 goroutine 占用了
		log.New(ctx).Warn("taskSub has been taken", "taskSubID", taskSub.ID, "taskID", taskSub.TaskID)
		return nil
	}

	// 发送邮件
	log.New(ctx).Info("start to send email", "taskID", taskSub.TaskID, "to", data.ToEmail)
	err = util.SendEmail(sysUser.Host, int(sysUser.Port), data.AliyunTagName, sysUser.Username, sysUser.Password, taskSub.FromEmail, data.ToEmail, data.Subject, data.ContentType, data.Content)

	if err != nil {
		log.New(ctx).Warn("first send email failed. setRetry to send email", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID, "error", err)
		if err = svc.retry(ctx, sysUser, taskSub, data); err != nil {
			//  开事务更新 subTask 和 task 状态
			svc.EmailTaskSubRepo.HandleFinal(ctx, taskSub.ID, enum.EmailTaskSubStatusFailure, taskSub.TaskID, err.Error(), 1, 0, -1, 0)
			log.New(ctx).Error("Fail to send email", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID, "error", err)
			return err
		}
		// 重试成功
		log.New(ctx).Info("setRetry setSuccess", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID)
	}
	taskSub.Status = enum.EmailTaskSubStatusSent
	//now := time.Now()
	//taskSub.SentTime = &now
	row, err = svc.EmailTaskSubRepo.Update(ctx, taskSub.Version, taskSub)
	if err != nil {
		log.New(ctx).Error("Fail to preconfigured taskSub status sent", "taskSubID", taskSub.ID, "error", err)
		return err
	}
	if row == 0 {
		// 说明已经被其他 goroutine 占用了
		log.New(ctx).Warn("taskSub has been taken", "taskSubID", taskSub.ID, "taskID", taskSub.TaskID)
		return nil
	}
	// 发送成功
	// 更新 emailTask 状态
	// 成功不使用事务
	//return svc.EmailTaskRepo.UpdateEmailTask(ctx, taskSub.TaskID, "", 1, -1)
	return nil
}

func (svc *EmailTaskSubService) retry(ctx context.Context, sysUser *model.EmlSysUser, taskSub *model.EmlTaskSub, data *EmailData) error {
	time.Sleep(time.Second)
	const (
		// add by 米畈2025-05-12 09:59:50 - 在此不进行重试, 重试机制通过eml_task_sub表统一实现，在此重试不方便进行流控
		maxRetries     = 1
		initialBackoff = 2 * time.Second
		maxBackoff     = 16 * time.Second
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// 指数退避 2 4 8 16
		time.Sleep(initialBackoff * time.Duration(1<<uint(attempt-1)))

		err := util.SendEmail(sysUser.Host, int(sysUser.Port), data.AliyunTagName, sysUser.Username, sysUser.Password, data.FromEmail, data.ToEmail, data.Subject, data.ContentType, data.Content)
		if err == nil {
			return nil
		}

		log.New(ctx).Warn("setRetry attempt failed", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID, "attempt", attempt, "error", err)

		if attempt >= maxRetries {
			return err
		}
	}

	return errors.New("max retries exceeded")
}

func (svc *EmailTaskSubService) aggregateData(ctx context.Context, sysUser *model.EmlSysUser, taskSub *model.EmlTaskSub) (*EmailData, error) {
	// 获取目标用户信息
	usrUser, err := svc.UsrUserRepo.Read(ctx, taskSub.Cid, taskSub.ToUserID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.New(ctx).Error("Fail to read usrUser", "cid", taskSub.Cid, "userID", taskSub.ToUserID, "error", err)
			return nil, err
		}
		log.New(ctx).Warn("usrUser not found", "cid", taskSub.Cid, "userID", taskSub.ToUserID)
		usrUser = &model.UsrUser{
			Language: enum.EmailTemplateEN,
			Email:    taskSub.ToEmail,
		}
	}

	// 获取模板信息
	template, err := svc.TemplateRepo.Read(ctx, taskSub.Cid, taskSub.TemplateID)
	if err != nil {
		log.New(ctx).Error("Fail to read template", "cid", taskSub.Cid, "templateID", taskSub.TemplateID, "error", err)
		return nil, err
	}

	//  根据模板类型组装数据
	var (
		data        map[string]interface{}
		specialData map[string]interface{}
	)
	// 组装通用数据和预设置数据
	data, err = svc.commonData(ctx, usrUser, sysUser, taskSub)
	if err != nil {
		log.New(ctx).Error("Fail to assemble common data", "error", err)
		return nil, err
	}
	switch template.Category {
	case enum.TemplateCategoryRedemption:
		// 兑换码邮件
		specialData, err = svc.redemptionCodeData(ctx, taskSub)
	case enum.TemplateCategoryCommon:
	case enum.TemplateCategoryRealTime:
	default:
		return nil, errors.New("unknown template category " + template.Category)
	}
	if err != nil {
		log.New(ctx).Error("Fail to aggregate data", "taskSub.ID", taskSub.ID, "templateID", template.ID, "error", err)
		return nil, err
	}
	// 补充特殊模板的数据
	// 预设置优先级 ＞ 模板数据
	for k, v := range specialData {
		if data[k] == nil {
			data[k] = v
		}
	}
	log.New(ctx).Info("setSuccess to aggregate data by template", "taskSub.ID", taskSub.ID, "data", data)

	// 根据默认语言选择模板
	subject, content := template.SubjectEn, template.ContentEn
	if usrUser.Language == enum.EmailTemplateZH {
		subject, content = template.SubjectZh, *template.ContentZh
	}
	// TXT 和 HTML 渲染步骤一致
	realContent, err := util.RenderEmailContent(context.Background(), content, data)
	if err != nil {
		log.New(context.Background()).Error("Fail to render toEmail content", "error", err)
		return nil, err
	}
	return &EmailData{
		FromEmail:     taskSub.FromEmail,
		ToEmail:       usrUser.Email,
		Subject:       subject,
		Content:       realContent,
		ContentType:   template.ContentType,
		AliyunTagName: template.TemplateName,
	}, nil
}

func (svc *EmailTaskSubService) commonData(ctx context.Context, usrUser *model.UsrUser, sysUser *model.EmlSysUser, taskSub *model.EmlTaskSub) (map[string]interface{}, error) {
	// 组装通用的数据
	decryptedVal := fmt.Sprintf("%s%d", usrUser.UserID, usrUser.Cid)
	encryptVal, err := util.Encrypt(decryptedVal, config.App.AESKEY)
	base64Str := util.EncodeBase64(encryptVal)
	if err != nil {
		log.New(context.Background()).Error("Fail to commonData", "error", err)
		return nil, err
	}
	data := map[string]interface{}{
		// 退订链接
		"Cancel": fmt.Sprintf("%s&key=%s", config.App.EventBridgeUrl, base64Str),
	}

	// 获取 taskSub 中预设置的 data
	if taskSub.Data != "" {
		err = json.Unmarshal([]byte(taskSub.Data), &data)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to unmarshal taskSub data. %s", err))
		}
	}
	log.New(ctx).Info("assemble commonData successfully", "data", data)

	taskSub.UpdateTime = time.Now()
	bytesData, _ := json.Marshal(data)
	taskSub.Data = string(bytesData)
	if _, err = svc.Query.EmlTaskSub.WithContext(ctx).Where(svc.Query.EmlTaskSub.ID.Eq(taskSub.ID)).Updates(taskSub); err != nil {
		log.New(ctx).Error("Fail to commonData", "error", err)
		return nil, err
	}

	return data, nil
}

func (svc *EmailTaskSubService) redemptionCodeData(ctx context.Context, taskSub *model.EmlTaskSub) (map[string]interface{}, error) {
	redemptionTaskSub, err := svc.EmailRedemptionTaskSubRepo.Read(ctx, taskSub.Cid, taskSub.ID)
	if err != nil {
		log.New(ctx).Error("Fail to read EmlRedemptionTaskSub", "taskSubID", taskSub.ID, "tepe", "error", err)
		return nil, err
	}

	// 如果已经存在，则直接返回
	var redemptionCode *model.EmlRedemptionCode
	redemptionCode, err = svc.Query.EmlRedemptionCode.WithContext(ctx).Where(svc.Query.EmlRedemptionCode.Cid.Eq(taskSub.Cid),
		svc.Query.EmlRedemptionCode.Type.Eq(redemptionTaskSub.Type), svc.Query.EmlRedemptionCode.Amount.Eq(redemptionTaskSub.Amount), svc.Query.EmlRedemptionCode.TaskSubID.Eq(taskSub.ID)).First()
	if err == nil {
		log.New(ctx).Debug("code already exists", "code", redemptionCode)
		return map[string]interface{}{
			"code": redemptionCode.Code,
		}, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.New(ctx).Error("Fail to read redemptionCode", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID, "error", err)
		return nil, err
	}

	// TODO 开启事务，更新兑换码和邮件任务
	tx := svc.Query.Begin()
	defer tx.Rollback()
	// 没有被发送的兑换码 ---> task_sub_id = 0
	redemptionCode, err = tx.EmlRedemptionCode.WithContext(ctx).Where(tx.EmlRedemptionCode.Cid.Eq(taskSub.Cid),
		tx.EmlRedemptionCode.Type.Eq(redemptionTaskSub.Type), tx.EmlRedemptionCode.Amount.Eq(redemptionTaskSub.Amount), tx.EmlRedemptionCode.TaskSubID.Eq(0)).First()
	if err != nil {
		log.New(ctx).Error("Fail to read redemptionCode", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID, "error", err)
		return nil, err
	}

	// 更新 redemptionCode 表
	redemptionCode.TaskSubID = taskSub.ID
	redemptionCode.UpdateTime = time.Now()
	info, err := tx.EmlRedemptionCode.WithContext(ctx).Where(tx.EmlRedemptionCode.ID.Eq(redemptionCode.ID),
		tx.EmlRedemptionCode.TaskSubID.Eq(0)).Updates(redemptionCode)
	if err != nil {
		return nil, err
	}
	if info.RowsAffected == 0 {
		log.New(ctx).Warn("redemptionCode not found", "taskID", taskSub.TaskID, "taskSubID", taskSub.ID)
		return nil, gorm.ErrRecordNotFound
	}

	// 更新 redemptionTaskSub 表
	redemptionTaskSub.UpdateTime = time.Now()
	redemptionTaskSub.Code = redemptionCode.Code
	if _, err = tx.EmlRedemptionTaskSub.WithContext(ctx).Where(tx.EmlRedemptionTaskSub.ID.Eq(redemptionTaskSub.ID)).Updates(redemptionTaskSub); err != nil {
		log.New(ctx).Error("Fail to update redemptionTaskSub", "taskSubID", taskSub.ID, "error", err)
		return nil, err
	}

	// 更新 taskSub
	data := map[string]interface{}{
		"Code": redemptionCode.Code,
	}
	bytesData, _ := json.Marshal(data)
	taskSub.Data = string(bytesData)
	taskSub.UpdateTime = time.Now()
	if info, err = tx.EmlTaskSub.WithContext(ctx).Where(tx.EmlTaskSub.ID.Eq(taskSub.ID)).Updates(taskSub); err != nil {
		log.New(ctx).Error("Fail to update taskSub", "taskSubID", taskSub.ID, "error", err)
		return nil, err
	}
	tx.Commit()

	return data, nil
}

func (svc *EmailTaskSubService) Unsubscribe(ctx context.Context, entity *model.EmlUnsubscribeUsrUser, key string) error {
	decryptedVal, err := util.Decrypt(key, config.App.AESKEY)
	if err != nil {
		log.New(ctx).Error("Fail to decrypt", "error", err)
		return err
	}
	uid, cidStr := decryptedVal[:len(decryptedVal)-1], decryptedVal[len(decryptedVal)-1:]
	cid, err := strconv.Atoi(cidStr)
	if err != nil {
		log.New(ctx).Error("Fail to convert cid", "error", err)
		return err
	}
	entity.Cid, entity.UID = int64(cid), &uid
	if _, err = svc.EmailUnsubscribeUsrUserRepo.Create(ctx, entity); err != nil {
		log.New(ctx).Error("Fail to create emailUnsubscribeUsrUser", "error", err)
		return err
	}
	log.New(ctx).Info("setSuccess to unsubscribe", "cid", cid, "uid", uid)
	return nil
}
