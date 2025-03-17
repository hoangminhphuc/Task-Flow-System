package biz

import (
	"context"
	"first-proj/module/userlikeitem/model"
)

type UserLikeItemStore interface {
	Create(ctx context.Context, data *model.Like) error
}

type userLikeItemBiz struct {
	store UserLikeItemStore
}

func NewUserLikeItemBiz(store UserLikeItemStore) *userLikeItemBiz {
	return &userLikeItemBiz{store: store}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	return nil
}


