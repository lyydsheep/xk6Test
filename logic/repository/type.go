package repository

import (
	"context"
	"email/dal/dao"
	"email/dal/model"
	"email/event"
	"gorm.io/gen"
	"time"
)

type EmailTaskSubRepo interface {
	Create(ctx context.Context, taskSub *model.EmlTaskSub) (*model.EmlTaskSub, error)
	Update(ctx context.Context, version int32, taskSub *model.EmlTaskSub) (int64, error)
	GetBatchByFromEmail(ctx context.Context, fromEmail string) ([]*model.EmlTaskSub, error)
	Read(ctx context.Context, taskID int64, version int32) (int64, error)
	HandleFinal(ctx context.Context, taskSubID int64, status string, taskID int64, processInfo string, failNumDelta int32, invalidNumDelta int32, pendingNumDelta int32, successNumDelta int32) error
	CreateWithEvent(ctx context.Context, paid *event.PaidEvent, priority int32) (bool, error)
	ReadDoneAndPredoneAndOutdated(ctx context.Context, fromEmail string, endTime time.Time) ([]*model.EmlTaskSub, error)
	ReadInitCount(ctx context.Context, fromEmail string) (int64, error)
	CreateWithTx(ctx context.Context, task *model.EmlTaskSub, tx *dao.QueryTx) error
}

type EmailTemplateRepo interface {
	Read(ctx context.Context, cid int64, templateID int64) (*model.EmlTemplate, error)
	ReadPriority(ctx context.Context, templateID int64) (int32, error)
	ReadMaxRetry(ctx context.Context, templateID int64) (int32, error)
}

type UsrUserRepo interface {
	ReadTagUser(ctx context.Context, cid int64, tag string) ([]*model.UsrUser, error)
	Read(ctx context.Context, cid int64, userID string) (*model.UsrUser, error)
}

type EmailSysUserRepo interface {
	Read(ctx context.Context, cid int64, fromEmail string) (*model.EmlSysUser, error)
	ReadByType(ctx context.Context, cid int64, _type string) (*model.EmlSysUser, error)
	ReadAll(ctx context.Context) ([]*model.EmlSysUser, error)
	UpdateStartTime(ctx context.Context, id int64, startTime time.Time) error
}

type EmailTaskRepo interface {
	Update(ctx context.Context, task *model.EmlTask, version int32) (int64, error)
	GetEmailTask(ctx context.Context) (*model.EmlTask, error)
	ReadByID(ctx context.Context, taskID int64) (*model.EmlTask, error)
	UpdateEmailTask(ctx context.Context, taskID int64, status string, successNumDelta int32, pendingNumDelta int32) error
	CreateWithTx(ctx context.Context, task *model.EmlTask, tx *dao.QueryTx) error
}

type EmailRedemptionCodeRepo interface {
	Read(ctx context.Context, taskSub *model.EmlTaskSub, codeType string, amount int32) (*model.EmlRedemptionCode, error)
}

type EmailRedemptionTaskSubRepo interface {
	Read(ctx context.Context, CID int64, TaskSubID int64) (*model.EmlRedemptionTaskSub, error)
}

type EmailUnsubscribeUsrUserRepo interface {
	Read(ctx context.Context, cid int64, uid string) (*model.EmlUnsubscribeUsrUser, error)
	Create(ctx context.Context, entity *model.EmlUnsubscribeUsrUser) (*model.EmlUnsubscribeUsrUser, error)
}

type EmailDomainCredibilityRepo interface {
	Create(ctx context.Context, entity *model.EmlDomainCredibility) (*model.EmlDomainCredibility, error)
	Read(ctx context.Context, cid int64, domain string, fromAddress string) (*model.EmlDomainCredibility, error)
	UpdateLastSentTime(ctx context.Context, entity *model.EmlDomainCredibility) (gen.ResultInfo, error)
	UpdateSpeed(ctx context.Context, entity *model.EmlDomainCredibility) (gen.ResultInfo, error)
	ReadAll(ctx context.Context) ([]*model.EmlDomainCredibility, error)
}

type TaskConfigRepo interface {
	Create(ctx context.Context, entity *model.EmlTaskConfig) (*model.EmlTaskConfig, error)
	Read(ctx context.Context, cid int64, category string, _type string) (*model.EmlTaskConfig, error)
	Update(ctx context.Context, entity *model.EmlTaskConfig) (gen.ResultInfo, error)
	ReadAll(ctx context.Context) ([]*model.EmlTaskConfig, error)
}

type EmailResultRepo interface {
	Create(ctx context.Context, entity *model.OpsEmlResult) error
	Update(ctx context.Context, entity *model.OpsEmlResult) (gen.ResultInfo, error)
	ReadByTime(ctx context.Context, fromAddress string, domain string, startTime time.Time, endTime time.Time) ([]*model.OpsEmlResult, error)
	GetEndTime(ctx context.Context, fromAddress string) (time.Time, error)
	ReadByTo(ctx context.Context, toAddress string) ([]*model.OpsEmlResult, error)
	ReadBatchByRange(ctx context.Context, accountName string, startTime time.Time, endTime time.Time) ([]*model.OpsEmlResult, error)
}

type OpsEmlNumberSuccessRepo interface {
	Create(ctx context.Context, entity *model.OpsEmlNumberSuccess) error
	Read(ctx context.Context, cid int64, fromAddress string, startTime time.Time, endTime time.Time) ([]*model.OpsEmlNumberSuccess, error)
}
