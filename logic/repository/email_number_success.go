package repository

import (
	"context"
	"email/dal/dao"
	"email/dal/model"
	"time"
)

type OpsEmlNumberSuccessRepoV1 struct {
	Query *dao.Query
}

func (repo *OpsEmlNumberSuccessRepoV1) Read(ctx context.Context, cid int64, fromAddress string, startTime time.Time, endTime time.Time) ([]*model.OpsEmlNumberSuccess, error) {
	qNumber := repo.Query.OpsEmlNumberSuccess
	return qNumber.WithContext(ctx).Where(
		qNumber.Cid.Eq(cid),
		qNumber.FromAddress.Eq(fromAddress),
		qNumber.CreateTime.Between(startTime, endTime),
	).Find()
}

func (repo *OpsEmlNumberSuccessRepoV1) Create(ctx context.Context, entity *model.OpsEmlNumberSuccess) error {
	qNumber := repo.Query.OpsEmlNumberSuccess
	return qNumber.WithContext(ctx).Create(entity)
}

func NewOpsEmlNumberSuccessRepoV1(query *dao.Query) OpsEmlNumberSuccessRepo {
	return &OpsEmlNumberSuccessRepoV1{
		Query: query,
	}
}
