package main

import "./Driver"
import "fmt"

func main() {
	Driver.Elev_init()
	for{
		
		//fmt.Println(Driver.Order_outer_list)
		//fmt.Println(Driver.Order_inner_list)
		for Driver.Elev_get_floor_sensor_signal() != 0{
			Driver.Elev_set_motor_dir(Driver.DIRN_DOWN)
			Driver.Floor_tracking()
			Driver.Register_button()
			Driver.Order_set_outer_order()
			Driver.Order_set_inner_order()
			Driver.JUNIORRRR_aka_Order_complete()
			fmt.Println(Driver.Order_inner_list)
		}
		for Driver.Elev_get_floor_sensor_signal() != 3{
			Driver.Elev_set_motor_dir(Driver.DIRN_UP)
			Driver.Floor_tracking()
			Driver.Register_button()
			Driver.Order_set_outer_order()
			Driver.Order_set_inner_order()
			Driver.JUNIORRRR_aka_Order_complete()
			fmt.Println(Driver.Order_inner_list)
		}
	}
}