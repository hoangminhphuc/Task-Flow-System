package ginuserlikeitem

import (
	"first-proj/common"
	// itemStorage "first-proj/module/item/storage"
	"first-proj/module/userlikeitem/biz"
	"first-proj/module/userlikeitem/storage"
	"first-proj/pubsub"
	"net/http"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UnlikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
			id, err := common.FromBase58(c.Param("id"))

			if err != nil {
				panic(common.ErrInvalidRequest(err))
			}
		

			requester := c.MustGet(common.CurrentUser).(common.Requester)
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
			ps := serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)

			store := storage.NewSQLStore(db)
			// itemStore := itemStorage.NewSQLStore(db)
			business := biz.NewUserUnlikeItemBiz(store, ps)

			if err := business.UnlikeItem(c.Request.Context(), 
			requester.GetUserId(), int(id.GetLocalID())); err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
