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

func DecreaseLikeCountAfterUserUnlikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
			Title: "Decrease like count after user unlikes item",
			Hld: func(ctx context.Context, message *pubsub.Message) error {
					db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

					data := message.Data().(HasItemID)

					return storage.NewSQLStore(db).DecreaseLikeCount(ctx, data.GetItemID())
			},
	}
}
