package repository

import (
	"context"
	"email/dal/dao"
	"email/dal/model"
	"gorm.io/gen"
	"time"
)

type EmailResultRepoV1 struct {
	Query *dao.Query
}

func (repo *EmailResultRepoV1) ReadBatchByRange(ctx context.Context, accountName string, startTime time.Time, endTime time.Time) ([]*model.OpsEmlResult, error) {
	qOps := repo.Query.OpsEmlResult
	return qOps.WithContext(ctx).Where(
		qOps.AccountName.Eq(accountName),
		qOps.SentTime.Gte(startTime),
		qOps.SentTime.Lt(endTime),
	).Find()
}

func (repo *EmailResultRepoV1) ReadByTo(ctx context.Context, toAddress string) ([]*model.OpsEmlResult, error) {
	return repo.Query.OpsEmlResult.WithContext(ctx).Where(
		repo.Query.OpsEmlResult.ToAddress.Eq(toAddress),
	).Find()
}

func (repo *EmailResultRepoV1) GetEndTime(ctx context.Context, fromAddress string) (time.Time, error) {
	var (
		result time.Time
	)
	err := repo.Query.OpsEmlResult.WithContext(ctx).
		Where(repo.Query.OpsEmlResult.AccountName.Eq(fromAddress)).
		Order(repo.Query.OpsEmlResult.SentTime.Desc()).
		Limit(1).
		Select(repo.Query.OpsEmlResult.SentTime).
		Scan(&result)
	return result, err
}

func (repo *EmailResultRepoV1) ReadByTime(ctx context.Context, fromAddress string, domain string, startTime time.Time, endTime time.Time) ([]*model.OpsEmlResult, error) {
	return repo.Query.OpsEmlResult.WithContext(ctx).Where(repo.Query.OpsEmlResult.AccountName.Eq(fromAddress),
		repo.Query.OpsEmlResult.Domain.Eq(domain), repo.Query.OpsEmlResult.SentTime.Between(startTime, endTime)).Find()
}

func (repo *EmailResultRepoV1) Create(ctx context.Context, entity *model.OpsEmlResult) error {
	return repo.Query.OpsEmlResult.WithContext(ctx).Create(entity)
}

func (repo *EmailResultRepoV1) Update(ctx context.Context, entity *model.OpsEmlResult) (gen.ResultInfo, error) {
	return repo.Query.OpsEmlResult.WithContext(ctx).Where(repo.Query.OpsEmlResult.ID.Eq(entity.ID)).Updates(entity)
}

func NewEmailResultRepoV1(q *dao.Query) EmailResultRepo {
	return &EmailResultRepoV1{
		Query: q,
	}
}
