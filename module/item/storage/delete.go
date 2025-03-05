package storage

import (
	"context"
	"first-proj/common"
	"first-proj/module/item/model"
)

func (s *sqlStore) DeleteItem(ctx context.Context, cond map[string]interface{}) error {
	deleteStatus := "Deleted"

	if err := s.db.Table(model.TodoItem{}.TableName()).
	Where(cond).
	Updates(map[string]interface{}{
		"status": deleteStatus,
	}).Error; err != nil {
		
		return common.ErrDB(err)
	}

	return nil
}
