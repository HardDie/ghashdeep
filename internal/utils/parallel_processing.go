package utils

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
)

type processingFunc func(m *sync.Mutex) error

// Processing Wrapper for running downloads in parallel
type Processing struct {
	// List of functions to start
	routines map[string]processingFunc
	// Process tasks one by one
	singleThread bool
}

func NewProcessing() *Processing {
	return &Processing{
		routines: make(map[string]processingFunc),
	}
}

// Push Add a new function to start
func (pr *Processing) Push(name string, fn processingFunc) {
	pr.routines[name] = fn
}

// SingleThread Set the flag so that tasks are executed sequentially
func (pr *Processing) SingleThread() {
	pr.singleThread = true
}

// Run Start all added functions and check them for errors
func (pr *Processing) Run() error {
	m := sync.Mutex{}
	var errs []error

	maxRoutines := runtime.NumCPU()
	limiter := make(chan struct{}, maxRoutines)

	var wg sync.WaitGroup
	for name, fn := range pr.routines {
		if pr.singleThread {
			if err := fn(&m); err != nil {
				errs = append(errs, fmt.Errorf("%s %w", name, err))
			}
			continue
		}

		// Do not run new goroutine if all CPU is busy
		limiter <- struct{}{}
		wg.Add(1)
		go func(name string, fn processingFunc) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			if err := fn(&m); err != nil {
				m.Lock()
				errs = append(errs, fmt.Errorf("[%s] %w", name, err))
				m.Unlock()
			}
		}(name, fn)
	}
	wg.Wait()

	return errors.Join(errs...)
}
