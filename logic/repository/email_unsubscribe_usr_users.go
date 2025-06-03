package repository

import (
	"context"
	"email/dal/dao"
	"email/dal/model"
	"gorm.io/gorm/clause"
)

type EmailUnsubscribeUsrUserRepoV1 struct {
	Query *dao.Query
}

func (repo *EmailUnsubscribeUsrUserRepoV1) Read(ctx context.Context, cid int64, uid string) (*model.EmlUnsubscribeUsrUser, error) {
	entity, err := repo.Query.EmlUnsubscribeUsrUser.WithContext(ctx).Where(repo.Query.EmlUnsubscribeUsrUser.Cid.Eq(cid),
		repo.Query.EmlUnsubscribeUsrUser.UID.Eq(uid)).First()
	return entity, err
}

func (repo *EmailUnsubscribeUsrUserRepoV1) Create(ctx context.Context, entity *model.EmlUnsubscribeUsrUser) (*model.EmlUnsubscribeUsrUser, error) {
	err := repo.Query.EmlUnsubscribeUsrUser.WithContext(ctx).Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func NewEmailUnsubscribeUsrUserRepoV1(query *dao.Query) EmailUnsubscribeUsrUserRepo {
	return &EmailUnsubscribeUsrUserRepoV1{
		Query: query,
	}
}
