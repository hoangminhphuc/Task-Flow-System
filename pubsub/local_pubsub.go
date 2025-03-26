package pubsub

import (
	"context"
	"first-proj/common"
	"log"
	"sync"
)

/* 
* Local Pub/Sub (in-mem) means:
*   - All publishers and subscribers share the same process and memory space.
*   - No external broker (like Redis or Kafka) is used.
*/
// It has a queue (buffer channel) at its core and many groups of subscribers.
// Because we want to send a message with a specific topic for many subscribers in a group to handle.
type localPubSub struct {
	name 				 string

	// Messages flow from Publish → messageQueue → subscribers.
	messageQueue chan *Message

	// A map where the key is a Topic,
	// and the value is a slice of channels for subscribers who care about that topic.
	mapChannel   map[Topic][]chan *Message

	// Multi readers, 1 writes only
	locker       *sync.RWMutex
}

func NewPubSub(name string) *localPubSub {
	pb := &localPubSub{
		name:         name,
		messageQueue: make(chan *Message, 10000),
		mapChannel:   make(map[Topic][]chan *Message),
		locker:       new(sync.RWMutex),
	}

	// Starts reading messages from queue
	// pb.run()

	return pb
}

func (ps *localPubSub) Publish(ctx context.Context, topic Topic, data *Message) error {
	
	// Before publishing, the message’s channel is set.
	data.SetChannel(topic)

	go func() {
			defer common.Recovery()
			//sent real data (or message) to queue
			ps.messageQueue <- data
			log.Println("New message published:", data.String())
	}()

	return nil
}

/* 
! Each time this function called, it generates new channel specificially to 1 subscriber
*/
func (ps *localPubSub) Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, unsubscribe func()) {
	
	c := make(chan *Message)

	// Only 1 can write (append subscriber channel c)
	ps.locker.Lock()


// Check if there's already an entry for the topic in the mapChannel map.
	if val, ok := ps.mapChannel[topic]; ok {
		// appends the new subscriber channel c to the slice.
			val = append(ps.mapChannel[topic], c)
			ps.mapChannel[topic] = val
	} else {
		// If there are no subscribers for this topic yet, 
		// creates a new slice containing just the new subscriber channel c
			ps.mapChannel[topic] = []chan *Message{c}
	}

	ps.locker.Unlock()

// Channel c is returned as read-only so subscribers can receive messages.
	return c, func() {
			log.Println("Unsubscribe")

			if chans, ok := ps.mapChannel[topic]; ok {
					for i := range chans {
							if chans[i] == c {
									chans = append(chans[:i], chans[i+1:]...)

									ps.locker.Lock()
									ps.mapChannel[topic] = chans
									ps.locker.Unlock()

									close(c)
									break
							}
					}
			}
	}
}


func (ps *localPubSub) run() error {
	go func() {
			defer common.Recovery()
			for {
				// waits (blocks) until a message is available in the channel.
					mess := <-ps.messageQueue
					log.Println("Message dequeue:", mess.String())

					ps.locker.RLock()
					
					// mess.Channel() retrieves the Topic for the message.
					// ps.mapChannel map is checked for a key matching that topic.
					// if ok, subs will hold a slice of subscribers.
					if subs, ok := ps.mapChannel[mess.Channel()]; ok {

					// Loops through all channels of subscribers
						for i := range subs {
							go func(c chan *Message) {
								defer common.Recovery()
								
								//The subscriber's channel receives the message.
								c <- mess
								}(subs[i]) // Pass like this will avoid closure capture problem.
							}
						}


					ps.locker.RUnlock()
			}
	}()
	
	return nil
}

func (ps *localPubSub) GetPrefix() string {
	return ps.name
}

func (ps *localPubSub) Get() interface{} {
	return ps
}

func (ps *localPubSub) Name() string {
	return ps.name
}

func (*localPubSub) InitFlags() {
}

func (*localPubSub) Configure() error {
	return nil
}

func (ps *localPubSub) Run() error {
	return ps.run()
}

func (*localPubSub) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
			c <- true
	}()
	return c
}
