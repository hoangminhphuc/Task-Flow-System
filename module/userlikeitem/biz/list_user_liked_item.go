package biz

import (
	"context"
	"first-proj/common"
	"first-proj/module/userlikeitem/model"
)

type ListUserLikedItemStore interface {
	ListUsers(ctx context.Context, itemId int, paging *common.Paging) ([]common.SimpleUser, error)
}

type listUserLikedItemBiz struct {
	store ListUserLikedItemStore
}

func NewListUserLikedItemBiz(store ListUserLikedItemStore) *listUserLikedItemBiz {
	return &listUserLikedItemBiz{store: store}
}

func (biz *listUserLikedItemBiz) ListUserLikedItem(
	ctx context.Context, 
	itemId int, 
	paging *common.Paging) ([]common.SimpleUser, error) {
	result, err := biz.store.ListUsers(ctx, itemId, paging)

	if err != nil {
		return nil, common.ErrCannatListEntity(model.EntityName, err)
	}

	return result, nil
}


