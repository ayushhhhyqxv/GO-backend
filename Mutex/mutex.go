package main

import (
	"fmt"
	"sync"
)

var value int64
var mutex sync.Mutex

func lock(wg *sync.WaitGroup){
	mutex.Lock()
	value++
	defer mutex.Unlock()

	wg.Done()
}

func main() {
	var wg sync.WaitGroup

	const n = 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		go lock(&wg)
	}
	wg.Wait()

	fmt.Println("value:", value)
}