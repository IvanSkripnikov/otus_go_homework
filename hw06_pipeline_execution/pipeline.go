package main

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var breakOut = make(Bi)
	// Если отсутствуют тапы пайплайна - ничего не делаем
	if len(stages) == 0 {
		return breakOut
	}

	var out Out

	for _, stage := range stages {
		if done != nil {

			for {
				select {
				case <-done:
					return breakOut
				default:
					out = stage(out)
					break
				}
			}

		} else {
			out = stage(out)
		}
	}

	return out
}
