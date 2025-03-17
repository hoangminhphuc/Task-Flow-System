package biz

import (
	"context"
	"first-proj/common"
	"first-proj/plugin/tokenprovider"
	"first-proj/module/user/model"
)

type LoginStorage interface {
    FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

type loginBusiness struct {
    storeUser     LoginStorage
    tokenProvider tokenprovider.Provider
    hasher        Hasher
    expiry        int
}

func NewLoginBusiness(storeUser LoginStorage, tokenProvider tokenprovider.Provider, hasher Hasher, expiry int) *loginBusiness {
    return &loginBusiness{
        storeUser:     storeUser,
        tokenProvider: tokenProvider,
        hasher:        hasher,
        expiry:        expiry,
    }
}

// 1. Find user, email
// 2. Hash pass from input and compare with pass in db
// 3. Provider: issue JWT Token for client
// 3.1. access token and refresh token
// 4. return token(s)

func (business *loginBusiness) Login(ctx context.Context, data *model.UserLogin) (tokenprovider.Token, error) {
	user, err := business.storeUser.FindUser(ctx, map[string]interface{}{"email": data.Email})
	if err != nil {
			return nil, model.ErrEmailOrPasswordInvalid
	}

	//the hash always follows this structure: password + salt
	if !business.hasher.Compare(user.Password, data.Password + user.Salt) {
			return nil, model.ErrEmailOrPasswordInvalid
	}

	payload := &common.TokenPayLoad{
			UId:  user.ID,
			URole: user.Role.String(),
	}

	accessToken, err := business.tokenProvider.Generate(payload, business.expiry)
	if err != nil {
			return nil, common.ErrInternal(err)
	}

	//should use refresh token along with access token

	// refreshToken, err := business.tokenProvider.Generate(payload, business.tkCfg.GetRtExp())
	// if err != nil {
	//     return nil, common.ErrInternal(err)
	// }

	// account := usermodel.NewAccount(accessToken, refreshToken)


	return accessToken, nil
}

