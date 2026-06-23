package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"time"
)


func worker(id int){
	trace.WithRegion(context.Background(),fmt.Sprintf("Worker: %d",id),func(){
		fmt.Printf("goRoutine %d started\n",id)
		for i:= 1; i<=3 ; i++ {
			fmt.Printf("goRoutine %d and Working iteration %d\n",id,i)
		}
		fmt.Printf("goRoutine %d executed\n",id)
	})
}

func main(){
	runtime.GOMAXPROCS(2)

	f,err := os.Create("trace.out")

	if err!=nil {
		panic(err)
	}

	fmt.Println("Tracing the events")

	if err:= trace.Start(f);err!=nil {
		panic(err)
	}

	defer f.Close()
	defer trace.Stop()

	trace.WithRegion(context.Background(),"main",func(){
		fmt.Println("Main routine started")

		for i:=1;i<=5;i++ {
			go worker(i)
		}

		for i:=1;i<=3;i++ {
			fmt.Printf("Main goRoutine working iteration %d\n",i)
			time.Sleep(200*time.Millisecond)
		}
		fmt.Println("Execution Done")
	})

}