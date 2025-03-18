package biz

import (
	"context"
	"first-proj/common"
	"first-proj/module/userlikeitem/model"
	"log"
)

type UserUnlikeItemStore interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type DecreaseItemStorage interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userUnlikeItemBiz struct {
	store 		UserUnlikeItemStore
	itemStore DecreaseItemStorage
}

func NewUserUnlikeItemBiz(store UserUnlikeItemStore, itemStore DecreaseItemStorage) *userUnlikeItemBiz {
	return &userUnlikeItemBiz{store: store, itemStore: itemStore}
}

func (biz *userUnlikeItemBiz) UnlikeItem(ctx context.Context, userId, itemId int) error {
	_, err := biz.store.Find(ctx, userId, itemId)

	if err == common.RecordNotFound {
		return model.ErrDidNotLikeItem(err)
	}

	if err != nil {
		return model.ErrCannotUnlikeItem(err)
	}


	
	if err := biz.store.Delete(ctx, userId, itemId); err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	//Nghiệp vụ phụ, chạy được hay không không quan tâm
	go func ()  {
		defer common.Recovery()
		
		if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
			log.Println(err)
		}
	}()

	return nil
}


