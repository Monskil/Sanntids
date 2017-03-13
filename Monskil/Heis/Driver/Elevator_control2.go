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
	"time"
	//"../FSM"
	//"../Network"
)

const N_FLOORS int = 4
const N_BUTTONS int = 3
const MOTOR_SPEED int = 2800

var Current_floor int = 0

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
	Current_floor = 0
	Direction = 0
	fmt.Println("Initialization successfull")
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

func Lights_tracking() {
	for {
		Elev_set_floor_indicator(Elev_get_floor_sensor_signal())

		for floor := 0; floor < N_FLOORS; floor++ {
			if Order_shared_outer_list[floor][0] == 1 {
				Elev_set_button_lamp(BUTTON_CALL_UP, floor, 1)
			} else {
				Elev_set_button_lamp(BUTTON_CALL_UP, floor, 0)
			}
			if Order_shared_outer_list[floor][1] == 1 {
				Elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 1)
			} else {
				Elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
			}
			if Order_inner_list[floor] == 1 {
				Elev_set_button_lamp(BUTTON_COMMAND, floor, 1)
			} else {
				Elev_set_button_lamp(BUTTON_COMMAND, floor, 0)
			}
		}
	}
	time.Sleep(50 * time.Millisecond)
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

var New_order_elev bool = false

func Set_new_order_var() {
	New_order_elev = true
}

func Register_button(Order_chan chan bool /*, New_order_chan chan bool, New_order_print_chan chan bool*/) {

	for {
		for floor := 0; floor < N_FLOORS; floor++ {
			if Check_all_buttons() == Button_channel_matrix[floor][0] {
				/*if floor != Elev_get_floor_sensor_signal() {
					//Elev_set_button_lamp(BUTTON_CALL_UP, floor, 1)
				}*/
				//New_order_print_chan <- true
				//New_order_chan <- true
				//New_order_elev = true
				if IO_read_bit(LIGHT_DOOR_OPEN) == 0 {
					Order_chan <- true
				}

			} else if Check_all_buttons() == Button_channel_matrix[floor][1] {
				/*if floor != Elev_get_floor_sensor_signal() {
					//Elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 1)
				}*/
				//New_order_print_chan <- true
				//New_order_chan <- true
				if IO_read_bit(LIGHT_DOOR_OPEN) == 0 {
					Order_chan <- true
				}
				//New_order_elev = true

			} else if Check_all_buttons() == Button_channel_matrix[floor][2] {
				/*if floor != Elev_get_floor_sensor_signal() {
					//Elev_set_button_lamp(BUTTON_COMMAND, floor, 1)
				}*/
				//New_order_print_chan <- true
				//New_order_chan <- true
				if IO_read_bit(LIGHT_DOOR_OPEN) == 0 {
					Order_chan <- true
				}
				//New_order_elev = true
			}

			if New_order_elev == true {
				if IO_read_bit(LIGHT_DOOR_OPEN) == 0 {
					New_order_elev = false
					Order_chan <- true
				}
			}

		}
		if Get_stop_signal() != 0 {
			Elev_set_stop_lamp(true)
		}
	}
}

/*func Elev_idle_with_order(Order_chan chan bool) {
for floor := 0; floor < N_FLOORS-1; floor++ {
	if Elev_is_idle() == false /* && (Order_outer_list[floor][0] == 1 || Order_outer_list[floor][1] == 1)*/ // {
/*Order_chan <- true
		}
	}
}*/

var Order_inner_list = [N_FLOORS]int{0, 0, 0, 0} //command for floor 1, 2, 3, 4

func Order_set_outer_order() {
	for {
		for floor := 0; floor < N_FLOORS-1; floor++ {
			if Check_all_buttons() == Button_channel_matrix[floor][0] && floor != Elev_get_floor_sensor_signal() {
				Order_shared_outer_list[floor][0] = 1

			}
		}
		for floor := 1; floor < N_FLOORS; floor++ {
			if Check_all_buttons() == Button_channel_matrix[floor][1] && floor != Elev_get_floor_sensor_signal() {
				Order_shared_outer_list[floor][1] = 1

			}
		}
	}
}

func Elev_test_set_order_outer_list(floor int, button int, value int, button_type Elev_button_type_t) {
	//fmt.Println("ja")
	if floor != Elev_get_floor_sensor_signal() {
		Order_shared_outer_list[floor][button] = value
		Elev_set_button_lamp(button_type, floor, value)
	}
}

