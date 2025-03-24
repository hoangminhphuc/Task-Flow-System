package model

import (
	"first-proj/common"
	"fmt"
	"time"
)

const (
	EntityName = "UserLikeItem"
)

type Like struct {
	UserId    int        					`json:"user_id" gorm:"column:user_id"`
	ItemId    int        					`json:"item_id" gorm:"column:item_id"`
	CreatedAt *time.Time 					`json:"created_at" gorm:"column:created_at"`
	User 			*common.SimpleUser 	`json:"-" gorm:"foreignKey:UserId;"`
}

func (l *Like) GetItemID() int { return l.ItemId }

func (Like) TableName() string { return "user_like_items" }

func ErrCannotLikeItem(err error) *common.AppError {
    return common.NewCustomError(
        err,
        fmt.Sprintf("Cannot like this item"),
        fmt.Sprintf("ErrCannotLikeItem"),
    )
}

func ErrCannotUnlikeItem(err error) *common.AppError {
    return common.NewCustomError(
        err,
        fmt.Sprintf("Cannot unlike this item"),
        fmt.Sprintf("ErrCannotUnlikeItem"),
    )
}

func ErrDidNotLikeItem(err error) *common.AppError {
		return common.NewCustomError(
				err,
				fmt.Sprintf("You have not like this item"),
				fmt.Sprintf("ErrDidNotLikeItem"),
		)
}