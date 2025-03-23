package biz

import (
	"context"
	// "first-proj/common"
	"first-proj/common/asyncjob"
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


	//Nghiệp vụ phụ, chạy được hay ko ko quan tâm nên nó sẽ là job chạy ở background

	/* 
	! When job grows, this will also grow exponentially. To avoid this, we should
	! separate job to another place, and use pub/sub.
	*/
	job := asyncjob.NewJob(func(ctx context.Context) error {
		if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
			return err
		}

		return nil
	})

	if err := asyncjob.NewGroup(true, job).Run(ctx); err != nil {
		log.Println(err)
	}


	return nil
}


