package biz

import (
	"context"
	"first-proj/common"
	"first-proj/module/userlikeitem/model"
	"log"
)

type UserLikeItemStore interface {
	Create(ctx context.Context, data *model.Like) error
}

type IncreaseItemStorage interface {
	IncreaseLikeCount(ctx context.Context, id int) error
}

type userLikeItemBiz struct {
	store 		UserLikeItemStore
	itemStore IncreaseItemStorage
}

func NewUserLikeItemBiz(store UserLikeItemStore, itemStore IncreaseItemStorage) *userLikeItemBiz {
	return &userLikeItemBiz{store: store, itemStore: itemStore}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	
	//Nghiệp vụ chính
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}


	//Nghiệp vụ phụ, chạy được hay không không quan tâm
		//Để không làm giảm hiệu suất của business like item này, ta cho vào go routine
	go func ()  {
		defer common.Recovery()

		if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
			log.Println(err)
		}
	}()


	return nil
}


