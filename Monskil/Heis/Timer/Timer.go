package Timer

import (
	"fmt"
	"time"
)

func Timer(timeout chan bool, set_timer chan bool) {
	const door_open_time = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-set_timer:
			fmt.Println("case: set_timer")
			timer.Reset(door_open_time)
		case <-timer.C:
			fmt.Println("case: timer.C")
			timer.Stop()
			fmt.Println("STOP")
			timeout <- true
			fmt.Println("TRUE")
		}
	}
}

// func Timer_external_order(timeout_order chan bool, set_order_timer chan bool) {
// 	const execute_order_time = 15 * time.Second
// 	timer := time.NewTimer(0)
// 	timer.Stop()

// 	for {
// 		select {
// 		case <-set_order_timer:
// 			timer.Reset(execute_order_time)
// 		case <-timer.C:
// 			timer.Stop()
// 			timeout_order <- true
// 		}
// 	}
// }
