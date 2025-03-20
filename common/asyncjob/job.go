package asyncjob

import (
	"context"
	"log"
	"time"
)

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
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
	StateInit JobState = iota
	StateRunning
	StateFailed
	StateTimeout
	StateCompleted
	StateRetryFailed
)

func (js JobState) String() string {
	return [6]string{"Init", "Running", "Failed", "Timeout", "Completed", "RetryFailed"}[js]
}

type jobConfig struct {
	Name        string
	MaxTimeout  time.Duration
	Retries     []time.Duration
}

type job struct {
	config     jobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler, options ...OptionHdl) *job {
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

	for i := range options {
			options[i](&j.config)
	}

	return &j
}

func (j *job) Execute(ctx context.Context) error {
	log.Printf("execute %s\n", j.config.Name)
	j.state = StateRunning

	var err error
	err = j.handler(ctx)

	if err != nil {
			j.state = StateFailed
			return err
	}

	j.state = StateCompleted

	return nil
}


type OptionHdl func(*jobConfig)

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
