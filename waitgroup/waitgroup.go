package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	for _, item := range [5]string{"h", "e", "l", "l", "o"} {
		wg.Add(1)
		go func(item string, wg *sync.WaitGroup) {
			defer wg.Done()
			fmt.Printf("Doing %s\n", item)
		}(item, &wg)
	}
	fmt.Println("Waiting...")
	wg.Wait()
	fmt.Println("Done...")
}
