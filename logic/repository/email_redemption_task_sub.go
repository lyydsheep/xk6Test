package repository

import (
	"context"
	"email/dal/dao"
	"email/dal/model"
)

type EmailRedemptionTaskSubRepoV1 struct {
	query *dao.Query
}

func (repo *EmailRedemptionTaskSubRepoV1) Read(ctx context.Context, CID int64, TaskSubID int64) (*model.EmlRedemptionTaskSub, error) {
	return repo.query.EmlRedemptionTaskSub.WithContext(ctx).Where(repo.query.EmlRedemptionTaskSub.Cid.Eq(CID), repo.query.EmlRedemptionTaskSub.TaskSubID.Eq(TaskSubID)).First()
}

func NewEmailRedemptionTaskSubV1(query *dao.Query) EmailRedemptionTaskSubRepo {
	return &EmailRedemptionTaskSubRepoV1{
		query: query,
	}
}
