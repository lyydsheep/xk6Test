package repository

import (
	"context"
	"email/common/log"
	"email/dal/cache"
	"email/dal/dao"
	"email/dal/model"
	"errors"
	"fmt"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type EmailTaskConfigRepoV1 struct {
	Query *dao.Query
	Cache cache.Cache
}

func (repo *EmailTaskConfigRepoV1) ReadAll(ctx context.Context) ([]*model.EmlTaskConfig, error) {
	return repo.Query.EmlTaskConfig.WithContext(ctx).Find()
}

func (repo *EmailTaskConfigRepoV1) Create(ctx context.Context, entity *model.EmlTaskConfig) (*model.EmlTaskConfig, error) {
	return entity, repo.Query.EmlTaskConfig.WithContext(ctx).Create(entity)
}

func (repo *EmailTaskConfigRepoV1) Read(ctx context.Context, cid int64, category string, _type string) (*model.EmlTaskConfig, error) {
	key := fmt.Sprintf("%dConfig%s%s", cid, category, _type)
	val, err := repo.Cache.Get(key)
	if err != nil {
		return nil, fmt.Errorf("查询模板配置失败: %w", err)

	}
	taskConfig, ok := val.(*model.EmlTaskConfig)
	if ok {
		log.New(ctx).Debug("hit cache in TaskConfigRepo")
		return taskConfig, nil
	}
	log.New(ctx).Debug("miss cache in TaskConfigRepo")

	entity, err := repo.Query.EmlTaskConfig.WithContext(ctx).Where(repo.Query.EmlTaskConfig.Cid.Eq(cid), repo.Query.EmlTaskConfig.Category.Eq(category),
		repo.Query.EmlTaskConfig.Type.Eq(_type)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("template not found")
		} else {
			return nil, fmt.Errorf("failed to read template: %w", err)
		}
	}
	repo.Cache.Set(key, entity, 0)
	return entity, nil
}

func (repo *EmailTaskConfigRepoV1) Update(ctx context.Context, entity *model.EmlTaskConfig) (gen.ResultInfo, error) {
	info, err := repo.Query.EmlTaskConfig.WithContext(ctx).Where(repo.Query.EmlTaskConfig.ID.Eq(entity.ID)).Updates(entity)
	if err != nil {
		return info, fmt.Errorf("更新模板配置失败: %w", err)
	}
	key := fmt.Sprintf("%dConfig%s%s", entity.Cid, entity.Category, entity.Type)
	repo.Cache.Set(key, entity, 0)
	return info, err
}

func NewTaskConfigRepoV1(q *dao.Query, cache cache.Cache) TaskConfigRepo {
	return &EmailTaskConfigRepoV1{
		Query: q,
		Cache: cache,
	}
}
