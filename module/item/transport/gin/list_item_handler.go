package ginitem

import (
	"first-proj/common"
	"first-proj/module/item/biz"
	"first-proj/module/item/model"
	"first-proj/module/item/repository"
	"first-proj/module/item/storage"
	"first-proj/module/item/storage/restapi"
	"net/http"


	goservice "github.com/200Lab-Education/go-sdk"
	// userLikeStore "first-proj/module/userlikeitem/storage"


	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListItems(serviceCtx goservice.ServiceContext) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		apiItemCaller := serviceCtx.MustGet(common.PluginItemAPI).(interface {
			GetServiceURL() string
		})

		var queryString struct {
			//embedded struct
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		queryString.Paging.Process()

		requester := c.MustGet(common.CurrentUser).(common.Requester)


		store := storage.NewSQLStore(db)
		likeStore := restapi.New(apiItemCaller.GetServiceURL(), serviceCtx.Logger("restapi.item"))
		repo := repository.NewListItemRepo(store, likeStore, requester)
		business := biz.NewListItemBiz(repo, requester)


		items, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		/* 
		! Transport layer is responsible for encode, decode id
		*/

		for i := range items {
			items[i].Mask()
		}

		//formatting the response JSON
		c.JSON(http.StatusOK, common.NewSuccessResponse(items, queryString.Paging, queryString.Filter))

	}
}