package storage

import (
	"context"
	"first-proj/common"
	"first-proj/module/user/model"

	"go.opencensus.io/trace"
	"gorm.io/gorm"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
    _, span := trace.StartSpan(ctx, "user.storage.find")
    defer span.End()
    
    db := s.db.Table(model.User{}.TableName())

    for i := range moreInfo {
        db = db.Preload(moreInfo[i])
    }

    var user model.User

    if err := db.Where(conditions).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, common.RecordNotFound
        }
        return nil, common.ErrDB(err)
    }

    return &user, nil
}
