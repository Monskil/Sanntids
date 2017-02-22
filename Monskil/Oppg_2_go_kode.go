package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i int = 0

func someGoRoutine1(ch chan int, chan_done1 chan bool) {
	for x := 0; x < 1000000; x++ {
		i = <-ch
		i += 1
		ch <- i

	}
	chan_done1 <- true
}

func someGoRoutine2(ch chan int, chan_done2 chan bool) {
	for x := 0; x < 1000003; x++ {
		i = <-ch
		i -= 1
		ch <- i
	}
	chan_done2 <- true
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ch := make(chan int, 1)
	ch <- i

	chan_done1 := make(chan bool, 1)
	chan_done2 := make(chan bool, 1)

	go someGoRoutine1(ch, chan_done1)
	go someGoRoutine2(ch, chan_done2)

	<-chan_done1
	<-chan_done2

	//time.Sleep(100 * time.Millisecond)

	Println(i)
}
