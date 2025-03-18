package storage

import (
	"context"
	"first-proj/common"
	"first-proj/module/userlikeitem/model"
	"time"

	"github.com/btcsuite/btcutil/base58"
)

const timeLayout = "2006-01-02T15:04:05.999999"

func (s *sqlStore) ListUsers(
	ctx context.Context,
	itemId int,
	paging *common.Paging,
) ([]common.SimpleUser, error) {
	var result []model.Like

	db := s.db.Where("item_id = ?", itemId)

	if err := db.Table(model.Like{}.TableName()).Select("user_id").Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if v := paging.FakeCursor; v != "" {
		// Converts a string into a time.Time object
    timeCreated, err := time.Parse(timeLayout, string(base58.Decode(v)))

    if err != nil {
        return nil, common.ErrDB(err)
    }
		// Converts a time.Time object into a formatted string
    db = db.Where("created_at < ?", timeCreated.Format("2006-01-02 15:04:05.999999"))
	} else {
			db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Select("*").
		Preload("User").	
    Order("created_at desc").
    Limit(paging.Limit).
    Find(&result).Error; err != nil {
			return nil, common.ErrDB(err)
	}

	users := make([]common.SimpleUser, len(result))

	for i := range result {
		users[i] = *result[i].User
		users[i].UpdatedAt = nil
		users[i].CreatedAt = result[i].CreatedAt
	}

	//Set next cursor
	if len(users) > 0 {
    users[len(result)-1].Mask()
		/* Takes the last person on the page, extract the CreatedAt field,
		format it to a time.Time object */
    paging.NextCursor = base58.Encode([]byte(users[len(result)-1].
													CreatedAt.Format(timeLayout)))

		
}

	return users, nil

}

func (s *sqlStore) GetItemLikes(ctx context.Context, ids []int) (map[int]int, error) {
	result := make(map[int]int)

	type sqlData struct {
			ItemId int `gorm:"column:item_id"`
			Count  int `gorm:"column:count"`
	}

	var listLike []sqlData

	if err := s.db.Table(model.Like{}.TableName()).
			Select("item_id, COUNT(item_id) as `count`").
			Where("item_id in (?)", ids).
			Group("item_id").
			Find(&listLike).Error; err != nil {
			return nil, common.ErrDB(err)
	}

	for _, item := range listLike {
			result[item.ItemId] = item.Count
	}

	return result, nil
}
