package main

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	emptyResult := make(Bi)
	close(emptyResult)

	// если в пайплайне нет этапов, выходим из функции
	if len(stages) == 0 {
		return emptyResult
	}

	for _, stage := range stages {
		if done != nil {
			result := make(Bi)
			// выполнение функции и переход на следующий этап
			go execute(out, done, result)
			out = stage(result)
		} else {
			out = stage(out)
		}
	}

	return out
}

func execute(out Out, done In, result Bi) {
	defer close(result)

	for {
		select {
		case <-done:
			return

		case res := <-out:
			select {
			case <-done:
				return

			case result <- res:
			}
		}
	}
}
