package service

import (
	"context"
	"email/common/enum"
	"email/common/log"
	"email/common/util"
	"email/config"
	"email/dal/model"
	"email/logic/repository"
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm20151123 "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
	"math"
	"strconv"
	"strings"
	"time"
)

type EmailResultService struct {
	EmailTaskSubRepo           repository.EmailTaskSubRepo
	EmailSysUserRepo           repository.EmailSysUserRepo
	EmailResultRepo            repository.EmailResultRepo
	TaskConfigRepo             repository.TaskConfigRepo
	OpsEmlNumberSuccessRepo    repository.OpsEmlNumberSuccessRepo
	EmailDomainCredibilityRepo repository.EmailDomainCredibilityRepo
	EmailTemplateRepo repository.EmailTemplateRepo
}

func NewEmailResultService(emailTaskSubRepo repository.EmailTaskSubRepo, emailSysUserRepo repository.EmailSysUserRepo,
	emailResultRepo repository.EmailResultRepo, taskConfigRepo repository.TaskConfigRepo, opsEmlNumberSuccessRepo repository.OpsEmlNumberSuccessRepo, emailDomainCredibilityRepo repository.EmailDomainCredibilityRepo,
	emailTemplateRepo repository.EmailTemplateRepo) *EmailResultService {
	return &EmailResultService{
		EmailTaskSubRepo:           emailTaskSubRepo,
		EmailSysUserRepo:           emailSysUserRepo,
		EmailResultRepo:            emailResultRepo,
		TaskConfigRepo:             taskConfigRepo,
		OpsEmlNumberSuccessRepo:    opsEmlNumberSuccessRepo,
		EmailDomainCredibilityRepo: emailDomainCredibilityRepo,
		EmailTemplateRepo: emailTemplateRepo,
	}
}

func (svc *EmailResultService) Run(ctx context.Context) {
	sysUsers, err := svc.EmailSysUserRepo.ReadAll(ctx)
	if err != nil {
		log.New(ctx).Error("Fail to read all sysUsers", "error", err)
		return
	}
	// TODO 默认所有都是阿里云 api，后续得支持更灵活的配置信息
	for _, sysUser := range sysUsers {
		if sysUser.EmailServer == enum.Ali {
			sysUser.Password, err = util.Decrypt(sysUser.Password, config.App.AESKEY)
			if err != nil {
				log.New(ctx).Error("sysUser decrypt password error", "error", err)
			}
			sysUser.AccessKeyID, err = util.Decrypt(sysUser.AccessKeyID, config.App.AESKEY)
			if err != nil {
				log.New(ctx).Error("sysUser decrypt accessKeyID error", "error", err)
			}
			sysUser.AccessKeySecret, err = util.Decrypt(sysUser.AccessKeySecret, config.App.AESKEY)
			if err != nil {
				log.New(ctx).Error("sysUser decrypt accessKeySecret error", "error", err)
			}
			go svc.fetchEmailResult(ctx, sysUser)
			go svc.changeTaskSubStatus(ctx, sysUser)
		}
	}
}

