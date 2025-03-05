package biz

import (
	"context"
	"first-proj/module/item/model"
	"first-proj/common"
)

type ListItemStorage interface {
	ListItem(ctx context.Context,
		filter *model.Filter, paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}

type listItemBiz struct {
	store ListItemStorage
	requester common.Requester
}

func NewListItemBiz(store ListItemStorage, requester common.Requester) *listItemBiz {
	return &listItemBiz{store: store, requester: requester}
}

func (biz *listItemBiz) ListItem(ctx context.Context,
	filter *model.Filter, paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	//more keys, you can wrap it again and again
	ctxStore := context.WithValue(ctx, common.CurrentUser, biz.requester)

	data, err := biz.store.ListItem(ctxStore, filter, paging, "Owner")

	if err != nil {
		return nil, common.ErrCannatListEntity(model.EntityName, err)
	}

	return data, nil
}