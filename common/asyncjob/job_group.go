package asyncjob

import (
	"context"
	"first-proj/common"
	"log"
	"sync"
)

type group struct {
	jobs         []Job
	isConcurrent bool
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
			isConcurrent: isConcurrent,
			jobs:         jobs,
			wg:           new(sync.WaitGroup),
	}
	return g
}

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	//buffer channel which collects error
	errChan := make(chan error, len(g.jobs))

	for i := range g.jobs {
			if g.isConcurrent {
				// Each job is launched in its own goroutine.
					go func(aj Job) {
							defer common.Recovery()

							// Run the job and send its error (if any) to errChan.
							errChan <- g.runJob(ctx, aj) 
							
							// Signal that this job is finished.
							g.wg.Done()
					}(g.jobs[i])

					continue
			}

			job := g.jobs[i]

			err := g.runJob(ctx, job)

			if err != nil {
					return err
			}

			errChan <- err 
			g.wg.Done()
		}

		
			g.wg.Wait()

			var err error

			for i := 1; i <= len(g.jobs); i++ {
					if v := <-errChan; v != nil {
							return v
					}
			}

		return err
}

// Execute the job, retry if failed, 
// and return error if reaches maximum number of retries allowed
func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
			for {
					log.Println(err)
					if j.State() == StateRetryFailed {
							return err
					}
			}

			// if j.Retry(ctx) == nil {
			// 		return nil
			// }
	}

	return nil
}
