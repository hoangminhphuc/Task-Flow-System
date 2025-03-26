package ginuserlikeitem

import (
	"first-proj/common"
	"first-proj/module/userlikeitem/biz"
	"first-proj/module/userlikeitem/storage"
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetItemLikes(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
			type RequestData struct {
				Ids []int `json:"ids"`
			}

			var data RequestData
		
			if err := c.ShouldBind(&data); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
					})
					return
			}
		
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			store := storage.NewSQLStore(db)
			business := biz.NewGetItemLikesBiz(store)

			result, err := business.GetItemLikes(c.Request.Context(), data.Ids)

			if err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
