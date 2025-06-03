package repository

import (
	"context"
	"email/common/log"
	"email/dal/dao"
	"email/dal/model"
)

type EmailRedemptionCodeV1 struct {
	Query *dao.Query
}

// 获取兑换码, codeType 用于区分不同面值、功能的兑换码
func (repo *EmailRedemptionCodeV1) Read(ctx context.Context, taskSub *model.EmlTaskSub, codeType string, amount int32) (*model.EmlRedemptionCode, error) {
	// 如果已经存在，则直接返回
	entity, err := repo.Query.EmlRedemptionCode.WithContext(ctx).Where(repo.Query.EmlRedemptionCode.Cid.Eq(taskSub.Cid),
		repo.Query.EmlRedemptionCode.Type.Eq(codeType), repo.Query.EmlRedemptionCode.Amount.Eq(amount), repo.Query.EmlRedemptionCode.TaskSubID.Eq(taskSub.ID)).First()
	if err == nil {
		log.New(ctx).Debug("code already exists", "code", entity)
		return entity, err
	}
	tx := repo.Query.Begin()
	defer tx.Rollback()
	// 没有被发送的兑换码 ---> task_sub_id = 0
	entity, err = tx.EmlRedemptionCode.WithContext(ctx).Where(tx.EmlRedemptionCode.Cid.Eq(taskSub.Cid),
		tx.EmlRedemptionCode.Type.Eq(codeType), tx.EmlRedemptionCode.Amount.Eq(amount), tx.EmlRedemptionCode.TaskSubID.Eq(0)).First()
	if err != nil {
		return nil, err
	}
	entity.TaskSubID = taskSub.ID
	// TODO where 更新
	info, err := tx.EmlRedemptionCode.WithContext(ctx).Where(tx.EmlRedemptionCode.ID.Eq(entity.ID),
		tx.EmlRedemptionCode.TaskSubID.Eq(0)).Updates(entity)
	if err != nil {
		return nil, err
	}
	if info.RowsAffected > 0 {
		tx.Commit()
		return entity, nil
	}
	return entity, err
}

func NewEmailRedemptionCodeV1(query *dao.Query) EmailRedemptionCodeRepo {
	return &EmailRedemptionCodeV1{
		Query: query,
	}
}
