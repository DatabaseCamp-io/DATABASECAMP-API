package loaders

import "sync"

type Concurrent struct {
	Wg    *sync.WaitGroup
	Err   *error
	Mutex *sync.Mutex
}
