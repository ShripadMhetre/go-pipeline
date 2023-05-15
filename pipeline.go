package gopipeline

import "errors"

type Dispatch interface {
}

type PipelineOptions struct {
	Concurrency int
}

type Pipeline interface {
	AddStage(pipe Stage, optss *PipelineOptions)
	Start() error
	Stop() error
	Input() chan<- Dispatch
	Output() <-chan Dispatch
}

var ErrPipelineEmpty = errors.New("concurrent pipeline empty")

// Implements "Pipeline"
type ConcurrentPipeline struct {
	workers []StageWorker
}

// Constructor function - ConcurrentPipeline
func New() Pipeline {
	return &ConcurrentPipeline{}
}

func (c *ConcurrentPipeline) AddStage(pipe Stage, opts *PipelineOptions) {
	if opts == nil {
		opts = &PipelineOptions{Concurrency: 1}
	}

	var input = make(chan Dispatch, 10)
	var output = make(chan Dispatch, 10)

	for _, i := range c.workers {
		input = i.Output()
	}

	worker := NewStageWorker(opts.Concurrency, pipe, input, output)
	c.workers = append(c.workers, worker)
}

func (c *ConcurrentPipeline) Output() <-chan Dispatch {
	sz := len(c.workers)
	return c.workers[sz-1].Output()
}

func (c *ConcurrentPipeline) Input() chan<- Dispatch {
	return c.workers[0].Input()
}

// Starts running the pipeline
func (c *ConcurrentPipeline) Start() error {
	if len(c.workers) == 0 {
		return ErrPipelineEmpty
	}

	for i := 0; i < len(c.workers); i++ {
		g := c.workers[i]
		g.Start()
	}

	return nil
}

// Stop running the pipeline
func (c *ConcurrentPipeline) Stop() error {
	for _, i := range c.workers {
		close(i.Input())
		i.WaitStop()
	}

	sz := len(c.workers)
	close(c.workers[sz-1].Output())
	return nil
}
