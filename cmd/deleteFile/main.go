package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	err := os.Remove("static/img/headPhoto/2cd0cb2f-1569-4e0f-aff3-15c16d7b7b97.jpg")
	if err != nil {
		fmt.Printf("remove ./file1.txt err : %v\n", err)
	}

}

//定时任务
func main2() {
	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second *1)
	go func() {
		i := 0
		for range ticker.C {
			i++
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			if i == 5 {
				ch <- 1
			}
		}

	}()
	<-ch
}
