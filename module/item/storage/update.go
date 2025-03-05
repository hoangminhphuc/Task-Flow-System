package storage

import (
	"context"
	"first-proj/common"
	"first-proj/module/item/model"
)

func (s *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	if err := s.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}