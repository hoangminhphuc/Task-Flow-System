package ginitem

import (
	"first-proj/common"
	"first-proj/module/item/biz"
	"first-proj/module/item/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		//id in the URL example: localhost:8080/v1/items/:id
		//this id is in string format => convert to int
		id, err := strconv.Atoi(c.Param("id")) 

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		//asign db to store
		store := storage.NewSQLStore(db)

		
		business := biz.NewGetItemBiz(store)

		data, err :=  business.GetItemById(c.Request.Context(), id)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)

			return
			// panic(err) => only for practice, not best practice
		}

		data.Mask()

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))

	}
}