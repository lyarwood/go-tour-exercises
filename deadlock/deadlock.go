package main

// https://medium.com/golangspec/detect-locks-passed-by-value-in-go-efb4ac9a3f2b

import (
	"fmt"
	"sync"
	"time"
)

type T struct {
	lock *sync.Mutex
}

func (t *T) Lock() {
	t.lock.Lock()
}
func (t *T) Unlock() {
	t.lock.Unlock()
}
func main() {
	t := T{lock: &sync.Mutex{}}
	c := make(chan bool)
	defer close(c)

	go func(t *T, c chan bool) {
		t.Lock()
		fmt.Println("locked inside routine.....")
		time.Sleep(1 * time.Second)
		t.Unlock()
		fmt.Println("unlocked inside routine.....")
		c <- true
	}(&t, c)
	fmt.Println("Waiting outside routune....")
	<-c
	t.Lock()
	fmt.Println("Locked outside routune....")
	t.Unlock()
	fmt.Println("Unlocked outside routune....")
}
