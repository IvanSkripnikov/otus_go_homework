package main

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCount := len(tasks)
	errorTaskCount := 0
	allHandledCount := 0

	completeFlagCh := make(chan struct{})
	tasksCh := make(chan Task, n)
	errorTaskCh := make(chan error, tasksCount)

	for i := 0; i < n; i++ {
		go taskHandler(tasksCh, errorTaskCh, completeFlagCh)
	}

	go taskManager(tasks, tasksCh, completeFlagCh)

	for err := range errorTaskCh {
		allHandledCount++

		if err != nil {
			errorTaskCount++
		}

		// Проверяем на лимит по ошибкам
		if m <= 0 && errorTaskCount > 0 || m > 0 && errorTaskCount >= m {
			close(completeFlagCh)
			return ErrErrorsLimitExceeded
		}

		// Защищаемся от вечного цикла
		if allHandledCount >= tasksCount {
			close(completeFlagCh)
			break
		}
	}

	close(errorTaskCh)

	return nil
}

func taskHandler(tasksCh chan Task, errorsCh chan error, completeFlagCh chan struct{}) {
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
