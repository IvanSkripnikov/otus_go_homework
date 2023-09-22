package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
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

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, atomic.LoadInt32(&runTasksCount), int32(workersCount+maxErrorsCount),
			"extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, atomic.LoadInt32(&runTasksCount), int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks with maximum errors 0", func(t *testing.T) {
		tasksCount := 20
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(170)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 7
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, atomic.LoadInt32(&runTasksCount), int32(workersCount+maxErrorsCount),
			"extra tasks were started")
	})

	t.Run("tasks with negative maximum errors", func(t *testing.T) {
		tasksCount := 30
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(150)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 7
		maxErrorsCount := -1
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, atomic.LoadInt32(&runTasksCount), int32(workersCount+maxErrorsCount),
			"extra tasks were started")
	})

	t.Run("tasks with random count errors and success tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		rand.Seed(time.Now().UnixNano())

		var runSuccessCount int32
		var runErrorCount int32
		maxErrorsCount := 5

		for i := 0; i < tasksCount; i++ {
			if i < maxErrorsCount {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(170)))
					atomic.AddInt32(&runErrorCount, 1)

					return fmt.Errorf("error from task %d", i)
				})
				continue
			}
			randNumber := rand.Intn(100) % 2

			if randNumber == 0 {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(170)))
					atomic.AddInt32(&runErrorCount, 1)

					return fmt.Errorf("error from task %d", i)
				})
			} else {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(170)))
					atomic.AddInt32(&runSuccessCount, 1)

					return nil
				})
			}
		}

		workersCount := 10

		err := Run(tasks, workersCount, maxErrorsCount)
		require.True(t, errors.Is(err, ErrErrorsLimitExceeded))
		require.LessOrEqual(t, atomic.LoadInt32(&runSuccessCount), int32(workersCount+maxErrorsCount),
			"extra tasks were started")
		require.LessOrEqual(t, int32(maxErrorsCount), atomic.LoadInt32(&runErrorCount),
			"extra tasks were started")
	})
}
