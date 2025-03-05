package ginitem

import (
	"first-proj/common"
	"first-proj/module/item/biz"
	"first-proj/module/item/model"
	"first-proj/module/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateItem(db *gorm.DB) func(ctx *gin.Context) {
	return func(c *gin.Context) {

		//c.Param returns string so we need to convert to int
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var updatedData model.TodoItemUpdate

		if err := c.ShouldBind(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		store := storage.NewSQLStore(db)
		business := biz.NewUpdateItemBiz(store, requester)


		if err := business.UpdateItemById(c.Request.Context(), id, &updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))

	}
}