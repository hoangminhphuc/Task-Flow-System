package biz

import (
	"context"
	"first-proj/common"
	"first-proj/module/user/model"
)

type RegisterStorage interface {
    FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
    CreateUser(ctx context.Context, data *model.UserCreate) error
}

type Hasher interface {
    Hash(data string) string
		//comparing before and after hash value
		Compare(hashedValue, plainText string) bool
}

type registerBusiness struct {
    registerStorage RegisterStorage
    hasher          Hasher
}

func NewRegisterBusiness(registerStorage RegisterStorage, hasher Hasher) *registerBusiness {
	return &registerBusiness{
			registerStorage: registerStorage,
			hasher:          hasher,
	}
}

func (business *registerBusiness) Register(ctx context.Context, data *model.UserCreate) error {
	user, _ := business.registerStorage.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if user != nil {
			// if user.Status == 0 {
			//     return error user has been disabled
			// }
			return model.ErrEmailExisted
	}

	salt := common.GenSalt(50)

	data.Password = business.hasher.Hash(data.Password + salt)
	data.Salt = salt
	data.Role = "user" // hard-coded role

	if err := business.registerStorage.CreateUser(ctx, data); err != nil {
			return common.ErrCannotCreateEntity(model.EntityName, err)
	}

	return nil
}
