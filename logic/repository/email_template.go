package repository

import (
	"context"
	"email/common/log"
	"email/dal/cache"
	"email/dal/dao"
	"email/dal/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type EmailTemplateRepoV1 struct {
	Query *dao.Query
	Cache cache.Cache
}

func (repo *EmailTemplateRepoV1) ReadPriority(ctx context.Context, templateID int64) (int32, error) {
	key := fmt.Sprintf("%dpriority", templateID)
	val, _ := repo.Cache.Get(key)
	priority, ok := val.(int32)
	if ok {
		return priority, nil
	}
	qTemplate := repo.Query.EmlTemplate
	entity, err := qTemplate.WithContext(ctx).Where(repo.Query.EmlTemplate.ID.Eq(templateID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("template " + err.Error())
		}
		return 0, errors.New("template " + err.Error())
	}
	priority = entity.Priority
	repo.Cache.Set(key, priority, 0)
	return priority, nil
}

func (repo *EmailTemplateRepoV1) ReadMaxRetry(ctx context.Context, templateID int64) (int32, error) {
	key := fmt.Sprintf("%dmaxretry", templateID)
	val, _ := repo.Cache.Get(key)
	maxRetry, ok := val.(int32)
	if ok {
		return maxRetry, nil
	}
	qTemplate := repo.Query.EmlTemplate
	entity, err := qTemplate.WithContext(ctx).Where(repo.Query.EmlTemplate.ID.Eq(templateID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("template " + err.Error())
		}
		return 0, errors.New("template " + err.Error())
	}
	maxRetry = entity.MaxRetry
	repo.Cache.Set(key, maxRetry, 0)
	return maxRetry, nil
}

func (repo *EmailTemplateRepoV1) Read(ctx context.Context, cid int64, templateID int64) (*model.EmlTemplate, error) {
	key := fmt.Sprintf("%dtemplate%d", cid, templateID)
	val, err := repo.Cache.Get(key)
	if err != nil {
		return nil, errors.New("template " + err.Error())

	}
	template, ok := val.(*model.EmlTemplate)
	if ok {
		log.New(ctx).Debug("hit cache in EmailTemplateRepo GetEmailTask")
		return template, nil
	}
	log.New(ctx).Debug("miss cache in EmailTemplateRepo GetEmailTask")

	entity, err := repo.Query.WithContext(ctx).EmlTemplate.Where(repo.Query.EmlTemplate.Cid.Eq(cid), repo.Query.EmlTemplate.ID.Eq(templateID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			return nil, errors.New("template " + err.Error())
		}
	}
	repo.Cache.Set(key, entity, 0)
	return entity, nil
}

func NewEmailTemplateRepoV1(q *dao.Query, cache cache.Cache) EmailTemplateRepo {
	return &EmailTemplateRepoV1{
		Query: q,
		Cache: cache,
	}
}
