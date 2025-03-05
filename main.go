package main

import (
	// "encoding/json"
	// "fmt"
	"log"
	// "net/http"
	"os"

	// "strconv"

	// "time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	//user-defined packages
	// "first-proj/common"
	// "first-proj/module/item/model"
	"first-proj/component/tokenprovider/jwt"
	"first-proj/middleware"
	"first-proj/module/user/storage"
	ginitem "first-proj/module/item/transport/gin"
	"first-proj/module/upload"
	ginuser "first-proj/module/user/transport/gin"
)


func main() {
	godotenv.Load(".env")

	db_url := os.Getenv("DATABASE_URL")
	systemSecret := os.Getenv("SECRET")
	
	dsn := db_url
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db = db.Debug()

	if err != nil {
		//similar to print but exit program
		log.Fatalln(err)
	}
	
	log.Println("Connected to DB", db)

	authStore := storage.NewSQLStore(db)
	tokenProvider := jwt.NewTokenJWTProvider("jwt", systemSecret)
	middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

	
	r := gin.Default()
	r.Use(middleware.Recover())

	r.Static("/static", "./static")


	v1 := r.Group("/v1")
{
		v1.PUT("/upload", upload.Upload(db))

		v1.POST("/register", ginuser.Register(db))
		v1.POST("/login", ginuser.Login(db, tokenProvider))
		v1.GET("/profile", middlewareAuth, ginuser.Profile())

    items := v1.Group("/items", middlewareAuth)
    {
        items.POST("", ginitem.CreateItem(db))
        items.GET("", ginitem.ListItems(db) )
        items.GET("/:id", ginitem.GetItem(db))
        items.PATCH("/:id", ginitem.UpdateItem(db))
        items.DELETE("/:id", ginitem.DeleteItem(db))
    }
}


  // r.GET("/ping", func(c *gin.Context) {
  //   c.JSON(http.StatusOK, gin.H{
  //     "message": "pong",
  //   })
  // })
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
