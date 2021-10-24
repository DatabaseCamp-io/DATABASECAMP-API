package models

import "sync"

type Concurrent struct {
	Wg  *sync.WaitGroup
	Err *error
}
