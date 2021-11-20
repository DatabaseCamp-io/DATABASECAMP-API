package general

import (
	"DatabaseCamp/database"
	"sync"
)

type Concurrent struct {
	Wg    *sync.WaitGroup
	Err   *error
	Mutex *sync.Mutex
}

type ConcurrentTransaction struct {
	Concurrent  *Concurrent
	Transaction database.ITransaction
}
