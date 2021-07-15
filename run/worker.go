package run

import (
	"runtime"
	"sync"
)

type taskParams struct {
	f               FieldI
	removeCombIndex int
}

type taskReturn struct {
	f FieldI
}

// Workers takes inputs on a channel and runs a task on those inputs.
// Given a task and an output channel, Workers will return an input channel on which to send inputs.
// The tasks will be run concurrently, when they are done, the output channel is closed.
// todo: add a context to kill the goroutines early
func Workers(task func(*taskParams), output chan *taskReturn) chan *taskParams {
	// create channels
	inputs := make(chan *taskParams)
	goRoutines := runtime.NumCPU()

	// Start launching goroutines
	// Each goroutine waits for inputs on the input channel
	// When one arrives, run the task on that input
	// It's up to the task to return any results to the output channel
	// When the input channel is closed, send true to the done channel indicating all tasks are done
	wg := sync.WaitGroup{}
	wg.Add(goRoutines)
	for i := 0; i < goRoutines; i++ {
		go func() {
			for input := range inputs {
				task(input)
			}
			wg.Done()
		}()
	}
	// this goroutine waits until all tasks are done, then closes the output channel
	go func() {
		wg.Wait()
		close(output)
	}()

	return inputs
}
