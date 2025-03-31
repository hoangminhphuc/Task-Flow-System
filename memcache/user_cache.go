package memcache

import (
	"context"
	"first-proj/module/user/model"
	"fmt"
	"sync"
)

type RealStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type userCaching struct {
	store     Caching
	realStore RealStore
	once      *sync.Once
}

func NewUserCaching(store Caching, realStore RealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

func (uc *userCaching) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
	userId := conditions["id"].(int)
	key := fmt.Sprintf("user-%d", userId)
	userInCache := uc.store.Read(key)

	if userInCache != nil {
			return userInCache.(*model.User), nil
	}

	uc.once.Do(func() {
			user, err := uc.realStore.FindUser(ctx, conditions, moreInfo...)
			if err != nil {
					panic(err)
			}
			// Update cache
			uc.store.Write(key, user)
	})
	
	return uc.store.Read(key).(*model.User), nil
}