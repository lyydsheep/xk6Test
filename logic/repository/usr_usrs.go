package repository

import (
	"context"
	"email/dal/dao"
	"email/dal/model"
	"errors"
)

type UsrUserRepoV1 struct {
	Query *dao.Query
}

func (u *UsrUserRepoV1) Read(ctx context.Context, cid int64, userID string) (*model.UsrUser, error) {
	// 查询用户信息
	user, err := u.Query.UsrUser.WithContext(ctx).Where(u.Query.UsrUser.Cid.Eq(cid), u.Query.UsrUser.UserID.Eq(userID)).First()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsrUserRepoV1) ReadTagUser(ctx context.Context, cid int64, tag string) ([]*model.UsrUser, error) {
	entities, err := u.Query.UsrUser.WithContext(ctx).Where(u.Query.UsrUser.Cid.Eq(cid), u.Query.UsrUser.Tags.Like("%"+tag+"%")).Find()
	if err != nil {
		return nil, errors.New("usrUser " + err.Error())
	}
	return entities, nil
}

func NewUsrUserRepoV1(query *dao.Query) UsrUserRepo {
	return &UsrUserRepoV1{
		Query: query,
	}
}
