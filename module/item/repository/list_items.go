package repository

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

type ItemLikeStorage interface {
	GetItemLikes(ctx context.Context, ids []int) (map[int]int, error)
}

type listItemRepo struct {
	store 		ListItemStorage
	likeStore ItemLikeStorage
	requester common.Requester
}

func NewListItemRepo(store ListItemStorage, likeStore ItemLikeStorage, requester common.Requester) *listItemRepo {
	return &listItemRepo{store: store, likeStore: likeStore, requester: requester}
}

func (repo *listItemRepo) ListItem(ctx context.Context,
	filter *model.Filter, paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	//more keys, you can wrap it again and again
	ctxStore := context.WithValue(ctx, common.CurrentUser, repo.requester)

	data, err := repo.store.ListItem(ctxStore, filter, paging, moreKeys...)

	if err != nil {
		return nil, common.ErrCannatListEntity(model.EntityName, err)
	}

	if len(data) == 0 {
    return data, nil
	}

	ids := make([]int, len(data))

	for i := range ids {
			ids[i] = data[i].ID
	}

	likeUserMap, err := repo.likeStore.GetItemLikes(ctxStore, ids)

	if err != nil {
	}

	//Faster than append, especially on large scale project
	for i := range data {
			data[i].LikedCount = likeUserMap[data[i].ID]
	}


	return data, nil
}