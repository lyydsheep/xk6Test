package service

import (
	"context"
	"email/common/enum"
	"email/common/log"
	"email/config"
	"email/dal/model"
	"email/logic/repository"
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

type AliFlowControlService struct {
	EmailDomainCredibilityRepo repository.EmailDomainCredibilityRepo
	EmailTaskSubRepo           repository.EmailTaskSubRepo
	EmailResultRepo            repository.EmailResultRepo
	EmailSysUserRepo           repository.EmailSysUserRepo
	OpsEmlNumberSuccessRepo    repository.OpsEmlNumberSuccessRepo
}

func NewFlowControlService(emailDomainCredibilityRepo repository.EmailDomainCredibilityRepo, emailTaskSubRepo repository.EmailTaskSubRepo,
	emailResultRepo repository.EmailResultRepo, emailSysUserRepo repository.EmailSysUserRepo, opsEmlNumberSuccessRepo repository.OpsEmlNumberSuccessRepo) *AliFlowControlService {
	return &AliFlowControlService{
		EmailDomainCredibilityRepo: emailDomainCredibilityRepo,
		EmailTaskSubRepo:           emailTaskSubRepo,
		EmailResultRepo:            emailResultRepo,
		EmailSysUserRepo:           emailSysUserRepo,
		OpsEmlNumberSuccessRepo:    opsEmlNumberSuccessRepo,
	}
}

func (svc *AliFlowControlService) Run(ctx context.Context) {
	m := map[string]struct{}{}
	for {
		sysUsers, err := svc.EmailSysUserRepo.ReadAll(ctx)
		if err != nil {
			log.New(ctx).Error("Fail to read all sysUsers", "error", err)
			return
		}
		// TODO 默认所有都是阿里云 api，后续得支持更灵活的配置信息
		// TODO goroutine 泄露风险
		for _, sysUser := range sysUsers {
			if sysUser.EmailServer == enum.Ali {
				key := fmt.Sprintf("%dkey%s", sysUser.Cid, sysUser.Username)
				if _, ok := m[key]; !ok {
					// 不存在
					log.New(ctx).Info("start to consume taskSub", "accountName", sysUser.Username)
					m[key] = struct{}{}
					go svc.work(ctx, sysUser)
				}
			}
		}
		time.Sleep(time.Hour)
	}
}

// 两个定时任务
func (svc *AliFlowControlService) work(ctx context.Context, sysUser *model.EmlSysUser) {
	c := cron.New()
	// 每天 8 点和 20点执行一次加速操作
	_, err := c.AddFunc(config.App.SpeedUpTime, func() {
		// 扫描eml_success_rate表 12小时内的记录
		svc.speedUp(ctx, sysUser)
	})
	if err != nil {
		log.New(ctx).Error("Failed to add cron job", "error", err)
		return
	}

	c.Start()
	defer c.Stop()
	// 防止 goroutine 退出
	select {}
}

// 加速逻辑实现
func (svc *AliFlowControlService) speedUp(ctx context.Context, sysUser *model.EmlSysUser) {
	// 扫描eml_success_rate表 12小时内的记录
	startTime, endTime := time.Now().Add(-12*time.Hour), time.Now()
	successRecords, err := svc.OpsEmlNumberSuccessRepo.Read(ctx, sysUser.Cid, sysUser.FromAddress, startTime, endTime)
	if err != nil {
		log.New(ctx).Error("fail to read ops_eml_number_success", "accountName", sysUser.Username, "startTime", startTime, "endTime", endTime, "error", err)
		return
	}
	if len(successRecords) == 0 {
		log.New(ctx).Info("No success records found, skipping speed up", "accountName", sysUser.Username)
		return
	}

	log.New(ctx).Info("Successfully fetched the NumberSuccess record", "accountName", sysUser.Username, "count", len(successRecords), "startTime", startTime, "endTime", endTime)
	// 统计 (domain, eml_success_rate_cnt)
	memo := map[string]struct {
		recordCnt       int32
		successCnt      int32
		temporaryErrCnt int32
		permanentErrCnt int32
		totalCnt        int32
		domain          string
	}{}
	for _, record := range successRecords {
		val := memo[record.Domain]
		memo[record.Domain] = struct {
			recordCnt       int32
			successCnt      int32
			temporaryErrCnt int32
			permanentErrCnt int32
			totalCnt        int32
			domain          string
		}{
			recordCnt:       val.recordCnt + 1,
			successCnt:      val.successCnt + record.SuccessCnt,
			temporaryErrCnt: val.temporaryErrCnt + record.TemporaryErr,
			permanentErrCnt: val.permanentErrCnt + record.PermanentErr,
			totalCnt:        val.totalCnt + record.TotalCnt,
			domain:          record.Domain,
		}
	}
	// 理想情况每一个 domain 会有 12 * 60 / 5 = 144 条记录
	for domain, val := range memo {
		log.New(ctx).Info("get domain information", "accountName", sysUser.Username, "domain", domain, "recordCnt", val.recordCnt, "successCnt", val.successCnt, "totalCnt", val.totalCnt)
		// 如果少于 36个，就不加速
		if val.recordCnt < 36 {
			log.New(ctx).Info("recordCnt not enough. cancellation of acceleration", "accountName", sysUser.Username, "domain", domain, "recordCnt", val.recordCnt)
			continue
		}

		// 		计算平均成功率，＞ 98%就加速
		// rate = 成功数 / （成功数 + 永久错误数）
		if val.totalCnt == 0 {
			log.New(ctx).Info("total is 0. cancellation of acceleration", "accountName", sysUser.Username, "domain", domain)
			continue
		}
		rate := int(val.successCnt * 100 / val.totalCnt)
		log.New(ctx).Info("Calculated success rate", "accountName", sysUser.Username, "domain", domain, "rate", rate)
		if rate > config.App.UpRate {
			credibility, err := svc.EmailDomainCredibilityRepo.Read(ctx, sysUser.Cid, domain, sysUser.FromAddress)
			if err != nil {
				log.New(ctx).Error("fail to read domain credibility", "domain", domain, "accountName", sysUser.Username)
				continue
			}
			oldSpeed := credibility.Speed
			credibility.Speed = enum.SpeedUp(credibility.Speed)
			_, err = svc.EmailDomainCredibilityRepo.UpdateSpeed(ctx, credibility)
			if err != nil {
				log.New(ctx).Error("fail to update domain credibility", "accountName", sysUser.Username, "domain", domain)
				continue
			}
			log.New(ctx).Info("domain speed up", "accountName", sysUser.Username, "domain", domain, "oldSpeed", oldSpeed, "newSpeed", credibility.Speed)
		}
	}
}
