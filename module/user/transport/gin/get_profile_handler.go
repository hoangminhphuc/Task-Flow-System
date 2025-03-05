package ginuser

import (
	"first-proj/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Profile() gin.HandlerFunc {
	return func(c *gin.Context) {
			u := c.MustGet(common.CurrentUser)

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
