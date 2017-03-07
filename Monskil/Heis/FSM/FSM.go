package FSM

import (
	"../Driver"
	"fmt"
)

func Function_state_machine() {

	//at_floor_chan := make(chan int)
	order_chan := make(chan bool)
	Button_chan := make(chan bool)
	Floor_chan := make(chan bool)
	go Driver.Check_all_buttons()
	go Driver.Go_to_buttons(Button_chan)
	go Driver.Go_to_order(order_chan)
	go Driver.JUNIORRRR_aka_Order_complete(Floor_chan, Button_chan)
	for {

		//go Driver.Register_button(Button_chan)
		select {

		case <-Button_chan:
			fmt.Println(Driver.Order_inner_list)
			Driver.Set_current_floor()
			Driver.Order_set_outer_order()
			Driver.Order_set_inner_order()
			Driver.Register_button()
			Driver.JUNIORRRR_aka_Order_complete(Floor_chan, Button_chan)

		case <-order_chan:
			Driver.Order_handling()

		case <-Floor_chan:
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)
			fmt.Println("JUNIOORRR")
		}
	}
}

/*
Driver.Order_set_outer_order()
Driver.Order_set_inner_order()
Driver.JUNIORRRR_aka_Order_complete()
Driver.Floor_tracking(at_floor_chan)

Driver.Floor_tracking(at_floor_chan)
*/