func Order_set_inner_order() {
	for {
		for floor := 0; floor < N_FLOORS; floor++ {
			if Check_all_buttons() == Button_channel_matrix[floor][2] && floor != Elev_get_floor_sensor_signal() {
				Order_inner_list[floor] = 1
			}
		}
	}
}

var Order_shared_outer_list = [N_FLOORS][N_BUTTONS - 1]int{
	{0, 0},
	{0, 0},
	{0, 0},
	{0, 0},
}

/////////////MANN ER GUL/////////////

var Order_outer_list = [N_FLOORS][N_BUTTONS - 1]int{
	{0, 0 /*FINNES IKKE*/}, //Venstre kolonne er "Opp", høyre er "Ned"
	{0, 0},
	{0, 0},
	{0 /*FINNES IKKE*/, 0},
}
var Direction int = 0

func Elev_is_idle(Order_chan chan bool) bool {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Order_inner_list[floor] == 1 || Order_outer_list[floor][0] == 1 || Order_outer_list[floor][0] == 1 {
			if IO_read_bit(LIGHT_DOOR_OPEN) == 0 {
				Order_chan <- true
			}
			return false
		}
	}
	return true
}

func Next_order() Elev_motor_direction_t {
	//Sjekke bestillinger over
	//fmt.Println(Current_floor)
	//fmt.Println(Direction)
	var More_orders_up bool = false
	for floor := Current_floor + 1; floor < N_FLOORS; floor++ {
		if Order_outer_list[floor][0] == 1 || Order_outer_list[floor][1] == 1 || Order_inner_list[floor] == 1 {
			More_orders_up = true
		}
	}
	//sjekke bestillinger under
	var More_orders_down bool = false
	for floor := Current_floor - 1; floor >= 0; floor-- {
		if Order_outer_list[floor][0] == 1 || Order_outer_list[floor][1] == 1 || Order_inner_list[floor] == 1 {
			More_orders_down = true
		}
	}
	for floor := Current_floor - 1; floor >= 0; floor-- {
		if More_orders_down == false {
			if Direction != 1 { /*HER*/
				Direction = 1 /*HER*/
			}
		}
		if Current_floor == 0 {
			if Direction != 1 { /*HER*/
				Direction = 1 /*HER*/
			}
			IO_clear_bit(MOTORDIR)
		}
		if ((Order_inner_list[floor] == 1) || (Order_outer_list[floor][0] == 1) || (Order_outer_list[floor][1] == 1)) && (floor < Current_floor) && (Direction != 1) && (More_orders_down == true /*counter_down_inner != current_floor*/) {
			//fmt.Println("feil")
			if Direction != -1 {
				Direction = -1
			}
			return DIRN_DOWN
		}
	}
	for floor := Current_floor + 1; floor < N_FLOORS; floor++ {
		if More_orders_up == false {
			if Direction != -1 { /*HER*/
				Direction = -1 /*HER*/
			}
		}
		if Current_floor == 3 {
			if Direction != -1 { /*HER*/
				Direction = -1 /*HER*/
			}
			IO_set_bit(MOTORDIR)
		}
		if ((Order_inner_list[floor] == 1) || (Order_outer_list[floor][0] == 1) || (Order_outer_list[floor][1] == 1)) && (floor > Current_floor) && (Direction != -1) && (More_orders_up == true /*counter_up_inner != N_FLOORS-current_floor*/) {

			if Direction != 1 {
				Direction = 1
			}
			return DIRN_UP
		}
	}
	Direction = 0
	return DIRN_STOP
}

