package asyncjob

import (
	"context"
	"log"
	"time"
)


/* 
! Running tasks/jobs asynchronously with built-in mechanisms
*/

type Job interface {
	Execute(ctx context.Context) error // runs the job
	Retry(ctx context.Context) error // try to run the job again
	State() JobState // returns current state 
	SetRetryDurations(times []time.Duration)
}

const (
	defaultMaxTimeout = time.Second * 10
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 2, time.Second * 4}
)

type JobHandler func(ctx context.Context) error

type JobState int

const (
	StateInit JobState = iota // 1, 2, 3...
	StateRunning
	StateFailed
	StateTimeout
	StateCompleted
	StateRetryFailed
)

func (js JobState) String() string {
	//creates an array and indexing into that array with js integer
	return [6]string{"Init", "Running", "Failed", "Timeout", "Completed", "RetryFailed"}[js]
}

type jobConfig struct {
	Name        string
	MaxTimeout  time.Duration 
	//Slice of durations that define the waiting period between retries.
	Retries     []time.Duration 
}

type job struct {
	config     jobConfig 
	handler    JobHandler // A function that actually performs the job work.
	state      JobState 
	/* Tracks which retry duration to use next 
	(initialized at -1 to indicate no attempt has been made yet). */
	retryIndex int
	stopChan   chan bool
}

/* 
* This way of writing constructor is easier to maintain
* than when we accept a bunch of arguments
*/
func NewJob(handler JobHandler, options ...OptionHdl) *job {
	 // Create a default configuration.
	j := job{
			config: jobConfig{
					MaxTimeout: defaultMaxTimeout,
					Retries:    defaultRetryTime,
			},
			handler:    handler,
			retryIndex: -1,
			state:      StateInit,
			stopChan:   make(chan bool),
	}

	/* 
	* Because options are slice of OptionHdl (which is a function that accepts a 
	* job config as a param), it applies each option function to modify the configuration.
	*/
	for i := range options {
			options[i](&j.config)
	}

	return &j
}

func (j *job) Execute(ctx context.Context) error {
	log.Printf("execute %s\n", j.config.Name)
	j.state = StateRunning

	var err error
	//Run the handler which is the actual function that do thing
	err = j.handler(ctx)

	if err != nil {
			j.state = StateFailed
			return err
	}

	j.state = StateCompleted

	return nil
}

func (j *job) Retry(ctx context.Context) error {
	// tracks how many retries have been attempted
	// if reach the maximum number of retries, return nil
	// if j.retryIndex == len(j.config.Retries)-1 {
	// 		return nil // TODO: we should save the last error of execution
	// }

	j.retryIndex += 1
	/* 
* Backoff strategy: If a job fails, immediately retry might just hit the same error.
* Waiting a little while gives the system a chance to recover.
	*/
	time.Sleep(j.config.Retries[j.retryIndex])

	err := j.Execute(ctx)

	if err == nil {
			j.state = StateCompleted
			return nil
	}

	if j.retryIndex == len(j.config.Retries)-1 {
			j.state = StateRetryFailed
			return err
	}

	j.state = StateFailed
	return err
}

func (j *job) State() JobState { return j.state }
func (j *job) RetryIndex() int { return j.retryIndex }

func (j *job) SetRetryDurations(times []time.Duration) {
    if len(times) == 0 {
        return
    }
    j.config.Retries = times
}


type OptionHdl func(*jobConfig)


/* 
	* Returns an option function that, when called, sets the name of the job
*/
func WithName(name string) OptionHdl {
	return func(cf *jobConfig) {
			cf.Name = name
	}
}

func WithRetriesDuration(items []time.Duration) OptionHdl {
	return func(cf *jobConfig) {
			cf.Retries = items
	}
}
