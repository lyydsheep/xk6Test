package repository

import (
	"context"
	"email/dal/dao"
	"email/dal/model"
	"gorm.io/gen"
	"time"
)

type EmailDomainCredibilityRepoV1 struct {
	Query *dao.Query
}

func (repo *EmailDomainCredibilityRepoV1) UpdateSpeed(ctx context.Context, entity *model.EmlDomainCredibility) (gen.ResultInfo, error) {
	return repo.Query.EmlDomainCredibility.WithContext(ctx).Where(repo.Query.EmlDomainCredibility.ID.Eq(entity.ID)).Updates(map[string]interface{}{
		"speed":       entity.Speed,
		"update_time": time.Now(),
	})
}

func (repo *EmailDomainCredibilityRepoV1) ReadAll(ctx context.Context) ([]*model.EmlDomainCredibility, error) {
	return repo.Query.EmlDomainCredibility.WithContext(ctx).Find()
}

func (repo *EmailDomainCredibilityRepoV1) Create(ctx context.Context, entity *model.EmlDomainCredibility) (*model.EmlDomainCredibility, error) {
	err := repo.Query.EmlDomainCredibility.WithContext(ctx).Create(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *EmailDomainCredibilityRepoV1) Read(ctx context.Context, cid int64, domain string, fromAddress string) (*model.EmlDomainCredibility, error) {
	return repo.Query.EmlDomainCredibility.WithContext(ctx).Where(repo.Query.EmlDomainCredibility.Cid.Eq(cid),
		repo.Query.EmlDomainCredibility.Domain.Eq(domain),
		repo.Query.EmlDomainCredibility.FromAddress.Eq(fromAddress)).First()
}

func (repo *EmailDomainCredibilityRepoV1) UpdateLastSentTime(ctx context.Context, entity *model.EmlDomainCredibility) (gen.ResultInfo, error) {
	return repo.Query.EmlDomainCredibility.WithContext(ctx).Where(repo.Query.EmlDomainCredibility.ID.Eq(entity.ID)).Updates(map[string]interface{}{
		"last_sent_time": entity.LastSentTime,
		"update_time":    time.Now(),
	})
}

func NewEmailDomainCredibilityRepoV1(q *dao.Query) EmailDomainCredibilityRepo {
	return &EmailDomainCredibilityRepoV1{
		Query: q,
	}
}
