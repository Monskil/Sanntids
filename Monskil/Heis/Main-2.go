package main

import (
	//"./Timer"
	"./Driver"
	"./FSM"
	//"fmt"
)

func main() {
	Driver.Elev_init()
	FSM.Function_state_machine()

	/*for {
		//Driver.Elev_get_floor_sensor_signal() != 0{
		//Driver.Elev_set_motor_dir(Driver.DIRN_DOWN)
		//Driver.Test_elev()
		//Driver.Floor_tracking(at_floor_chan)
		Driver.Register_button()
		Driver.Order_set_outer_order()
		Driver.Order_set_inner_order()
		Driver.JUNIORRRR_aka_Order_complete()

		//fmt.Println(Driver.Order_outer_list)

		//fmt.Println(Driver.Order_outer_list)
		//fmt.Println(Driver.Order_inner_list)
		FSM.Function_state_machine()
	}
	*/
}
