package subscriber

import (
	"context"
	"first-proj/common"
	"first-proj/module/item/storage"
	// "first-proj/module/userlikeitem/model"
	"first-proj/pubsub"
	// "log"

	goservice "github.com/200Lab-Education/go-sdk"
	"gorm.io/gorm"
)

type HasItemID interface {
	GetItemID() int
}

// func IncreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext, ctx context.Context) {
// 	ps := serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)
// 	db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

// 	c, _ := ps.Subscribe(ctx, common.TopicUserLikedItem)

// 	go func() {
// 			defer common.Recovery()
// 			for msg := range c {
// 				data := msg.Data().(*model.Like) 
// 				=> This casts data into lower layer interface to get item id, 
// 				which is not recommeneded
				
			
// 					data := msg.Data().(HasItemID)

// 					if err := storage.NewSQLStore(db).IncreaseLikeCount(ctx, data.GetItemID()); err != nil {
// 							log.Println(err)
// 					}
// 			}
// 	}()
// }

func IncreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
			Title: "Increase like count after user likes item",
			Hld: func(ctx context.Context, message *pubsub.Message) error {
					db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

					data := message.Data().(HasItemID)

					return storage.NewSQLStore(db).IncreaseLikeCount(ctx, data.GetItemID())
			},
	}
}
