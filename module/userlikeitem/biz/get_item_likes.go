package biz

import (
	"context"
	"first-proj/common"
	"first-proj/module/item/model"
)

type GetItemLikesStore interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type getItemLikesBiz struct {
	store GetItemLikesStore
}

func NewGetItemLikesBiz(store GetItemLikesStore) *getItemLikesBiz {
	return &getItemLikesBiz{store: store}
}

func (biz *getItemLikesBiz) GetItemLikes(ctx context.Context, 
	ids []int) (map[int]int, error) {
		result, err := biz.store.GetItemLikes(ctx, ids)

		if err != nil {
			return nil, common.ErrCannotGetEntity(model.EntityName, err)
		}

		return result, nil
}


