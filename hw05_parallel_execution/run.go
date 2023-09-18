package main

import (
	"errors"
	"runtime"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	lessOrZeroM := m <= 0
	tasksCount := len(tasks)
	errorTaskCount := 0
	allHandledCount := 0
	numCPU := runtime.NumCPU()
	tasksCh := make(chan Task, n)

	go taskManager(tasks, tasksCh)

	// ставим размер в 2 раза больший, чем количество системных обработчиков (чтоб уж наверняка)
	errorTaskCh := make(chan error, numCPU*2)
	for i := 0; i < n; i++ {
		go taskHandler(tasksCh, errorTaskCh)
	}

	for err := range errorTaskCh {
		allHandledCount++

		if err != nil {
			errorTaskCount++
		}

		isErrErrorsLimitExceed := isErrErrorsLimitExceed(m, errorTaskCount, lessOrZeroM)
		completeHandledCount := completeHandledCount(allHandledCount, tasksCount)

		// пока не выполнилось корнер условие пропускаем итерацию
		if !isErrErrorsLimitExceed && !completeHandledCount {
			continue
		}

		if isErrErrorsLimitExceed {
			return ErrErrorsLimitExceeded
		}

		if completeHandledCount {
			break
		}
	}

	defer close(errorTaskCh)

	return nil
}

func taskHandler(tasksCh chan Task, errorTaskCh chan error) {
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}
		errorTaskCh <- task()
	}
}

func taskManager(tasks []Task, tasksCh chan Task) {
	defer close(tasksCh)

	for _, task := range tasks {
		tasksCh <- task
	}
}

func isErrErrorsLimitExceed(m, errorTaskCount int, lessOrZeroM bool) bool {
	if lessOrZeroM || errorTaskCount >= m {
		return true
	}

	return false
}

func completeHandledCount(allHandledCount, tasksCount int) bool {
	return allHandledCount >= tasksCount
}
