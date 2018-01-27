package poloniex

import (
	"fmt"
	"sync"
)

type Logger struct {
	isOpen bool
	Lock   *sync.Mutex
}

func (l *Logger) LogRoutine(bus <-chan string) {
	if l.isOpen {
		for {
			message := <-bus
			l.Lock.Lock()
			fmt.Println(message)
			l.Lock.Unlock()
		}
	}
}
