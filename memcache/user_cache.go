package memcache

import (
	"context"
	"first-proj/module/user/model"
	"fmt"
	"log"
	"sync"
	"time"
)

type RealStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type userCaching struct {
	store     Cache
	realStore RealStore
	once      *sync.Once
}

func NewUserCaching(store Cache, realStore RealStore) *userCaching {
	return &userCaching{
		store:     store,
		realStore: realStore,
		once:      new(sync.Once),
	}
}

func (uc *userCaching) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
	var user model.User
	
	userId := conditions["id"].(int)
	key := fmt.Sprintf("user-%d", userId)
	err := uc.store.Get(ctx, key, &user)

	// Found in cache, returns immediately
	if err != nil && user.ID > 0 {
			return &user, nil
	}

	var userErr error

	uc.once.Do(func() {
			realUser, userErr := uc.realStore.FindUser(ctx, conditions, moreInfo...)

			if userErr != nil {
					log.Println(userErr)
					return
			} 
			// Update cache
			user = *realUser
			_ = uc.store.Set(ctx, key, realUser, time.Hour * 2)
	})

	if userErr != nil {
			return nil, userErr
	}

	// Ensures all goroutines get the latest cached result instead of making DB calls.
	err = uc.store.Get(ctx, key, &user)

	if err != nil && user.ID > 0 {
			return &user, nil
	}
	
	// Even after second attempt retrieving from cache, it still returns error
	// Then return nil
	return nil, err
}