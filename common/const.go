package common

import "fmt"

const (
	CurrentUser = "current_user"
)

type DbType int 

const (
	DbTypeItem DbType = 1
	DbTypeUser DbType = 2
)

const (
	PluginDBMain = "mysql"
	PluginPubSub = "pubsub"
	PluginItemAPI = "item-api"
	PluginRedis = "redis"

	TopicUserLikedItem = "TopicUserLikedItem"
	TopicUserUnlikedItem = "TopicUserUnlikedItem"
)

/* 
! Security const
*/

const (
	BCRYPT_COST = 12
)


func Recovery() {
	if r := recover(); r != nil {
		fmt.Println("Recovered: ", r)
	}
}

type TokenPayLoad struct {
	UId int	`json:"user_id"`
	URole string `json:"role"`
}

func (p TokenPayLoad) UserId() int {	
	return p.UId
}

func (p TokenPayLoad) Role() string {
	return p.URole
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}