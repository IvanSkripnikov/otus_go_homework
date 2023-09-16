package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorCount := 0
	allHandledCount := 0
	maxPossibleCount := n + m
	tasksCh := make(chan Task, n)
	errorsCh := make(chan error, len(tasks))

	for i := 0; i < n; i++ {
		go taskHandler(n, tasksCh, errorsCh)
	}

	go taskManager(tasks, tasksCh)

	time.Sleep(1 * time.Second)

	for err := range errorsCh {
		allHandledCount++

		if err != nil {
			errorCount++
		}

		if allHandledCount > maxPossibleCount || errorCount > m {
			return ErrErrorsLimitExceeded
		}
	}

	close(errorsCh)

	return nil
}

func taskHandler(n int, tasksCh chan Task, errorsCh chan error) {
	wg := sync.WaitGroup{}
	wg.Add(n)

	task := <-tasksCh
	if task != nil {
		errorsCh <- task()
	}

	wg.Wait()
}

func taskManager(tasks []Task, tasksCh chan Task) {
	for _, task := range tasks {
		tasksCh <- task
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
