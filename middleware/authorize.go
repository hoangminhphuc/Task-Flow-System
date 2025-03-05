package middleware

import (
	"context"
	"errors"
	"first-proj/common"
	"first-proj/component/tokenprovider"
	"first-proj/module/user/model"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthenStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error)
}

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("wrong authen header"),
		fmt.Sprintf("ErrWrongAuthHeader"),
	)
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	// "Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", ErrWrongAuthHeader(nil)
	}

	return parts[1], nil
}


// 1. Get token from header
// 2. Validate token and parse to payload
// 3. From the token payload, we use user_id to find from DB
func RequiredAuth(authStore AuthenStore, tokenProvider tokenprovider.Provider) func(c *gin.Context) {
	return func(c *gin.Context) {
			token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

			if err != nil {
					panic(err)
			}

			//db := appCtx.GetMaiDBConnection()
			//store := userstore.NewSQLStore(db)
			
			payload, err := tokenProvider.Validate(token)
			if err != nil {
					panic(err)
			}


			//This part is a cache-able part, it keeps going to db to find user
			user, err := authStore.FindUser(c.Request.Context(), map[string]interface{}{"id": payload.UserId()})

			if err != nil {
					panic(err)
			}

			if user.Status == 0 {
					panic(common.ErrNoPermission(errors.New("user has been deleted or banned")))
			}

			c.Set(common.CurrentUser, user)
			c.Next()
	}
}
