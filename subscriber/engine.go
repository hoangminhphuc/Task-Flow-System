package subscriber

import (
	"context"
	"first-proj/common"
	"first-proj/common/asyncjob"
	"first-proj/pubsub"
	"log"

	goservice "github.com/200Lab-Education/go-sdk"
)

type subJob struct {
	Title string
	Hld   func(ctx context.Context, message *pubsub.Message) error
}

type pbEngine struct {
	serviceCtx goservice.ServiceContext
}

func NewEngine(serviceCtx goservice.ServiceContext) *pbEngine {
	return &pbEngine{serviceCtx: serviceCtx}
}


// This is scalable since when project grows, it will have so many subscribers
// At that point, we just have to pass them as arguments into this function
func (engine *pbEngine) Start() error {
	engine.startSubTopic(common.TopicUserLikedItem, true, 
		IncreaseLikeCountAfterUserLikeItem(engine.serviceCtx),
	)

	engine.startSubTopic(common.TopicUserUnlikedItem, true, 
		DecreaseLikeCountAfterUserUnlikeItem(engine.serviceCtx),
	)

	return nil
}

type GroupJob interface {
	Run(ctx context.Context) error
}


func (engine *pbEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, jobs ...subJob) error {
	ps := engine.serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)


	// c is a channel for receiving messages from broker
	c, _ := ps.Subscribe(context.Background(), topic)

	for _, item := range jobs {
			log.Println("Setup subscriber for: ", item.Title)
	}


  // ! Converts a subJob into an asyncjob.JobHandler.
	/* subJob has a handler that takes context and a message; 
		while a job only takes context not a message. */

	/*
	* Why does we need to convert ?
		- Because a job doesnt know what the message is. It needs message so that it
			knows when to execute. 
	*/
	getJobHandler := func(job *subJob, message *pubsub.Message) asyncjob.JobHandler {
			return func(ctx context.Context) error {
					log.Println("running job for", job.Title, ". Value:", message.Data())
				
				// executes subJob
					return job.Hld(ctx, message)
			}
	}

// This will be initialized at the beginning, when service are started.
	go func() {
    for {
			//Listening for message, blocks the rest when receive a message from broker
        msg := <- c
        jobHdlArr := make([]asyncjob.Job, len(jobs))

			/* 
			* Convert each subJob to an asyncjob.Job
			 */
        for i := range jobs {
					// jobHdl is a function of type asyncjob.JobHandler
					// inside this function, it executes the subJob normally
            jobHdl := getJobHandler(&jobs[i], msg)
            jobHdlArr[i] = asyncjob.NewJob(jobHdl, asyncjob.WithName(jobs[i].Title))
        }

			// Runs job group
        group := asyncjob.NewGroup(isConcurrent, jobHdlArr...)

        if err := group.Run(context.Background()); err != nil {
            log.Println(err)
        }
			}
		}()

	return nil

}

