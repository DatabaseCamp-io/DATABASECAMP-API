package general

// controller.go
/**
 * 	This file is a part of models, used to collect model for controller
 */

import (
	"DatabaseCamp/database"
	"sync"
)

// Model for do concurrent
type Concurrent struct {
	Wg    *sync.WaitGroup
	Err   *error
	Mutex *sync.Mutex
}

// Model for do concurrent with transaction
type ConcurrentTransaction struct {
	Concurrent  *Concurrent
	Transaction database.ITransaction
}
