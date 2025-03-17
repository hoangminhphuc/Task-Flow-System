package ginuser

import (
	"first-proj/common"
	"first-proj/plugin/tokenprovider"
	"first-proj/module/user/biz"
	"first-proj/module/user/model"
	"first-proj/module/user/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(db *gorm.DB, tokenProvider tokenprovider.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserData model.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}


		store := storage.NewSQLStore(db)
		bcrypt := common.NewBcryptHash()

		// should be expired in 7 days, in this case we set it to 30 days
		business := biz.NewLoginBusiness(store, tokenProvider, bcrypt, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			// Handle error
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
