package main

import (
	"fmt"
	"sync"
)

var counter int64

func count(wg *sync.WaitGroup){
	counter++
	wg.Done()
}

func main(){
	var wg sync.WaitGroup

	for i:=1;i<=1000;i++{
		wg.Add(1)
		go count(&wg)
	}
	wg.Wait()
	fmt.Println("Done Executing Go Routines !: ",counter)
}

// Here output may vary because operation by goroutine is'nt atomic 