// 定时扫描 taskSub 库，获取 taskSub 的结果，并根据结果设置 taskSub 的状态（重试规则)
func (svc *EmailResultService) changeTaskSubStatus(ctx context.Context, sysUser *model.EmlSysUser) {
	log.New(ctx).Info("start to change TaskSubStatus", "sysUserID", sysUser.ID, "accountName", sysUser.Username)
	for {
		log.New(ctx).Info("start to changeTaskSubStatus", "accountName", sysUser.Username, "now", time.Now())

		// 查询 result 表最新时间 endTime
		endTime, err := svc.EmailResultRepo.GetEndTime(ctx, sysUser.Username)
		log.New(ctx).Info(fmt.Sprintf("get endTime %s", endTime), "accountName", sysUser.Username)
		if err != nil {
			log.New(ctx).Error("failed to get end time from email result repository", "accountName", sysUser.Username, "error", err)
			continue
		}

		// 规避任务 fetchTime 和任务结果落库的延迟
		realEnd := endTime.Add(-time.Minute * 10)
		if time.Since(endTime) > time.Minute*10 {
			realEnd = endTime
		}
		// 在 sub 表中查询小于等于 endTime 的任务
		log.New(ctx).Info("start to query taskSubs", "accountName", sysUser.Username, "endTime", endTime)
		taskSubs, err := svc.EmailTaskSubRepo.ReadDoneAndPredoneAndOutdated(ctx, sysUser.FromAddress, realEnd)
		log.New(ctx).Info(fmt.Sprintf("A total of %d taskSubs were acquired", len(taskSubs)), "fromAddress", sysUser.FromAddress, "endTime", endTime)
		if err != nil {
			log.New(ctx).Error("Fail to read done taskSub", "error", err)
			continue
		}
		if len(taskSubs) == 0 {
			log.New(ctx).Info("no taskSubs were found to change status. sleep one minute", "accountName", sysUser.Username)
			time.Sleep(time.Minute)
			continue
		}
		for _, taskSub := range taskSubs {
			log.New(ctx).Info("start to change taskSub status", "taskSubID", taskSub.ID, "status", "toEmail", taskSub.ToEmail, taskSub.Status, "fetchTime", taskSub.FetchTime)
			var targetRes *model.OpsEmlResult
			version := taskSub.Version
			// 在 result 表中查询任务结果
			results, err := svc.EmailResultRepo.ReadByTo(ctx, taskSub.ToEmail)
			log.New(ctx).Info(fmt.Sprintf("A total of %d email results were acquired", len(results)), "accountName", sysUser.Username)
			if err != nil {
				// TODO 进入这个分支可能会导致阻塞
				log.New(ctx).Error("failed to read email results by recipient", "toEmail", taskSub.ToEmail, "error", err)
				continue
			}

			abs := func(x int64) int64 {
				if x < 0 {
					return -x
				}
				return x
			}
			minGap, fetchTime := int64(math.MaxInt64), taskSub.FetchTime.UnixMilli()
			for _, result := range results {
				// 阿里云 api 记录的 sent_time 必定在 fetch_time 后面
				if sysUser.Username == result.AccountName && (taskSub.FetchTime.Before(result.SentTime) || taskSub.FetchTime.Equal(result.SentTime)) {
					gapTime := abs(result.SentTime.UnixMilli() - fetchTime)
					if gapTime < minGap {
						minGap = gapTime
						targetRes = result
					}
				}
			}

			//  结果不存在 --> 异常情况  --->  dingTalk && 置为 notFound
			if targetRes == nil {
				taskSub.Version++
				taskSub.UpdateTime = time.Now()
				taskSub.Status = enum.EmailTaskSubStatusNotFound
				log.New(ctx).Error(fmt.Sprintf("no result. accountName is %s, toEmail is %s", sysUser.Username, taskSub.ToEmail), "sysUserID", sysUser.ID, "accountName", sysUser.Username, "toEmail", taskSub.ToEmail, "userID", taskSub.ToUserID, "taskSubID", taskSub.ID)
				dingTalk.notifyError("error：邮件任务状态转换失败", errors.New(fmt.Sprintf("fail to find email result. accountName is %s, toEmail is %s. fetchTime is %s status is %s, number of retries is %d", sysUser.Username, taskSub.ToEmail, taskSub.FetchTime, taskSub.Status, taskSub.Retry)))
				if _, err = svc.EmailTaskSubRepo.Update(ctx, version, taskSub); err != nil {
					log.New(ctx).Error("Fail to update taskSub", "sysUserID", sysUser.ID, "accountName", sysUser.Username, "error", err)
				}
				continue
			}

			// 校验是否可以重试
			code, message := strings.Split(targetRes.Message, " ")[0], targetRes.Message
			smtpCode := enum.GetSmtpCode(message)
			log.New(ctx).Info("a result was found", "accountName", sysUser.Username, "toEmail", taskSub.ToEmail, "taskSubID", taskSub.ID, "sentTime", targetRes.SentTime, "code", code, "message", targetRes.Message)
			switch smtpCode.Category {
			case enum.Success:
				// 发送成功
				log.New(ctx).Info("the return code indicates successful transmission", "accountName", sysUser.Username, "toEmail", taskSub.ToEmail, "taskSubID", taskSub.ID, "code", code, "message", message)
				setSuccess(taskSub, targetRes.SentTime, code, message)
			case enum.TemporaryFailure:
				// 可以重试
				maxRetry, err := svc.EmailTemplateRepo.ReadMaxRetry(ctx, taskSub.TemplateID)
				if err != nil {
					log.New(ctx).Error("[createCodeTask]Fail to read max retry", "templateID", taskSub.TemplateID, "error", err)
					continue
				}
				log.New(ctx).Info(fmt.Sprintf("the return code: %s indicates that the email needs to be retried", code), "accountName", sysUser.Username, "toEmail", taskSub.ToEmail, "taskSubID", taskSub.ID, "code", code, "message", message, "maxRetry", maxRetry)
				setRetry(taskSub, maxRetry, code, message, targetRes.SentTime)
			case enum.PermanentFailure:
				// 邮件发送失败
				// 进入 unknown 分支
				fallthrough
			case enum.Unknown:
				// 其他错误
				// 将邮件任务状态设置为失败
				log.New(ctx).Info(fmt.Sprintf("smtp category is %d the return code: %s indicates that the email needs to be failed", smtpCode.Category, code), "accountName", sysUser.Username, "toEmail", taskSub.ToEmail, "taskSubID", taskSub.ID, "code", code, "message", message)

				Fail(taskSub, code, message, targetRes.SentTime)
			default:
				dingTalk.notifyError("error: Entering an illegal branch", errors.New(fmt.Sprintf("smtp code category is %d", smtpCode.Category)))
				continue
			}

			if _, err = svc.EmailTaskSubRepo.Update(ctx, version, taskSub); err != nil {
				log.New(ctx).Error("Fail to update taskSub", "sysUserID", sysUser.ID, "accountName", sysUser.Username, "error", err)
				continue
			}
		}
	}
}

