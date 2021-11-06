package models

import (
	"DatabaseCamp/database"
	"sync"
)

type Concurrent struct {
	Wg  *sync.WaitGroup
	Err *error
}

type ConcurrentTransaction struct {
	Concurrent  *Concurrent
	Transaction database.ITransaction
}
