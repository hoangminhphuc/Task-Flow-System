package ginitem

import (
	"first-proj/common"
	"first-proj/module/item/biz"
	"first-proj/module/item/model"
	"first-proj/module/item/storage"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var itemData model.TodoItemCreation

		//parse request data (JSON in this case) and bind it to a struct.
		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		itemData.UserId = requester.GetUserId()


		store := storage.NewSQLStore(db)
		business := biz.NewCreateItemBiz(store)

		if err := business.CreateNewItem(c.Request.Context(), &itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))

	}
}