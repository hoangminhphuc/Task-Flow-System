package model

import (
	"errors"
	"first-proj/common"
	"strings"
)

var (
	ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
	ErrItemIsDeleted = errors.New("item is deleted")
)

const (
	EntityName = "Item"
)


type TodoItem struct {
	//embedded struct
	common.SQLModel
	UserId			int 								`json:"-" gorm:"column:user_id;"`
	Title       string     					`json:"title" gorm:"column:title;"`
	Description string     					`json:"description" gorm:"column:description;"`
	Status      string     					`json:"status" gorm:"-"`
	Image 			*common.Image 			`json:"image" gorm:"column:image;"`
	LikedCount  int 								`json:"liked_count" gorm:"column:liked_count;"`
	Owner 			*common.SimpleUser	`json:"owner" gorm:"foreignKey:UserId;"`
}

func (TodoItem) TableName() string { return "todo_items" }

func (item *TodoItem) Mask() {
	item.SQLModel.Mask(common.DbTypeItem)

	if value := item.Owner; value != nil {
		value.Mask()
	}
}

type TodoItemCreation struct {
	Id          int        			`json:"id" gorm:"column:id;"`
	UserId			int 						`json:"-" gorm:"column:user_id;"`
	Title       string     			`json:"title" gorm:"column:title;"`
	Description string     			`json:"description" gorm:"column:description;"`
	Image 			*common.Image 	`json:"image" gorm:"column:image;"`
}


// Check if title empty
func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title) 
	
	if i.Title == "" {
		return ErrTitleCannotBeEmpty
	}

	return nil
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

type TodoItemUpdate struct {
	Title       *string     `json:"title" gorm:"column:title;"`
	Description *string     `json:"description" gorm:"column:description;"`
	Status      *string     `json:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }