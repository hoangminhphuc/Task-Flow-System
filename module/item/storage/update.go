package storage

import (
	"context"
	"first-proj/common"
	"first-proj/module/item/model"

	"gorm.io/gorm"
)

func (s *sqlStore) UpdateItem(ctx context.Context, cond map[string]interface{}, dataUpdate *model.TodoItemUpdate) error {
	if err := s.db.Where(cond).Updates(dataUpdate).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) IncreaseLikeCount(ctx context.Context, id int) error {
	db := s.db

	if err := db.Table(model.TodoItem{}.TableName()).Where("id = ?", id).
			Update("liked_count", gorm.Expr("liked_count + ?", 1)).Error; err != nil {
			return common.ErrDB(err)
	}

	return nil
}

func (s *sqlStore) DecreaseLikeCount(ctx context.Context, id int) error {
	db := s.db

	if err := db.Table(model.TodoItem{}.TableName()).Where("id = ?", id).
			Update("liked_count", gorm.Expr("liked_count - ?", 1)).Error; err != nil {
			return common.ErrDB(err)
	}

	return nil
}
