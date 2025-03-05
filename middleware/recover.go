package middleware

import (
	"first-proj/common"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func (c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)

					return
				}
				
				appErr := common.ErrInternal(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				//We panic here to start stack trace, but not exit our program.
				//Stack trace will trace back to exactly this recover.go file, because this 
				//is where it start panic.
				// panic(err) 
				return
		}

	}()

	c.Next()
}
}