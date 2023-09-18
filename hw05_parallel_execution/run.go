package main

import (
	"errors"
	"runtime"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCount := len(tasks)
	errorTaskCount := 0
	allHandledCount := 0
	numCPU := runtime.NumCPU()

	completeFlagCh := make(chan struct{})
	tasksCh := make(chan Task, n)

	// ставим размер в 2 раза больший, чем количество системных обработчиков (чтоб уж наверняка)
	errorTaskCh := make(chan error, numCPU*2)

	for i := 0; i < n; i++ {
		go taskHandler(tasksCh, errorTaskCh, completeFlagCh)
	}

	go taskManager(tasks, tasksCh, completeFlagCh)

	for err := range errorTaskCh {
		allHandledCount++

		if err != nil {
			errorTaskCount++
		}

		isErrErrorsLimitExceed := isErrErrorsLimitExceed(m, errorTaskCount)
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

	defer func() {
		close(errorTaskCh)
	}()

	return nil
}

func taskHandler(tasksCh chan Task, errorTaskCh chan error, completeFlagCh chan struct{}) {
	for {
		select {
		case task := <-tasksCh:
			if task != nil {
				errorTaskCh <- task()
			}
		case <-completeFlagCh:
			return
		}
	}
}

func taskManager(tasks []Task, tasksCh chan Task, completeFlagCh chan struct{}) {
	defer close(tasksCh)

	for _, task := range tasks {
		select {
		case tasksCh <- task:
		case <-completeFlagCh:
			return
		}
	}
}

func isErrErrorsLimitExceed(m, errorTaskCount int) bool {
	if m <= 0 && errorTaskCount > 0 || m > 0 && errorTaskCount >= m {
		return true
	}

	return false
}

func completeHandledCount(allHandledCount, tasksCount int) bool {
	return allHandledCount >= tasksCount
}
