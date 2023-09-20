package main

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	lessOrZeroM := m <= 0
	tasksCount := len(tasks)
	errorTaskCount := 0
	allHandledCount := 0
	tasksCh := make(chan Task, tasksCount)

	mu := sync.Mutex{}
	mu.Lock()
	// аписываем в канал задачи
	go taskProducer(tasks, tasksCh)
	mu.Unlock()

	mu.Lock()
	// заполняем канал с результататми выполнения задач
	errorTaskCh := make(chan error, tasksCount)
	for i := 0; i < n; i++ {
		go taskConsumer(tasksCh, errorTaskCh)
	}
	mu.Unlock()

	// проверяем канал с результатами выполнения задач
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

		// если число ошибок превысило допустимый лимит возвращаем ошибку
		if isErrErrorsLimitExceed {
			return ErrErrorsLimitExceeded
		}

		// если число успешно выполненных задач достигло необходимого количества - завершаем считывание из канала
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
		if !ok || task == nil {
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
