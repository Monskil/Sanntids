package Driver // where "driver" is the folder that contains IO.go, IO.c, IO.h, Channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "IO2.h"
*/
//import "C"

import (
	//"../Timer"
	"fmt"
	//"runtime"
	//"time"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3
const MOTOR_SPEED int = 2800

var current_floor int = -1

type Elev_button_type_t int

const (
	BUTTON_CALL_UP   = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND   = 2
)

type Elev_motor_direction_t int

const (
	DIRN_DOWN = -1
	DIRN_STOP = 0
	DIRN_UP   = 1
)

func Elev_init() (int, error) {
	var init_success int = IO_init()

	if init_success == 0 {
		return -1, fmt.Errorf("Unable to initialize elevator hardware")
	}

	for floor := 0; floor < N_FLOORS; floor++ {
		//var b Elev_button_type_t
		//for b := 0; b < N_BUTTONS; b++ {
		//Elev_set_button_lamp(b,f,0)
		//}
		if floor != 0 {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
		}
		if floor != N_FLOORS-1 {
			Elev_set_button_lamp(BUTTON_CALL_UP, floor, 0)
		}
		Elev_set_button_lamp(BUTTON_COMMAND, floor, 0)
	}
	Elev_set_stop_lamp(false)
	Elev_set_door_open_lamp(false)
	for Elev_get_floor_sensor_signal() != 0 {
		Elev_set_motor_dir(DIRN_DOWN)
	}
	Elev_set_motor_dir(DIRN_STOP)
	Elev_set_floor_indicator(0)
	current_floor = 0
	return 0, nil
}

func Elev_set_motor_dir(dir_n Elev_motor_direction_t) {
	if dir_n == 0 {
		IO_write_analog(MOTOR, 0)
	} else if dir_n > 0 {
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	} else if dir_n < 0 {
		IO_set_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func Elev_set_button_lamp(button Elev_button_type_t, floor int, value int) (int, error) {

	if floor < 0 || floor > N_FLOORS /*|| button < 0 || button > N_BUTTONS*/ {
		return -1, fmt.Errorf("Floor has an illegal value")
	}
	if button != BUTTON_CALL_UP && button != BUTTON_CALL_DOWN && button != BUTTON_COMMAND {
		return -1, fmt.Errorf("Button has an illegal value")
	}
	if button == BUTTON_CALL_UP && floor == N_FLOORS-1 {
		return -1, fmt.Errorf("Button up from top floor does not exist")
	}
	if button == BUTTON_CALL_DOWN && floor == 0 {
		return -1, fmt.Errorf("Button down from ground floor does not exist")
	}
	if value != 0 {
		IO_set_bit(Lamp_channel_matrix[floor][button])
	} else {
		IO_clear_bit(Lamp_channel_matrix[floor][button])
	}
	return 0, nil
}

func Elev_set_floor_indicator(floor int) (int, error) {

	// Binary encoding. One light must always be on.

	if floor < 0 || floor > N_FLOORS {
		return -1, fmt.Errorf("Floor has an illegal value")
	}
	if floor&0x02 != 0 {
		IO_set_bit(LIGHT_FLOOR_IND1)
	} else {
		IO_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor&0x01 != 0 {
		IO_set_bit(LIGHT_FLOOR_IND2)
	} else {
		IO_clear_bit(LIGHT_FLOOR_IND2)
	}
	return 0, nil
}

func Elev_set_door_open_lamp(value bool) {
	if value {
		IO_set_bit(LIGHT_DOOR_OPEN)
	} else {
		IO_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_set_stop_lamp(value bool) {
	if value {
		IO_set_bit(LIGHT_STOP)
	} else {
		IO_clear_bit(LIGHT_STOP)
	}
}

func Elev_get_button_signal(button Elev_button_type_t, floor int) (int, error, int) {

	if floor < 0 || floor > N_FLOORS /* || button < 0 || button > N_BUTTONS*/ {
		return -1, fmt.Errorf("Floor has an illegal value"), 0
	}
	if button != BUTTON_CALL_UP && button != BUTTON_CALL_DOWN && button != BUTTON_COMMAND {
		return -1, fmt.Errorf("Button has an illegal value"), 0
	}
	if button == BUTTON_CALL_UP && floor == N_FLOORS-1 {
		return -1, fmt.Errorf("Button up from top floor does not exist"), 0
	}
	if button == BUTTON_CALL_DOWN && floor == 0 {
		return -1, fmt.Errorf("Button down from ground floor does not exist"), 0
	}
	return 0, nil, IO_read_bit(Button_channel_matrix[floor][button])
}

func Elev_get_floor_sensor_signal() int {
	if IO_read_bit(SENSOR_FLOOR1) != 0 {
		return 0
	} else if IO_read_bit(SENSOR_FLOOR2) != 0 {
		return 1
	} else if IO_read_bit(SENSOR_FLOOR3) != 0 {
		return 2
	} else if IO_read_bit(SENSOR_FLOOR4) != 0 {
		return 3
	} else {
		return -1
	}
}

func Get_stop_signal() int {
	return IO_read_bit(STOP)
}

func Get_obstruction_signal() int {
	return IO_read_bit(OBSTRUCTION)
}

func Floor_tracking(at_floor_chan chan int) {
	Elev_set_floor_indicator(Elev_get_floor_sensor_signal())
	at_floor_chan <- Elev_get_floor_sensor_signal()
}

var Lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var Button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func Check_all_buttons() int {
	for {
		for floor := 0; floor < N_FLOORS; floor++ {
			if _, _, x := Elev_get_button_signal(BUTTON_CALL_UP, floor); x == 1 {
				return Button_channel_matrix[floor][0]
			} else if _, _, x := Elev_get_button_signal(BUTTON_CALL_DOWN, floor); x == 1 {
				return Button_channel_matrix[floor][1]
			} else if _, _, x := Elev_get_button_signal(BUTTON_COMMAND, floor); x == 1 {
				return Button_channel_matrix[floor][2]
			}
		}
	}
	return -1
}

func Go_to_buttons(Button_chan chan bool) {
	for {
		if Check_all_buttons() != -1 {
			Button_chan <- true
		}
	}
}

func Register_button() {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Check_all_buttons() == Button_channel_matrix[floor][0] {
			Elev_set_button_lamp(BUTTON_CALL_UP, floor, 1)
		} else if Check_all_buttons() == Button_channel_matrix[floor][1] {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 1)
		} else if Check_all_buttons() == Button_channel_matrix[floor][2] {
			Elev_set_button_lamp(BUTTON_COMMAND, floor, 1)
		}
	}
	if Get_stop_signal() != 0 {
		Elev_set_stop_lamp(true)
	}
}

func JUNIORRRR_aka_Order_complete(Floor_chan chan bool, Button_chan chan bool) {
	fmt.Println("j")
	for floor := 0; floor < N_FLOORS; floor++ {
		if Order_inner_list[floor] == 1 && Elev_get_floor_sensor_signal() == floor /*dooropen*/ {
			Elev_set_button_lamp(BUTTON_COMMAND, floor, 0)
			Order_inner_list[floor] = 0
			Floor_chan <- true
		}
		if Order_outer_list[floor][0] == 1 && Elev_get_floor_sensor_signal() == floor && IO_read_bit(MOTORDIR) == 0 /*DOOROPEN*/ {
			Elev_set_button_lamp(BUTTON_CALL_UP, floor, 0)
			Order_outer_list[floor][0] = 0
			Floor_chan <- true

		} else if Order_outer_list[floor][1] == 1 && Elev_get_floor_sensor_signal() == floor && IO_read_bit(MOTORDIR) == 1 /*DOOROPEN*/ {
			Elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
			Order_outer_list[floor][1] = 0
			Floor_chan <- true
		}
	}
}

var Order_outer_list = [N_FLOORS][N_BUTTONS - 1]int{
	{0, 0}, //Venstre kolonne er "Opp", høyre er "Ned"
	{0, 0},
	{0, 0},
	{0, 0},
}

var Order_inner_list = [4]int{0, 0, 0, 0} //command for floor 1, 2, 3, 4

func Order_set_outer_order() {
	for floor := 0; floor < N_FLOORS-1; floor++ {
		if Check_all_buttons() == Button_channel_matrix[floor][0] {
			Order_outer_list[floor][0] = 1
		}
	}
	for floor := 1; floor < N_FLOORS; floor++ {
		if Check_all_buttons() == Button_channel_matrix[floor][1] {
			Order_outer_list[floor][1] = 1
		}
	}
}

func Order_set_inner_order() {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Check_all_buttons() == Button_channel_matrix[floor][2] {
			Order_inner_list[floor] = 1
		}
	}
}

func Test_elev(test_chan chan bool) {
	if Get_obstruction_signal() == 1 && Elev_get_floor_sensor_signal() != 0 {
		Elev_set_motor_dir(DIRN_DOWN)
	} else if Get_obstruction_signal() == 0 && Elev_get_floor_sensor_signal() != 3 {
		Elev_set_motor_dir(DIRN_UP)
	} else {
		Elev_set_motor_dir(DIRN_STOP)
	}
}

func Elev_is_elevator_vacant() bool {
	var x = 0
	for floor := 0; floor < N_FLOORS; floor++ {
		if Order_inner_list[floor] == 1 {
			x++
		}
	}
	if x != 0 {
		return false
	}
	return true
}

func Go_to_order(order_chan chan bool) {
	for {
		if Elev_is_elevator_vacant() == false {
			order_chan <- true
		}
	}
}

func Elev_is_outer_orders() int {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Order_outer_list[floor][0] == 1 && Elev_get_floor_sensor_signal() == floor && IO_read_bit(MOTORDIR) == 0 /*DOOROPEN*/ {
			return 1 //Opp-ordre
		} else if Order_outer_list[floor][1] == 1 && Elev_get_floor_sensor_signal() == floor && IO_read_bit(MOTORDIR) == 1 /*DOOROPEN*/ {
			return 2 //Ned-ordre
		}
	}
	return 0 //Ingen ordre
}

func Begin_journey(floor int) {
	if current_floor < floor {
		Elev_set_motor_dir(DIRN_UP)
	} else if current_floor > floor {
		Elev_set_motor_dir(DIRN_DOWN)
	}
}

func Set_current_floor() {
	if Elev_get_floor_sensor_signal() != -1 {
		current_floor = Elev_get_floor_sensor_signal()
	}
}

func Order_handling() {
	for i := 0; i < N_FLOORS; i++ {
		if Order_inner_list[i] == 1 {
			Begin_journey(i)
		}
	}
}

/*
func Manage_door() {
	if Elev_get_floor_sensor_signal() != -1 && IO_read_bit(MOTOR) == 0 {
		timeout := make(chan bool)
		set_timer := make(chan bool)
		go Timer.Timer(timeout, set_timer)
		Elev_set_door_open_lamp(true)
		is_timeout := make(chan bool)
		for {
			select {
			case <-timeout:
				Elev_set_door_open_lamp(false)
				is_timeout <- false
			case <-is_timeout:
				return
			}
		}
	}
}
*/
