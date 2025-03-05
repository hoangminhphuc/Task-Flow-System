package storage

import (
	"first-proj/common"
	"first-proj/module/item/model"

	"context"

	"gorm.io/gorm"
)


/* 
	! Storage layer
*/


/* 
	* This implements GetItemStorage interface because this function 
	* is a method of struct sqlStore. 
	! Only struct can implement interface
 */


func (s *sqlStore) GetItem(ctx context.Context, cond map[string]interface{}) (*model.TodoItem, error) {
	var data model.TodoItem


	// First fetches the first matching row from db and stores it in data.
	if err := s.db.Where(cond).First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &data, nil
}