func setRetry(taskSub *model.EmlTaskSub, maxRetry int32, code, description string, sentTime time.Time) {
	taskSub.Version++
	taskSub.UpdateTime = time.Now()
	if taskSub.Retry < maxRetry {
		taskSub.Status = enum.EmailTaskSubStatusUnsent
		taskSub.Retry++
	} else {
		taskSub.Status = enum.EmailTaskSubStatusFailure
		taskSub.CodeSMTP, taskSub.CodeDescription = code, description
		taskSub.SentTime = &sentTime
	}
}

func Fail(taskSub *model.EmlTaskSub, code, description string, sentTime time.Time) {
	taskSub.Version++
	taskSub.Status = enum.EmailTaskSubStatusFailure
	taskSub.CodeSMTP, taskSub.CodeDescription = code, description
	taskSub.SentTime = &sentTime
	taskSub.UpdateTime = time.Now()
}

func setSuccess(taskSub *model.EmlTaskSub, sentTime time.Time, code, description string) {
	taskSub.SentTime = &sentTime
	taskSub.CodeSMTP = code
	taskSub.CodeDescription = description
	taskSub.Version++
	taskSub.UpdateTime = time.Now()
}

func (svc *EmailResultService) fetchEmailResult(ctx context.Context, sysUser *model.EmlSysUser) {
	// 实时获取新任务的结果
	client, err := createClient(sysUser)
	if err != nil {
		panic(err)
	}
	var (
		nextStart = ""
		startTime = sysUser.FetchStartTime
	)
	log.New(ctx).Info("start to fetchEmailResult", "accountName", sysUser.Username)
	for {
		time.Sleep(time.Second * 5)
		endTimeCeiling := time.Now().Add(-3 * time.Minute)
		var endTime time.Time
		if startTime.Add(time.Minute * 5).Before(endTimeCeiling) {
			endTime = startTime.Add(time.Minute * 5)
		} else {
			log.New(ctx).Info("endTimeCeiling too close to current time, sleep one minute", "accountName", sysUser.Username, "startTime", startTime, "endTime", endTimeCeiling, "now", time.Now())
			time.Sleep(time.Minute)
			continue
		}
		log.New(ctx).Info("start to fetchEmailResult", "accountName", sysUser.Username, "startTime", startTime.Format(enum.TimeFormatHyphenedYMDHI), "endTime", endTime.Format(enum.TimeFormatHyphenedYMDHI))
		resp, err := util.FetchResp(sysUser.Username, "", nextStart, startTime.Format(enum.TimeFormatHyphenedYMDHI), endTime.Format(enum.TimeFormatHyphenedYMDHI), client)
		if err != nil {
			log.New(ctx).Error("Fail to fetchEmailResult", "sysUserID", sysUser.ID, "accountName", sysUser.Username, "error", err)
			continue
		}
		// 先获取全部的数据，再插入到数据库
		details := make([]*dm20151123.SenderStatisticsDetailByParamResponseBodyDataMailDetail, 0)
		for nextStart, err = svc.processResp(ctx, resp, &details); nextStart != ""; {
			resp, err = util.FetchResp(sysUser.Username, "", nextStart, startTime.Format(enum.TimeFormatHyphenedYMDHI), endTime.Format(enum.TimeFormatHyphenedYMDHI), client)
			if err != nil {
				log.New(ctx).Error("Fail to fetch old response details", "sysUserID", sysUser.ID, "accountName", sysUser.Username, "error", err)
				break
			}
			nextStart, err = svc.processResp(ctx, resp, &details)
			if err != nil {
				log.New(ctx).Error("Fail to process resp", "sysUserID", sysUser.ID, "accountName", sysUser.Username, "error", err)
				break
			}
		}

		records := make([]*model.OpsEmlResult, 0, len(details))
		// 按照时间由远及近插入数据库
		for i := len(details) - 1; i >= 0; i-- {
			detail := details[i]
			// 将数据插入到数据库中
			sendTime, err := strconv.ParseInt(*detail.LastUpdateTime, 10, 64)
			if err != nil {
				log.New(ctx).Error("Fail to parse send time", "accountName", detail.AccountName, "error", err)
				continue
			}
			record := detailToResult(sendTime, detail)
			if err = svc.EmailResultRepo.Create(ctx, record); err != nil {
				log.New(ctx).Error("Fail to store ops eml result", "error", err)
				continue
			}
			records = append(records, record)
			log.New(ctx).Debug("store ops eml result", "accountName", detail.AccountName, "errorClassification", detail.ErrorClassification, "sentTime", detail.LastUpdateTime, "message", detail.Message, "status", detail.Status, "subject", detail.Subject, "toAddress", detail.ToAddress)
		}

		svc.calculateRate(ctx, sysUser, startTime, endTime, records)
		// 推迟 5 分钟
		startTime = endTime
		if err = svc.EmailSysUserRepo.UpdateStartTime(ctx, sysUser.ID, endTime); err != nil {
			log.New(ctx).Error("Fail to update start time", "sysUserID", sysUser.ID, "accountName", sysUser.Username, "error", err)
		}
	}
}

