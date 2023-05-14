package gopipeline

type dispatch interface {
}

type PipelineOpts struct {
	Concurrency int
}

type Pipeline interface {
	AddPipe(pipe Stage, opt *PipelineOpts)
	Start() error
	Stop() error
	Input() chan<- dispatch
	Output() <-chan dispatch
}

// Implements "Pipeline"
type ConcurrentPipeline struct {
	workerGroups []StageWorker
}
