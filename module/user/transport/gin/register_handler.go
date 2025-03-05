package ginuser

import (
	"first-proj/common"
	"first-proj/module/user/biz"
	"first-proj/module/user/model"
	"first-proj/module/user/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
			var data model.UserCreate


			if err := c.ShouldBind(&data); err != nil {
					panic(err)
			}


			store := storage.NewSQLStore(db)
			bcrypt := common.NewBcryptHash()
			biz := biz.NewRegisterBusiness(store, bcrypt)

			if err := biz.Register(c.Request.Context(), &data); err != nil {
					panic(err)
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.ID))
	}
}
