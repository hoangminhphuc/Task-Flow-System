package cmd

import (
	// "context"
	"first-proj/common"
	"first-proj/middleware"
	ginitem "first-proj/module/item/transport/gin"
	"first-proj/module/upload"
	userstorage "first-proj/module/user/storage"
	ginuser "first-proj/module/user/transport/gin"
	ginuserlikeitem "first-proj/module/userlikeitem/transport/gin"
	"first-proj/plugin/sdkgorm"
	"first-proj/plugin/simple"
	"first-proj/plugin/tokenprovider/jwt"
	"first-proj/pubsub"
	"first-proj/subscriber"
	"fmt"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	// "net/http"
	"os"
)


func newService() goservice.Service {
    service := goservice.New(
        goservice.WithName("social-todo-list"),
        goservice.WithVersion("1.0.0"),
        goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.PluginDBMain)),
				goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
				goservice.WithInitRunnable(simple.NewSimplePlugin("simple")),
    )

    return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
			systemSecret := os.Getenv("SECRET")

			service := newService()

			serviceLogger := service.Logger("service")

			if err := service.Init(); err != nil {
					serviceLogger.Fatalln(err)
			}

			service.HTTPServer().AddHandler(func(engine *gin.Engine) {
				engine.Use(middleware.Recover())
				engine.Use(middleware.CORSConfig())

			/* 
			! Example of how to use plugin
			*/		
				//inline implement interface
				service.MustGet("simple").(interface{
					GetValue() string
				}).GetValue()



				db := service.MustGet(common.PluginDBMain).(*gorm.DB)

				authStore := userstorage.NewSQLStore(db)
				tokenProvider := jwt.NewTokenJWTProvider("jwt", systemSecret)
				middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)


				v1 := engine.Group("/v1")
				{
						v1.PUT("/upload", upload.Upload(db))

						v1.POST("/register", ginuser.Register(db))
						v1.POST("/login", ginuser.Login(db, tokenProvider))
						v1.GET("/profile", middlewareAuth, ginuser.Profile(db))

						items := v1.Group("/items", middlewareAuth)
						{
								items.POST("", ginitem.CreateItem(db))
								items.GET("", ginitem.ListItems(db) )
								items.GET("/:id", ginitem.GetItem(db))
								items.PATCH("/:id", ginitem.UpdateItem(db))
								items.DELETE("/:id", ginitem.DeleteItem(db))

								//RPC
								items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
								items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
								items.GET("/:id/likers", ginuserlikeitem.ListUserLikedItem(service))
						}
				}
			})

			_ = subscriber.NewEngine(service).Start()

			if err := service.Start(); err != nil {
					serviceLogger.Fatalln(err)
			}
	},
}


func Execute() {
	// TransAddPoint outenv as a sub command
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
	}
}
