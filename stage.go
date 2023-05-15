package gopipeline

import (
	"log"
	"sync"
)

type Stage interface {
	Process(stage Dispatch) ([]Dispatch, error)
}

type StageWorker struct {
	wg          *sync.WaitGroup
	in          chan Dispatch
	out         chan Dispatch
	concurrency int
	stage       Stage
}

// Constructor function - StageWorker
func NewStageWorker(concurrency int, stage Stage, input chan Dispatch, output chan Dispatch) StageWorker {
	return StageWorker{
		wg:          &sync.WaitGroup{},
		in:          input,
		out:         output,
		concurrency: concurrency,
		stage:       stage,
	}

}

func (w *StageWorker) Start() error {

	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)

		go func() {
			defer w.wg.Done()
			for i := range w.Input() {
				result, err := w.stage.Process(i)
				if err != nil {
					log.Println(err)
					continue
				}
				for _, r := range result {
					w.Output() <- r
				}
			}
		}()
	}

	return nil

}

func (w *StageWorker) WaitStop() error {
	w.wg.Wait()
	return nil
}

func (w *StageWorker) Input() chan Dispatch {
	return w.in
}

func (w *StageWorker) Output() chan Dispatch {
	return w.out
}
