package Timer

import (
	"../Driver"
	//"../Network_main"
	//"fmt"
	"time"
)

func Timer(timeout chan bool, set_timer chan bool, Order_chan chan bool) {
	const door_open_time = 3 * time.Second //When the door opens, it will stay that way for 3 seconds
	timer := time.NewTimer(0)
	timer.Stop()

	for {
		select {
		case <-set_timer: //Sets timer when arrived at an order
			timer.Reset(door_open_time)
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)

		case <-timer.C:
			timer.Stop()
			timeout <- true //Times out after 3 seconds
			Order_chan <- true
		}
	}
}

/*
func Timer_2(timeout_2 chan bool, Set_timer_2_floor chan int, Set_timer_2_b chan int) {
	const dead_time = 2 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	//fmt.Println("in timer 2")
	for {
		select {
		case floor := <-Set_timer_2_floor:
			b := <-Set_timer_2_b
			if Driver.Order_outer_list[floor][b] == 0 {
				fmt.Println("order ok")
				timer.Reset(dead_time)
			}
		case <-timer.C:
			fmt.Println("order not ok")
			timer.Stop()
			Network_main.Dead_1 = true
			timeout_2 <- true
		}
	}
}*/
