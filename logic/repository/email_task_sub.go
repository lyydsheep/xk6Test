package repository

import (
	"context"
	"email/common/enum"
	"email/common/log"
	"email/dal/dao"
	"email/dal/model"
	"email/event"
	"errors"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

const (
	YearlyPro      = "yearly-pro-continuous"
	MonthlyPro     = "monthly-pro-continuous"
	YearlyPremium  = "yearly-premium-continuous"
	MonthlyPremium = "monthly-premium-continuous"
	Credit         = "credit"
)

type amount struct {
	Type   string
	Amount int32
}

var (
	Pro = amount{
		Type:   "credit",
		Amount: 300,
	}
	Premium = amount{
		Type:   "credit",
		Amount: 1200,
	}
)

type TaskSubRepoV1 struct {
	Query *dao.Query
}

func (repo *TaskSubRepoV1) CreateWithTx(ctx context.Context, task *model.EmlTaskSub, tx *dao.QueryTx) error {
	return tx.EmlTaskSub.WithContext(ctx).Create(task)
}

func (repo *TaskSubRepoV1) ReadInitCount(ctx context.Context, fromEmail string) (int64, error) {
	qTaskSub := repo.Query.EmlTaskSub
	return qTaskSub.WithContext(ctx).Where(qTaskSub.FromEmail.Eq(fromEmail),
		qTaskSub.Status.Eq(enum.EmailTaskSubStatusUnsent)).Count()
}

func (repo *TaskSubRepoV1) ReadDoneAndPredoneAndOutdated(ctx context.Context, fromEmail string, endTime time.Time) ([]*model.EmlTaskSub, error) {
	emlTaskSub := repo.Query.EmlTaskSub
	// 获取已发送的邮件
	res, err := emlTaskSub.WithContext(ctx).Where(emlTaskSub.FromEmail.Eq(fromEmail), emlTaskSub.CodeSMTP.Eq(""),
		emlTaskSub.Status.Eq(enum.EmailTaskSubStatusSent), emlTaskSub.FetchTime.Lte(endTime)).Limit(1).Find()
	if err != nil || len(res) > 0 {
		return res, err
	}

	// 获取 predone 邮件
	res, err = emlTaskSub.WithContext(ctx).Where(emlTaskSub.FromEmail.Eq(fromEmail), emlTaskSub.CodeSMTP.Eq(""),
		emlTaskSub.FetchTime.Lte(endTime), emlTaskSub.Status.Eq(enum.EmailTaskSubStatusPreDone)).Limit(1).Find()
	if err != nil || len(res) > 0 {
		return res, err
	}

	// 获取超时邮件
	return emlTaskSub.WithContext(ctx).Where(emlTaskSub.FromEmail.Eq(fromEmail), emlTaskSub.Status.Eq(enum.EmailTaskSubStatusProcess), emlTaskSub.FetchTime.Lte(time.Now().Add(-24*time.Hour))).Limit(1).Find()
}

func (repo *TaskSubRepoV1) ReadByFromAndTo(ctx context.Context, cid int64, fromEmail, toEmail string) (*model.EmlTaskSub, error) {
	return repo.Query.EmlTaskSub.WithContext(ctx).Where(repo.Query.EmlTaskSub.Cid.Eq(cid), repo.Query.EmlTaskSub.FromEmail.Eq(fromEmail),
		repo.Query.EmlTaskSub.ToEmail.Eq(toEmail), repo.Query.EmlTaskSub.CodeSMTP.Eq(""),
		repo.Query.EmlTaskSub.Status.In(enum.EmailTaskSubStatusSent, enum.EmailTaskSubStatusPreDone)).Order(repo.Query.EmlTaskSub.ID.Desc()).First()
}

func (repo *TaskSubRepoV1) ReadByStatusAndCode(ctx context.Context, status, code string, limit int) ([]*model.EmlTaskSub, error) {
	return repo.Query.EmlTaskSub.WithContext(ctx).Where(repo.Query.EmlTaskSub.CodeSMTP.Eq(code), repo.Query.EmlTaskSub.Status.Eq(status)).
		Or(repo.Query.EmlTaskSub.CodeSMTP.Eq(code), repo.Query.EmlTaskSub.Status.Eq(enum.EmailTaskSubStatusPreDone), repo.Query.EmlTaskSub.FetchTime.Lt(time.Now().Add(-time.Hour))).Order(repo.Query.EmlTaskSub.ID.Desc()).Limit(limit).Find()
}

func (repo *TaskSubRepoV1) ReadByTime(ctx context.Context, cid int64, fromEmail, domain string, start time.Time, end time.Time) ([]*model.EmlTaskSub, error) {
	return repo.Query.EmlTaskSub.WithContext(ctx).Where(repo.Query.EmlTaskSub.CodeSMTP.Neq(""), repo.Query.EmlTaskSub.FromEmail.Eq(fromEmail),
		repo.Query.EmlTaskSub.Cid.Eq(cid), repo.Query.EmlTaskSub.EmailDomain.Eq(domain), repo.Query.EmlTaskSub.SentTime.Between(start, end)).Find()
}

// TODO templateID 这里写死
func (repo *TaskSubRepoV1) CreateWithEvent(ctx context.Context, paid *event.PaidEvent, priority int32) (bool, error) {
	tx := repo.Query.Begin()
	defer tx.Rollback()
	now := time.Now()
	task := &model.EmlTask{
		Cid:        paid.Account.Cid,
		Status:     enum.EmailTaskStatusFinished,
		FetchTime:  &now,
		PendingNum: 1,
		TotalNum:   1,
		TemplateID: 2,
		Type:       enum.TaskConfigRedemption,
	}
	err := tx.EmlTask.WithContext(ctx).Create(task)
	if err != nil {
		log.New(ctx).Error("[createCodeTask]fail to create task when CreateWithEvent", "paidEvent", paid, "error", err)
		return false, err
	}
	sysUser, err := tx.EmlSysUser.WithContext(ctx).Where(tx.EmlSysUser.Cid.Eq(paid.User.Cid), tx.EmlSysUser.Type.Eq(enum.SysUserRedemption)).First()
	if err != nil {
		log.New(ctx).Error("[createCodeTask]fail to read sysUser when CreateWithEvent", "cId", paid.User.Cid, "error", err)
		return false, err
	}
	taskSub := &model.EmlTaskSub{
		Cid: paid.User.Cid,
		// TODO 这个-1是一个临时方案，后续需要修改
		TaskID:    task.ID,
		FromEmail: sysUser.FromAddress,
		ToUserID:  paid.User.Uid,
		ToEmail:   paid.User.Email,
		// TODO 将兑换码邮件模板写死为 2
		TemplateID: 2,
		Status:     enum.EmailTaskSubStatusUnsent,
		Version:    0,
		Type:       enum.TaskConfigRedemption,
		Retry:      0,
		Priority:   priority,
	}
	err = tx.EmlTaskSub.WithContext(ctx).Create(taskSub)
	if err != nil {
		log.New(ctx).Error("[createCodeTask]fail to create taskSub when CreateWithEvent", "paidEvent", paid, "error", err)
		return false, err
	}
	redemptionTaskSub := &model.EmlRedemptionTaskSub{
		Cid:       paid.User.Cid,
		TaskSubID: taskSub.ID,
		PaymentID: paid.Account.PaymentId,
	}
	if len(paid.GoodItems) == 0 {
		return false, errors.New("[createCodeTask]paidEvent.GoodItems is empty")
	}
	// TODO 根据商品类型，生成对应的兑换码类型
	switch paid.GoodItems[0].GoodsId {
	case MonthlyPro, YearlyPro:
		redemptionTaskSub.Type = Pro.Type
		redemptionTaskSub.Amount = Pro.Amount
	case MonthlyPremium, YearlyPremium:
		redemptionTaskSub.Type = Premium.Type
		redemptionTaskSub.Amount = Premium.Amount
	default:
		log.New(ctx).Error("[createCodeTask]unknown goodsId", "paidEvent", paid)
		return true, errors.New("unknown goodsId")
	}
	// 唯一索引冲突会报错  --->  回滚
	err = tx.EmlRedemptionTaskSub.WithContext(ctx).Create(redemptionTaskSub)
	if err != nil {
		log.New(ctx).Error("[createCodeTask]fail to create taskSub when CreateWithEvent", "paidEvent", paid, "error", err)
		return false, err
	}
	return false, tx.Commit()
}

func (repo *TaskSubRepoV1) HandleFinal(ctx context.Context, taskSubID int64, status string, taskID int64, processInfo string, failNumDelta int32, invalidNumDelta int32, pendingNumDelta int32, successNumDelta int32) error {
	tx := repo.Query.Begin()
	defer tx.Rollback()
	var (
		taskSub *model.EmlTaskSub
		err     error
		//task    *model.EmlTask
		info gen.ResultInfo
	)
	if taskSubID != 0 {
		taskSub, err = tx.EmlTaskSub.WithContext(ctx).Where(tx.EmlTaskSub.ID.Eq(taskSubID)).First()
		if err != nil {
			log.New(ctx).Error("fail to read taskSub when handleFinal", "taskSubID", taskSubID, "error", err)
			return err
		}
		taskSub.Status = status
		taskSub.Version++
		info, err = tx.EmlTaskSub.WithContext(ctx).Where(tx.EmlTaskSub.ID.Eq(taskSub.ID)).Updates(taskSub)
		if err != nil {
			log.New(ctx).Error("fail to update taskSub", "taskSubID", taskSub.ID, "error", err)
			return errors.New("taskSub " + err.Error())
		}
		if info.RowsAffected == 0 {
			log.New(ctx).Error("fail to update taskSub", "taskSubID", taskSub.ID, "affectedRows", info.RowsAffected)
			return errors.New("fail to update emailTaskSub when handleFinal")
		}
	} else {
		taskSub = &model.EmlTaskSub{ID: 0}
	}
	return nil
	//
	//errCount := 0
	//for errCount < 5 {
	//	task, err = tx.EmlTask.WithContext(ctx).Where(tx.EmlTask.ID.Eq(taskID)).First()
	//	if err != nil {
	//		if errors.Is(err, gorm.ErrRecordNotFound) {
	//			return errors.New("taskSub " + err.Error())
	//		}
	//		errCount++
	//		continue
	//	}
	//	version := task.Version
	//	{
	//		task.Version++
	//		task.FailNum += failNumDelta
	//		task.InvalidNum += invalidNumDelta
	//		task.PendingNum += pendingNumDelta
	//		task.SuccessNum += successNumDelta
	//		if processInfo != "" {
	//			task.ProcessInfo += fmt.Sprintf("taskSub.ID: %d  ProcessInfo: %s\n", taskSubID, processInfo)
	//		}
	//	}
	//	info, err = tx.EmlTask.WithContext(ctx).Where(tx.EmlTask.ID.Eq(taskID), tx.EmlTask.Version.Eq(version)).Updates(map[string]interface{}{
	//		"version":      task.Version,
	//		"fail_num":     task.FailNum,
	//		"invalid_num":  task.InvalidNum,
	//		"pending_num":  task.PendingNum,
	//		"process_info": task.ProcessInfo,
	//		"status":       task.Status,
	//		"success_num":  task.SuccessNum,
	//	})
	//	if err != nil {
	//		errCount++
	//		continue
	//	}
	//	if info.RowsAffected > 0 {
	//		tx.Commit()
	//		log.New(ctx).Debug("update taskSub success", "taskID", taskID, "affectedRows", info.RowsAffected)
	//		return nil
	//	}
	//	err = errors.New("rowsAffected is 0")
	//	errCount++
	//	time.Sleep((time.Duration(rand.Float32())*2 + time.Duration(1<<errCount)) * time.Second)
	//}
	//
	//log.New(ctx).Error("fail to handle fail", "taskID", taskID, "taskSub.ID", taskSubID, "err", err)
	//return errors.New("taskSub " + err.Error())
}

func (repo *TaskSubRepoV1) Read(ctx context.Context, taskID int64, version int32) (int64, error) {
	_, err := repo.Query.EmlTaskSub.WithContext(ctx).Where(repo.Query.EmlTaskSub.ID.Eq(taskID),
		repo.Query.EmlTaskSub.Version.Eq(version)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return 1, nil
}

// FetchBatch implements EmailTaskSubRepo.
func (repo *TaskSubRepoV1) GetBatchByFromEmail(ctx context.Context, fromEmail string) ([]*model.EmlTaskSub, error) {
	// 首次 兑换码 > 营销
	emlTaskSub := repo.Query.EmlTaskSub
	taskSubs, err := emlTaskSub.WithContext(ctx).Where(emlTaskSub.Status.Eq(enum.EmailTaskSubStatusUnsent),
		emlTaskSub.Retry.Eq(0), emlTaskSub.FromEmail.Eq(fromEmail)).
		Order(emlTaskSub.Priority.Asc()).Limit(1).Find()
	if err != nil {
		return nil, err
	}
	if len(taskSubs) == 0 {
		// 重试 次数少 > 次数多
		taskSubs, err = emlTaskSub.WithContext(ctx).Where(emlTaskSub.Status.Eq(enum.EmailTaskSubStatusUnsent),
			emlTaskSub.Retry.Gt(0), emlTaskSub.FromEmail.Eq(fromEmail)).
			Order(emlTaskSub.Retry.Asc()).Limit(1).Find()
	}
	if err != nil {
		return nil, err
	}
	taskSubList := make([]*model.EmlTaskSub, 0, len(taskSubs))

	// 乐观锁校验
	for i := range taskSubs {
		version := taskSubs[i].Version
		// 更新状态
		taskSubs[i].Version++
		taskSubs[i].Status = enum.EmailTaskSubStatusProcess
		taskSubs[i].FetchTime = time.Now()
		result, err := repo.Query.EmlTaskSub.WithContext(ctx).
			Where(repo.Query.EmlTaskSub.ID.Eq(taskSubs[i].ID), repo.Query.EmlTaskSub.Version.Eq(version)).Updates(taskSubs[i])
		if err != nil {
			log.New(ctx).Error("fail to update taskSub", "taskSubID", taskSubs[i].ID, "error", err)
			continue
		}
		if result.RowsAffected > 0 {
			// 有效更新
			taskSubList = append(taskSubList, taskSubs[i])
		}
	}
	return taskSubList, nil
}

// Create 需要基于幂等操作
func (repo *TaskSubRepoV1) Create(ctx context.Context, taskSub *model.EmlTaskSub) (*model.EmlTaskSub, error) {
	// 唯一索引(taskID, toEmail)
	// 幂等操作，如果发生冲突，则什么都不做
	err := repo.Query.EmlTaskSub.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(taskSub)

	if err != nil {
		return nil, errors.New("taskSub " + err.Error())
	}
	return taskSub, err
}

func (repo *TaskSubRepoV1) Update(ctx context.Context, version int32, taskSub *model.EmlTaskSub) (int64, error) {
	info, err := repo.Query.EmlTaskSub.WithContext(ctx).Where(repo.Query.EmlTaskSub.ID.Eq(taskSub.ID),
		repo.Query.EmlTaskSub.Version.Eq(version)).Updates(taskSub)
	if err != nil {
		return 0, err
	}
	return info.RowsAffected, nil
}

func NewTaskSubRepoV1(q *dao.Query) EmailTaskSubRepo {
	return &TaskSubRepoV1{
		Query: q,
	}
}
