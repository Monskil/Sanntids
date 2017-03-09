package Timer

import (
	"../Driver"
	//"fmt"
	"time"
)

func Timer(timeout chan bool, set_timer chan bool, Order_chan chan bool) {
	const door_open_time = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-set_timer:
			timer.Reset(door_open_time)
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)

		case <-timer.C:
			timer.Stop()
			timeout <- true
			Order_chan <- false
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
