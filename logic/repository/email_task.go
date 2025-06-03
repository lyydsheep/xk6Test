package repository

import (
	"context"
	"email/common/enum"
	"email/common/log"
	"email/dal/dao"
	"email/dal/model"
	"errors"
	"time"
)

type EmailTaskRepoV1 struct {
	Query *dao.Query
}

func (repo *EmailTaskRepoV1) CreateWithTx(ctx context.Context, task *model.EmlTask, tx *dao.QueryTx) error {
	return tx.EmlTask.WithContext(ctx).Create(task)
}

func (repo *EmailTaskRepoV1) UpdateEmailTask(ctx context.Context, taskID int64, status string, successNumDelta int32, pendingNumDelta int32) error {
	var (
		err       error
		errCount  int
		emailTask *model.EmlTask
	)
	for errCount < 3 {
		emailTask, err = repo.ReadByID(ctx, taskID)
		if err != nil {
			errCount += 1
			log.New(ctx).Error("fail to read emailTask when update emailTask", "taskID", taskID, "error", err)
			continue
		}
		if status != "" {
			emailTask.Status = status
		}
		version := emailTask.Version
		emailTask.SuccessNum += successNumDelta
		emailTask.PendingNum += pendingNumDelta
		emailTask.Version += 1
		row, err := repo.Update(ctx, emailTask, version)

		if err != nil {
			errCount += 1
			log.New(ctx).Error("fail to update emailTask when update emailTask", "taskID", taskID, "error", err)
		}
		if row > 0 {
			log.New(ctx).Warn("emailTask has been taken when update emailTask", "taskID", taskID)
			return nil
		}
		time.Sleep(time.Second * time.Duration(1<<errCount))
	}

	log.New(ctx).Info("success to update emailTask when update emailTask", "taskID", taskID, "taskSubID")
	return err
}

func (repo *EmailTaskRepoV1) ReadByID(ctx context.Context, taskID int64) (*model.EmlTask, error) {
	return repo.Query.EmlTask.WithContext(ctx).Where(repo.Query.EmlTask.ID.Eq(taskID)).First()
}

func (repo *EmailTaskRepoV1) GetEmailTask(ctx context.Context) (*model.EmlTask, error) {
	// 获取 pending 状态的任务
	// 获取 processing 状态的任务 && 超时 15 分钟
	tasks, err := repo.Query.EmlTask.WithContext(ctx).Where(repo.Query.EmlTask.Status.Eq(enum.EmailTaskStatusPending)).
		Or(repo.Query.EmlTask.Status.Eq(enum.EmailTaskStatusProcessing), repo.Query.EmlTask.FetchTime.Lt(time.Now().UTC().Add(-15*time.Hour))).Limit(1).Find()
	if err != nil {
		return nil, errors.New("emailTask " + err.Error())
	}
	if len(tasks) == 0 {
		log.New(ctx).Debug("no task to fetch")
		return nil, nil
	}
	task := tasks[0]
	log.New(ctx).Debug("fetch task", "taskID", task.ID)

	// 乐观锁
	version := task.Version
	task.Version = version + 1
	task.Status = enum.EmailTaskStatusProcessing
	now := time.Now().UTC()
	task.FetchTime = &now

	result, err := repo.Query.EmlTask.WithContext(ctx).
		Where(repo.Query.EmlTask.ID.Eq(task.ID), repo.Query.EmlTask.Version.Eq(version)).Updates(task)
	if err != nil {
		log.New(ctx).Error("fail to update task", "taskID", task.ID, "error", err)
		return nil, errors.New("emailTask " + err.Error())
	}
	if result.RowsAffected > 0 {
		// 成功获取锁
		log.New(ctx).Debug("success to fetch task", "taskID", task.ID)
		return task, nil
	}
	return nil, nil
}

func (repo *EmailTaskRepoV1) Update(ctx context.Context, task *model.EmlTask, version int32) (int64, error) {
	result, err := repo.Query.EmlTask.WithContext(ctx).Where(repo.Query.EmlTask.ID.Eq(task.ID),
		repo.Query.EmlTask.Version.Eq(version)).Updates(map[string]interface{}{
		"status":       task.Status,
		"fetch_time":   task.FetchTime,
		"success_num":  task.SuccessNum,
		"fail_num":     task.FailNum,
		"invalid_num":  task.InvalidNum,
		"pending_num":  task.PendingNum,
		"open_num":     task.OpenNum,
		"open_num_de":  task.OpenNumDe,
		"click_num":    task.ClickNum,
		"click_num_de": task.ClickNumDe,
		"version":      task.Version,
	})
	if err != nil {
		return 0, errors.New("emailTask " + err.Error())
	}
	return result.RowsAffected, nil
}

func NewEmailTaskRepoV1(q *dao.Query) EmailTaskRepo {
	return &EmailTaskRepoV1{
		Query: q,
	}
}
