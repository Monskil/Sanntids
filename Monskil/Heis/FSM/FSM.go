package FSM

import (
	"../Driver"
	//"fmt"
	//"time"
)

func Function_state_machine() {
	Arrived_chan := make(chan bool)
	Order_chan := make(chan bool)
	go Driver.Floor_tracking()
	go Driver.Order_set_inner_order()
	go Driver.Set_current_floor()
	go Driver.Register_button(Order_chan)
	go Driver.Is_arrived(Arrived_chan)
	//go Driver.Print_queue()
	for {
		select {

		case <-Arrived_chan:
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)
			dir := Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
			dir = Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
		case <-Order_chan:
			dir := Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)

		}
	}
}

/*
	//at_floor_chan := make(chan int)
	order_chan := make(chan int)
	Button_chan := make(chan bool)
	Floor_chan := make(chan bool)
	//	go Driver.Check_all_buttons()
	go Driver.Go_to_buttons(Button_chan)
	go Driver.Go_to_order(order_chan)
	go Driver.Floor_tracking()
	go Driver.Register_button()
	go Driver.Set_current_floor()
	//go Driver.JUNIORRRR_aka_Order_complete(Floor_chan)

	for {

		//go Driver.Register_button(Button_chan)
		select {

		case <-Button_chan: //Hvis knappetrykk
			Driver.Order_set_outer_order()
			Driver.Order_set_inner_order()

		case floor := <-order_chan: //Hvis bestilling
			Driver.Order_handling(floor)
			Driver.JUNIORRRR_aka_Order_complete(floor, Floor_chan)
		case <-Floor_chan: //Hvis ankommet bestilling
			fmt.Println("MANNERGUL")
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
