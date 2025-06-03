package repository

import (
	"context"
	"email/common/log"
	"email/common/util"
	"email/config"
	"email/dal/cache"
	"email/dal/dao"
	"email/dal/model"
	"errors"
	"fmt"
	"time"
)

type EmailSysUserRepoV1 struct {
	Query *dao.Query
	Cache cache.Cache
}

func (repo *EmailSysUserRepoV1) ReadByType(ctx context.Context, cid int64, _type string) (*model.EmlSysUser, error) {
	qSysUser := repo.Query.EmlSysUser
	return qSysUser.WithContext(ctx).Where(qSysUser.Cid.Eq(cid), qSysUser.Type.Eq(_type)).First()
}

func (repo *EmailSysUserRepoV1) UpdateStartTime(ctx context.Context, id int64, startTime time.Time) error {
	_, err := repo.Query.EmlSysUser.WithContext(ctx).Where(repo.Query.EmlSysUser.ID.Eq(id)).Updates(map[string]interface{}{
		"fetch_start_time": startTime,
		"update_time":      time.Now(),
	})
	return err
}

func (repo *EmailSysUserRepoV1) ReadAll(ctx context.Context) ([]*model.EmlSysUser, error) {
	return repo.Query.EmlSysUser.WithContext(ctx).Find()
}

// Read implements EmailSysUserRepo.
func (repo *EmailSysUserRepoV1) Read(ctx context.Context, cid int64, fromEmail string) (*model.EmlSysUser, error) {
	key := fmt.Sprintf("%dsysUser%s", cid, fromEmail)
	val, err := repo.Cache.Get(key)
	if err != nil {
		return nil, errors.New("sysUser " + err.Error())
	}
	sysUser, ok := val.(*model.EmlSysUser)
	if ok {
		log.New(ctx).Debug("hit cache in EmailSysUserRepo GetEmailTask")
		return sysUser, nil
	}

	log.New(ctx).Debug("miss cache in EmailSysUserRepo GetEmailTask")

	// 找不到就报错
	emailSysUser, err := repo.Query.EmlSysUser.WithContext(ctx).Where(repo.Query.EmlSysUser.Cid.Eq(cid),
		repo.Query.EmlSysUser.FromAddress.Eq(fromEmail)).First()
	if err != nil {
		return nil, errors.New("sysUser " + err.Error())
	}
	emailSysUser.Password, err = util.Decrypt(emailSysUser.Password, config.App.AESKEY)
	if err != nil {
		return nil, errors.New("sysUser " + err.Error())
	}
	emailSysUser.AccessKeyID, err = util.Decrypt(emailSysUser.AccessKeyID, config.App.AESKEY)
	if err != nil {
		return nil, errors.New("sysUser " + err.Error())
	}
	emailSysUser.AccessKeySecret, err = util.Decrypt(emailSysUser.AccessKeySecret, config.App.AESKEY)
	if err != nil {
		return nil, errors.New("sysUser " + err.Error())
	}
	repo.Cache.Set(key, emailSysUser, 0)
	return emailSysUser, nil
}

func NewEmailSysUserRepoV1(query *dao.Query, cache cache.Cache) EmailSysUserRepo {
	return &EmailSysUserRepoV1{
		Query: query,
		Cache: cache,
	}
}
