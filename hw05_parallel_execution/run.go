package main

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	lessOrZeroM := m <= 0
	tasksCount := len(tasks)
	errorTaskCount := 0
	allHandledCount := 0
	var flagNotWrite int32
	tasksCh := make(chan Task)

	// аписываем в канал задачи
	go taskProducer(tasks, tasksCh, flagNotWrite)

	// заполняем канал с результататми выполнения задач
	errorTaskCh := make(chan error, tasksCount)
	for i := 0; i < n; i++ {
		go taskConsumer(tasksCh, errorTaskCh, flagNotWrite)
	}

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

		// если же у нас выполнилось необходимое количество задач или произошло ошибок - выставляем флаг
		atomic.AddInt32(&flagNotWrite, 1)

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

func taskConsumer(tasksCh chan Task, errorTaskCh chan error, flagNotWrite int32) {
	for {
		// если выставлен флаг на окончание - не даём записывать в канал
		if atomic.LoadInt32(&flagNotWrite) > 0 {
			return
		}
		task, ok := <-tasksCh
		if !ok || task == nil {
			return
		}
		errorTaskCh <- task()
	}
}

func taskProducer(tasks []Task, tasksCh chan Task, flagNotWrite int32) {
	defer close(tasksCh)
	// если выставлен флаг на окончание - не даём записывать в канал
	if atomic.LoadInt32(&flagNotWrite) > 0 {
		return
	}

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
