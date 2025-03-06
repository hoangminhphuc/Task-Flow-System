package biz

import (
	"context"
	"first-proj/common"
	"first-proj/module/user/model"
)

type GetProfileStorage interface {
	GetStats(ctx context.Context, userId int) (*model.UserProfileStats, error)
}

type getProfileBiz struct {
	store GetProfileStorage
}

func NewGetProfileBiz(store GetProfileStorage) *getProfileBiz {
	return &getProfileBiz{store: store}
}

func (biz *getProfileBiz) GetStats(ctx context.Context, userId int) (*model.UserProfileStats, error) {
	data, err :=  biz.store.GetStats(ctx, userId)

	if err != nil {
		return nil, common.ErrCannotGetEntity(model.EntityName, err)
	}

	return data, nil
}