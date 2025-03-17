package ginuserlikeitem

import (
	"first-proj/common"
	"first-proj/module/userlikeitem/biz"
	"first-proj/module/userlikeitem/model"
	"first-proj/module/userlikeitem/storage"
	"net/http"
	"time"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
			id, err := common.FromBase58(c.Param("id"))

			if err != nil {
				panic(common.ErrInvalidRequest(err))
			}
		

			requester := c.MustGet(common.CurrentUser).(common.Requester)
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			store := storage.NewSQLStore(db)
			business := biz.NewUserLikeItemBiz(store)
			now := time.Now().UTC()

			if err := business.LikeItem(c.Request.Context(), &model.Like{
				UserId: requester.GetUserId(), 
				ItemId: int(id.GetLocalID()),
				CreatedAt: &now,
			}); err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
