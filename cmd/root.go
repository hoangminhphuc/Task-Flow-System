package cmd

import (
	// "context"
	"first-proj/common"
	"first-proj/memcache"
	"first-proj/middleware"
	ginitem "first-proj/module/item/transport/gin"
	"first-proj/module/upload"
	userstorage "first-proj/module/user/storage"
	ginuser "first-proj/module/user/transport/gin"
	ginuserlikeitem "first-proj/module/userlikeitem/transport/gin"
	ginuserlikeitem_rpc "first-proj/module/userlikeitem/transport/rpc"
	"first-proj/plugin/appredis"
	"first-proj/plugin/rpccaller"
	"first-proj/plugin/sdkgorm"
	"first-proj/plugin/simple"
	"first-proj/plugin/tokenprovider/jwt"
	"first-proj/pubsub"
	"first-proj/subscriber"
	"fmt"
	"log"

	// "time"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
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
				goservice.WithInitRunnable(rpccaller.NewApiItemCaller(common.PluginItemAPI)),
				goservice.WithInitRunnable(appredis.NewRedisDB("redis", common.PluginRedis)),
				
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

		// init all the services that have been registered in the initServices map.
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
				// service.MustGet("simple").(interface{
				// 	GetValue() string
				// }).GetValue()


				db := service.MustGet(common.PluginDBMain).(*gorm.DB)

				authStore := userstorage.NewSQLStore(db)
				authCache := memcache.NewUserCaching(memcache.NewRedisCache(service), authStore)
				
				
				tokenProvider := jwt.NewTokenJWTProvider("jwt", systemSecret)
				middlewareAuth := middleware.RequiredAuth(authCache, tokenProvider)


				v1 := engine.Group("/v1")
				{
						v1.PUT("/upload", upload.Upload(db))

						v1.POST("/register", ginuser.Register(db))
						v1.POST("/login", ginuser.Login(db, tokenProvider))
						v1.GET("/profile", middlewareAuth, ginuser.Profile(db))

						items := v1.Group("/items", middlewareAuth)
						{
							items.POST("", ginitem.CreateItem(service))
							items.GET("", ginitem.ListItems(service))
							items.GET("/:id", ginitem.GetItem(service))
							items.PATCH("/:id", ginitem.UpdateItem(db))
							items.DELETE("/:id", ginitem.DeleteItem(db))

							//RPC
							items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
							items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
							items.GET("/:id/likers", ginuserlikeitem.ListUserLikedItem(service))
						}

						rpc := v1.Group("/rpc")
						{
							rpc.POST("/get_item_likes", ginuserlikeitem_rpc.GetItemLikes(service))
						}
				}
			})

			je, err := jaeger.NewExporter(jaeger.Options{
				AgentEndpoint: os.Getenv("JAEGER_AGENT_URL"),
				Process: jaeger.Process{ServiceName: "Task-Management-System"},	
			})

			if err != nil {
				log.Println(err)
			}

			trace.RegisterExporter(je)
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

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
