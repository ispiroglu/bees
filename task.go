package bees

import (
	"log"
	"sync"
)

// There might be multiple type of Tasks. In order to cover
// all of them, we use interface.
type ITask interface {
	Execute() error
	OnFailure(error)
}

// This tasks can be implemented in client side maybe?
// This way, every client can provide its own structure to use
// worker pool
type Task struct {
	fn func() error // What should be the template of this fn field?
	wg *sync.WaitGroup
}

func (t *Task) Execute() error {
	defer t.wg.Done()
	return t.fn()
}

func (t *Task) OnFailure(err error) {
	log.Printf("Got error on worker pool, %v", err)
}
