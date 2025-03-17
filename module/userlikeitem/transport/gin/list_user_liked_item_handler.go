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

func ListUserLikedItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
			id, err := common.FromBase58(c.Param("id"))

			if err != nil {
				panic(common.ErrInvalidRequest(err))
			}

			var queryString struct {
				common.Paging
			}
		
			if err := c.ShouldBind(&queryString); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
							"error": err.Error(),
					})
					return
			}

			queryString.Process()
		
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			store := storage.NewSQLStore(db)
			business := biz.NewListUserLikedItemBiz(store)

			result, err := business.ListUserLikedItem(c.Request.Context(), 
				int(id.GetLocalID()), 
				&queryString.Paging)

			if err != nil {
				panic(err)
			}

			for i := range result {
				result[i].Mask()
			}


			c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, nil))
	}
}
