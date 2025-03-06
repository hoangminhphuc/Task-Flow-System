package ginuser

import (
	"first-proj/common"
	"first-proj/module/user/biz"
	"first-proj/module/user/model"
	"first-proj/module/user/storage"
	"log"
	"net/http"

	// "strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Profile(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {

		currentUser := c.MustGet(common.CurrentUser) 
		user, ok := currentUser.(*model.User)       

		if !ok {
				log.Println("Failed to cast")
				return
		}
		userId := user.GetUserId()

		store := storage.NewSQLStore(db)
		biz := biz.NewGetProfileBiz(store)


		userStats, err := biz.GetStats(c.Request.Context(), userId)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}

		userProfileResponse := model.ProfileResponse{
			User:  user,
			Stats: *userStats,
	}

		
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(userProfileResponse))
	}
}
