package FSM

import (
	"../Driver"
	"../Network"
	"../Timer"
	"fmt"
	"time"
)

func Function_state_machine() {
	Arrived_chan := make(chan bool, 100000000)
	Order_chan := make(chan bool, 100000000)
	Set_timeout_chan := make(chan bool, 100000000)
	Set_timer_chan := make(chan bool, 100000000)

	go Network.Network_server_main()
	go Network.Network_server_main_2()
	go Network.Network_client_main()
	go Network.Network_client_main_2()
	go Order_compare_outer_lists(Order_chan)
	go Driver.Lights_tracking()
	go Driver.Is_arrived(Arrived_chan, Set_timeout_chan)
	go Driver.Order_set_inner_order()
	go Driver.Order_set_outer_order()
	go Driver.Set_current_floor()
	go Driver.Register_button(Order_chan)
	go Timer.Timer(Set_timeout_chan, Set_timer_chan, Order_chan)
	//go Driver.Bursdagskvinn()

	//go Driver.Print_queue()
	for {
		select {

		case <-Arrived_chan:
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)
			dir := Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
			dir = Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
			Set_timer_chan <- true
			Driver.Elev_set_door_open_lamp(true)
		case <-Order_chan:
			dir := Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
		case <-Set_timeout_chan:
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)
			Driver.Elev_set_door_open_lamp(false)
		}
	}
}

func Order_compare_outer_lists(Order_chan chan bool) {
	for {
		time.Sleep(1 * time.Second)
		counter := 0
		for floor := 0; floor < 4; floor++ {
			if (Driver.Order_outer_list[floor][0] != Network.Server_list[floor][0]) && (Driver.Order_outer_list[floor][0] != 1) {
				Driver.Order_outer_list[floor][0] = Network.Server_list[floor][0]
				Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_UP, floor, 0)
				counter++

			}
			if (Driver.Order_outer_list[floor][1] != Network.Server_list[floor][1]) && (Driver.Order_outer_list[floor][1] != 1) {
				Driver.Order_outer_list[floor][1] = Network.Server_list[floor][1]
				Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_DOWN, floor, 0)
				counter++
			}
		}
		if counter != 0 {
			Driver.Set_new_order_var()
		}
	}
}
func SLETT_DENNE() {
	fmt.Println("SLETT DEN DAA")
	time.Sleep(1 * time.Second)
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
