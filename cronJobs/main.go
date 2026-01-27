package main

import (
    "fmt"
    "time"
    "github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds()) // cron scheduler !

	c.AddFunc("* * * * * *", func() {
		fmt.Println("Function is Running",time.Now().Format("15:04:05"))
	})

	c.Start()

	select {


	}
}