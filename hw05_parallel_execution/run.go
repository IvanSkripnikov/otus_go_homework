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
	completeFlagCh := make(chan struct{})
	tasksCh := make(chan Task, n)

	go taskProducer(tasks, tasksCh, completeFlagCh)

	errorTaskCh := make(chan error, numCPU*2)
	for i := 0; i < n; i++ {
		go taskConsumer(tasksCh, errorTaskCh, completeFlagCh)
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

		close(completeFlagCh)

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

func taskConsumer(tasksCh chan Task, errorsCh chan error, completeFlagCh chan struct{}) {
	for {
		select {
		case task := <-tasksCh:
			if task != nil {
				errorsCh <- task()
			}
		case <-completeFlagCh:
			return
		}
	}
}

func taskProducer(tasks []Task, tasksCh chan Task, completeFlagCh chan struct{}) {
	defer close(tasksCh)

	for _, task := range tasks {
		select {
		case tasksCh <- task:
		case <-completeFlagCh:
			return
		}
	}
}

func isErrErrorsLimitExceed(m, errorTaskCount int, lessOrZeroM bool) bool {
	return lessOrZeroM || errorTaskCount >= m
}

func completeHandledCount(allHandledCount, tasksCount int) bool {
	return allHandledCount >= tasksCount
}
