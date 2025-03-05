package biz

import (
	"context"
	"first-proj/common"
	"first-proj/module/item/model"
)

/** 
	 *! Business layer
**/


type GetItemStorage interface {
	//use map here to support multiple conditions passing in (not just id)
	GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error)
}

type getItemBiz struct {
	store GetItemStorage
}

func NewGetItemBiz(store GetItemStorage) *getItemBiz {
	return &getItemBiz{store : store}
}



/* 
	* A function of getItemBiz struct with business logic

	Given an item ID, return the item if it exists, 
	otherwise return a business-specific error.
*/


//A function of getItemBiz struct with business logic
func (biz *getItemBiz) GetItemById(ctx context.Context, id int) (*model.TodoItem, error) {
	data, err := biz.store.GetItem(ctx, map[string]interface{}{"id": id})

	if err != nil {
		return nil, common.ErrCannotGetEntity(model.EntityName, err)
	}

	return data, nil

}