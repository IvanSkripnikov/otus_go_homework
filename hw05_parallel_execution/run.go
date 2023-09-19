package main

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	lessOrZeroM := m <= 0
	tasksCount := len(tasks)
	errorTaskCount := 0
	allHandledCount := 0
	tasksCh := make(chan Task, n)

	go taskProducer(tasks, tasksCh)

	errorTaskCh := make(chan error, n)
	for i := 0; i < n; i++ {
		go taskConsumer(tasksCh, errorTaskCh)
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

func taskConsumer(tasksCh chan Task, errorTaskCh chan error) {
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}
		errorTaskCh <- task()
	}
}

func taskProducer(tasks []Task, tasksCh chan Task) {
	defer close(tasksCh)

	for _, task := range tasks {
		tasksCh <- task
	}
}

func isErrErrorsLimitExceed(m, errorTaskCount int, lessOrZeroM bool) bool {
	return lessOrZeroM || errorTaskCount >= m
}

func completeHandledCount(allHandledCount, tasksCount int) bool {
	return allHandledCount >= tasksCount
}
