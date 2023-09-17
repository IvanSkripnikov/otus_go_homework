package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorTaskCount := 0
	allHandledCount := 0
	maxPossibleCount := n + m

	completeFlagCh := make(chan struct{})
	tasksCh := make(chan Task, n)
	errorTaskCh := make(chan error, len(tasks))

	for i := 0; i < n; i++ {
		go taskHandler(tasksCh, errorTaskCh, completeFlagCh)
	}

	go taskManager(tasks, tasksCh, completeFlagCh)

	time.Sleep(1 * time.Second)

	for err := range errorTaskCh {
		allHandledCount++

		if err != nil {
			errorTaskCount++
		}

		if allHandledCount > maxPossibleCount || errorTaskCount > m {
			return ErrErrorsLimitExceeded
		}

		if allHandledCount >= len(tasks) {
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

func main() {
	tasksCount := 5
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
	}

	workersCount := 5
	maxErrorsCount := 2
	err := Run(tasks, workersCount, maxErrorsCount)
	fmt.Println(err)
}
