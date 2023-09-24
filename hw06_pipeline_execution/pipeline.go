package main

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out
	for key, stage := range stages {
		if key == 0 {
			out = stage(in)
		} else {
			out = stage(out)
		}
	}

	if out != nil {
		return out
	}

	return nil
}