// 计算发送成功率
func (svc *EmailResultService) calculateRate(ctx context.Context, sysUser *model.EmlSysUser, startTime, endTime time.Time, results []*model.OpsEmlResult) {
	log.New(ctx).Info("start to calculateRate", "accountName", sysUser.Username, "results count", len(results), "startTime", startTime, "endTime", endTime)

	// 统计出 (domain, success_cnt, total)
	type domain struct {
		Domain       string
		SuccessCnt   int32
		TemporaryErr int32
		PermanentErr int32
		TotalCnt     int32
	}
	m := make(map[string]domain)
	for _, result := range results {
		do, ok := m[result.Domain]
		if !ok {
			do = domain{
				Domain:       result.Domain,
				SuccessCnt:   0,
				TemporaryErr: 0,
				PermanentErr: 0,
				TotalCnt:     0,
			}
		}
		do.TotalCnt++
		smtpCode := enum.GetSmtpCode(result.Message)
		switch smtpCode.Category {
		case enum.Success:
			do.SuccessCnt++
		case enum.TemporaryFailure:
			do.TemporaryErr++
		case enum.PermanentFailure:
			do.PermanentErr++
		default:
			do.PermanentErr++
			log.New(ctx).Warn("unknown smtp code category", "result.ID", result.ID, "result.AccountName", result.AccountName, "result.ToAddress", result.ToAddress, "smtpCode.Category", smtpCode.Category, "message", result.Message)
		}

		m[result.Domain] = do
	}

	// 处理每一个(domain, success_cnt, total)
	for _, do := range m {
		// 针对每一个(domain, success_cnt, total) ，并存入eml_success_rate表
		entity := &model.OpsEmlNumberSuccess{
			Cid:          sysUser.Cid,
			FromAddress:  sysUser.FromAddress,
			Domain:       do.Domain,
			StartTime:    startTime,
			EndTime:      endTime,
			SuccessCnt:   do.SuccessCnt,
			TemporaryErr: do.TemporaryErr,
			PermanentErr: do.PermanentErr,
			TotalCnt:     do.TotalCnt,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err := svc.OpsEmlNumberSuccessRepo.Create(ctx, entity)
		if err != nil {
			log.New(ctx).Error("fail to create ops_eml_number_success", "accountName", sysUser.Username, "error", err, "domain", do.Domain, "successCnt", do.SuccessCnt, "totalCnt", do.TotalCnt)
			continue
		}

		// 如果过低就减速
		// rate = 成功数 / （成功数 + 永久错误数）
		if do.TotalCnt == 0 {
			log.New(ctx).Warn("total is zero", "accountName", sysUser.Username, "domain", do.Domain, "successCnt", do.SuccessCnt, "permanentErr", do.PermanentErr)
			continue
		}
		rate := int(do.SuccessCnt * 100 / do.TotalCnt)
		log.New(ctx).Info("Calculated success rate", "accountName", sysUser.Username, "domain", do.Domain, "rate", rate)
		if rate < config.App.DownRate {
			log.New(ctx).Info("domain slow down", "accountName", sysUser.Username, "domain", do.Domain)
			credibility, err := svc.EmailDomainCredibilityRepo.Read(ctx, sysUser.Cid, do.Domain, sysUser.FromAddress)
			if err != nil {
				log.New(ctx).Error("fail to read email domain credibility", "accountName", sysUser.Username, "domain", do.Domain, "error", err)
				continue
			}
			oldSpeed := credibility.Speed
			credibility.Speed = enum.SlowDown(oldSpeed)
			if _, err = svc.EmailDomainCredibilityRepo.UpdateSpeed(ctx, credibility); err != nil {
				log.New(ctx).Error("fail to update email domain credibility", "accountName", sysUser.Username, "domain", do.Domain, "error", err)
			}
			log.New(ctx).Info("update email domain credibility", "accountName", sysUser.Username, "domain", do.Domain, "newSpeed", credibility.Speed, "oldSpeed", oldSpeed)
		} else {
			log.New(ctx).Info("Success rate greater than downRate, no need to slow down", "rate", rate, "downRate", config.App.DownRate)
		}
	}
}

func (svc *EmailResultService) processResp(ctx context.Context, resp *dm20151123.SenderStatisticsDetailByParamResponse, results *[]*dm20151123.SenderStatisticsDetailByParamResponseBodyDataMailDetail) (string, error) {
	log.New(ctx).Info(fmt.Sprintf("%d email results successfully fetched", len(resp.Body.Data.MailDetail)))
	for _, detail := range resp.Body.Data.MailDetail {
		// 将写入数据库操作替换成写内存
		*results = append(*results, detail)
	}
	return *resp.Body.NextStart, nil
}

func detailToResult(sendTime int64, detail *dm20151123.SenderStatisticsDetailByParamResponseBodyDataMailDetail) *model.OpsEmlResult {
	return &model.OpsEmlResult{
		AccountName:         *detail.AccountName,
		Domain:              strings.SplitN(*detail.ToAddress, "@", 2)[1],
		ErrorClassification: *detail.ErrorClassification,
		SentTime:            time.UnixMilli(sendTime),
		Message:             *detail.Message,
		Status:              *detail.Status,
		Subject:             *detail.Subject,
		ToAddress:           *detail.ToAddress,
	}
}

func createClient(sysUser *model.EmlSysUser) (_result *dm20151123.Client, err error) {
	credential, err := credential.NewCredential(&credential.Config{
		Type:            tea.String("access_key"),
		AccessKeyId:     tea.String(sysUser.AccessKeyID),
		AccessKeySecret: tea.String(sysUser.AccessKeySecret),
	})
	if err != nil {
		return _result, err
	}

	config := &openapi.Config{
		Credential: credential,
	}
	config.Endpoint = tea.String("dm.ap-southeast-1.aliyuncs.com")
	_result = &dm20151123.Client{}
	_result, err = dm20151123.NewClient(config)
	return _result, err
}
