package biz

import (
	"context"
	"errors"
	"first-proj/common"
	"first-proj/module/item/model"
)

type UpdateItemStorage interface {
	//use map here to support multiple conditions passing in (not just id)
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
	UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error
}

type updateItemBiz struct {
	store UpdateItemStorage
	requester common.Requester
}

func NewUpdateItemBiz(store UpdateItemStorage, requester common.Requester) *updateItemBiz {
	return &updateItemBiz{store : store, requester: requester}
}

func (biz *updateItemBiz) UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return common.ErrCannotGetEntity(model.EntityName, err)
	}

	if data.Status == "Deleted" {
		return model.ErrItemIsDeleted
	}


	//Check user role
	isOwner := biz.requester.GetUserId() == data.UserId
	if !isOwner && !common.IsAdmin(biz.requester){
		return common.ErrNoPermission(errors.New("no permission"))
	}


	if err := biz.store.UpdateItem(ctx, map[string]interface{}{"id": id}, dataUpdate); err != nil {
		return common.ErrCannotUpdateEntity(model.EntityName, err)
	}

	return nil

}