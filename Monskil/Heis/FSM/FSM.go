package FSM

import (
	"../Driver"
	"../Network"
	"../Timer"
	"fmt"
	"time"
)

//var New_order bool = false

func Function_state_machine() {
	Arrived_chan := make(chan bool)
	Order_chan := make(chan bool)
	Set_timeout_chan := make(chan bool)
	Set_timer_chan := make(chan bool)
	kodd_chan := make(chan bool)
	//New_order_chan := make(chan bool)
	//New_order_print_chan := make(chan bool)
	go Network.Network_server_main( /*New_order*/ )
	go Network.Network_client_main( /*New_order*/ ) //Network.Network_client_main( /*New_order*/ )
	go Order_compare_outer_lists(Order_chan, kodd_chan)
	go Driver.Floor_tracking()
	go Driver.Is_arrived(Arrived_chan, Set_timeout_chan)
	go Driver.Order_set_inner_order()
	go Driver.Order_set_outer_order()
	go Driver.Set_current_floor()
	go Driver.Register_button(Order_chan /*, New_order_chan, New_order_print_chan*/)
	go Timer.Timer(Set_timeout_chan, Set_timer_chan, Order_chan)
	//go Driver.Bursdagskvinn()

	//go Network.Network_client_2_main()
	go Driver.Print_queue()
	for {
		select {

		case <-Arrived_chan:
			//Order_compare_outer_lists(Order_chan, kodd_chan)
			//New_order_print_chan <- true
			//New_order_chan <- true
			//Network.Network_client_main( /*New_order*/ )
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)
			dir := Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
			dir = Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
			Set_timer_chan <- true
			Driver.Elev_set_door_open_lamp(true)
			//Driver.Time_var = int(time.After(3 * time.Second))
		case <-Order_chan:
			//Order_compare_outer_lists()
			//New_order_print_chan <- true
			//New_order_chan <- true
			dir := Driver.Next_order()
			Driver.Elev_set_motor_dir(dir)
		case <-Set_timeout_chan:
			Driver.Elev_set_motor_dir(Driver.DIRN_STOP)
			Driver.Elev_set_door_open_lamp(false)
			//Order_compare_outer_lists(Order_chan, kodd_chan)
		case <-kodd_chan:

		}
	}
}

func Order_compare_outer_lists(Order_chan chan bool, kodd_chan chan bool) {
	for {
		time.Sleep(1 * time.Second)
		counter := 0
		for floor := 0; floor < 4; floor++ {
			if (Driver.Order_outer_list[floor][0] != Network.Server_list[floor][0]) && (Driver.Order_outer_list[floor][0] != 1) {
				//fmt.Println("lol")
				//Driver.Elev_test_set_order_outer_list(floor, 0, Network.Server_list[floor][0], Driver.BUTTON_CALL_UP)
				Driver.Order_outer_list[floor][0] = 0
				Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_UP, floor, 0)
				counter++

			}
			//fmt.Println("hei")
			if (Driver.Order_outer_list[floor][1] != Network.Server_list[floor][1]) && (Driver.Order_outer_list[floor][1] != 1) {
				//Driver.Elev_test_set_order_outer_list(floor, 1, Network.Server_list[floor][1], Driver.BUTTON_CALL_DOWN)
				Driver.Order_outer_list[floor][1] = 0
				Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_DOWN, floor, 0)

				counter++
			} /*else {
				kodd_chan <- true
			}*/
		}
		if counter != 0 {
			Driver.Set_new_order_var()
			//Order_chan <- true
		}
	}
}
func SLETT_DENNE() {
	fmt.Println("SLETT DEN DAA")
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
