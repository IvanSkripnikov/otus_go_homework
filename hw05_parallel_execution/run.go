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
	errorCount := 0

	for i := 0; i < n; i++ {
		go taskHandler(tasks, &errorCount)
	}
	fmt.Println(errorCount)
	time.Sleep(2 * time.Second)

	return nil
}

func taskHandler(tasks []Task, errorCount *int) {
	fmt.Println(tasks)
	*errorCount++
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