func Is_arrived(Arrived_chan chan bool, Set_timeout_chan chan bool) {
	for {
		for floor := 0; floor < N_FLOORS; floor++ {
			if IO_read_bit(MOTORDIR) == 0 && Order_inner_list[floor] == 1 && floor == Elev_get_floor_sensor_signal() {
				//Elev_set_button_lamp(BUTTON_COMMAND, floor, 0)
				Arrived_chan <- true
				select {
				case <-Set_timeout_chan:
					Order_inner_list[floor] = 0
					Elev_set_door_open_lamp(false)
				}

			}
			if IO_read_bit(MOTORDIR) == 0 && Order_outer_list[floor][0] == 1 && floor == Elev_get_floor_sensor_signal() {
				//Elev_set_button_lamp(BUTTON_CALL_UP, floor, 0)
				Arrived_chan <- true
				select {
				case <-Set_timeout_chan:
					Order_outer_list[floor][0] = 0
					Order_shared_outer_list[floor][0] = 0
					Elev_set_door_open_lamp(false)
				}
			}
			if IO_read_bit(MOTORDIR) == 0 && Order_outer_list[floor][1] == 1 && floor == Elev_get_floor_sensor_signal() {
				if floor == 3 {
					//Elev_set_button_lamp(BUTTON_CALL_DOWN, 3, 0)
					Arrived_chan <- true
					select {
					case <-Set_timeout_chan:
						Order_outer_list[3][1] = 0
						Order_shared_outer_list[3][1] = 0
						Elev_set_door_open_lamp(false)
					}
				} else if floor == 2 && Order_outer_list[3][1] == 0 {
					//Elev_set_button_lamp(BUTTON_CALL_DOWN, 2, 0)
					Arrived_chan <- true
					select {
					case <-Set_timeout_chan:
						Order_outer_list[2][1] = 0
						Order_shared_outer_list[2][1] = 0
						Elev_set_door_open_lamp(false)
					}

				} else if floor == 1 && Order_outer_list[3][1] == 0 && Order_outer_list[2][1] == 0 {
					//Elev_set_button_lamp(BUTTON_CALL_DOWN, 1, 0)
					Arrived_chan <- true
					select {
					case <-Set_timeout_chan:
						Order_outer_list[1][1] = 0
						Order_shared_outer_list[1][1] = 0
						Elev_set_door_open_lamp(false)
					}

				}
			}
		}
		for floor := N_FLOORS - 1; floor >= 0; floor-- {
			if IO_read_bit(MOTORDIR) == 1 && Order_inner_list[floor] == 1 && floor == Elev_get_floor_sensor_signal() {
				//Elev_set_button_lamp(BUTTON_COMMAND, floor, 0)
				Arrived_chan <- true
				select {
				case <-Set_timeout_chan:
					Order_inner_list[floor] = 0
					Elev_set_door_open_lamp(false)
				}
			}
			if IO_read_bit(MOTORDIR) == 1 && Order_outer_list[floor][1] == 1 && floor == Elev_get_floor_sensor_signal() {
				//Elev_set_button_lamp(BUTTON_CALL_DOWN, floor, 0)
				Arrived_chan <- true
				select {
				case <-Set_timeout_chan:
					Order_outer_list[floor][1] = 0
					Order_shared_outer_list[floor][1] = 0
					Elev_set_door_open_lamp(false)
				}
			}
			if IO_read_bit(MOTORDIR) == 1 && Order_outer_list[floor][0] == 1 && floor == Elev_get_floor_sensor_signal() {
				if floor == 0 {
					//Elev_set_button_lamp(BUTTON_CALL_UP, 0, 0)
					Arrived_chan <- true
					select {
					case <-Set_timeout_chan:
						Order_outer_list[0][0] = 0
						Order_shared_outer_list[0][0] = 0
						Elev_set_door_open_lamp(false)
					}

				} else if floor == 1 && Order_outer_list[0][0] == 0 {
					//Elev_set_button_lamp(BUTTON_CALL_UP, 1, 0)
					Arrived_chan <- true
					select {
					case <-Set_timeout_chan:
						Order_outer_list[1][0] = 0
						Order_shared_outer_list[1][0] = 0
						Elev_set_door_open_lamp(false)
					}

				} else if floor == 2 && Order_outer_list[1][0] == 0 && Order_outer_list[0][0] == 0 {
					//Elev_set_button_lamp(BUTTON_CALL_UP, 2, 0)
					Arrived_chan <- true
					select {
					case <-Set_timeout_chan:
						Order_outer_list[2][0] = 0
						Order_shared_outer_list[2][0] = 0
						Elev_set_door_open_lamp(false)
					}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

/////////////GULT ER KULT////////////

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

func Set_current_floor() {
	for {
		temp := Elev_get_floor_sensor_signal()
		if temp != -1 {
			Current_floor = temp
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func Order_handling(floor int) {

	for {
		//var Current_floor int = 0 //trenger vi å sette den her? den står uttafor funksjoner i denne filen også
		if Elev_get_floor_sensor_signal() != -1 {
			Current_floor = Elev_get_floor_sensor_signal()
		}
		if Current_floor < floor {
			Elev_set_motor_dir(DIRN_UP)
			return
		} else if Current_floor > floor {
			Elev_set_motor_dir(DIRN_DOWN)
			return
		}
	}
}

func Print_queue() {
	for {
		fmt.Println("Shared:     ", Order_shared_outer_list)
		fmt.Println("Egen outer :", Order_outer_list)
		time.Sleep(1 * time.Second)
	}
}
