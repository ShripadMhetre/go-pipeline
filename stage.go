package gopipeline

import "sync"

type Stage interface {
	Process(stage dispatch) ([]dispatch, error)
}

type StageWorker struct {
	wg          *sync.WaitGroup
	input       chan dispatch
	output      chan dispatch
	concurrency int
	pipe        Stage
}

func NewWorkerGroup(concurrency int, pipe Stage, input chan dispatch, output chan dispatch) StageWorker {
	return StageWorker{
		wg:          &sync.WaitGroup{},
		input:       input,
		output:      output,
		concurrency: concurrency,
		pipe:        pipe,
	}

}
