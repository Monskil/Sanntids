package main

import (
	//"./Timer"
	"./Driver"
	"./FSM"
	//"./Network"
	//"fmt"
	//"strings"
)

func main() {

	Driver.Elev_init()
	FSM.Function_state_machine()

	//for {
	//var teststreng string = Orders_to_string()
	//	var x string = Orders_to_string()
	//fmt.Println("x: ", x[:2]+"k")
	//Driver.Order_set_outer_order()
	//Driver.Order_set_inner_order()
	//	fmt.Println(x)
	//}
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
