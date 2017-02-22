package Driver // where "driver" is the folder that contains IO.go, IO.c, IO.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "IO2.h"
*/
//import "C"

import (
	. "fmt"
	"runtime"
	"time"
)

type Elev_button_type_t int 
const ( 
	BUTTON_CALL_UP = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_COMMAND = 2
)

type Elev_motor_direction_t int
const ( 
	DIRN_DOWN = -1
	DIRN_STOP = 0
	DIRN_UP = 1
)


func Elev_init() (int, error){
	var init_success int = IO_init()  
	if init_success = 0{											
		return -1, fmt.Errorf("Unable to initialize elevator hardware")
	}														

	for f int := 0; f < N_FLOORS; f++ {
		for b Elev_button_type_t := 0; b < N_BUTTONS; b++ {
			Elev_set_button_lamp(b,f,0)
		}
	}
	Elev_set_stop_lamp(0)
	Elev_set_door_open_lamp(0)
	Elev_set_floor_indicator(0)
	
	return nil,_
}

func Elev_set_motor_dir(dir_n Elev_motor_direction_t){
	if dir_n = 0{
		IO_write_analog(MOTOR , 0)
	} else if dir_n > 0{
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	} else if dir_n < 0{
		IO_set_bit(MOTORDIR)
		IO_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func Elev_set_button_lamp(button elev_button_type_t, floor int, value int) (int, error) {

    if floor < 0 || floor > N_FLOORS || button < 0 || button > N_BUTTONS{
    	return -1, fmt.Errorf("Floor or button has an illegal value")
    if value{
    	IO_set_bit(lamp_channel_matrix[floor][button])
    } else{
    	IO_clear_bit(lamp_channel_matrix[floor][button])
    }
    return _,nil
}

func Elev_set_floor_indicator(floor int) (int, error){
	
    // Binary encoding. One light must always be on.

	if floor < 0 || floor > N_FLOORS{
		return -1, fmt.Errorf("Floor has an illegal value")
	}
	if floor & 0x02{
		IO_set_bit(LIGHT_FLOOR_IND1)
	} else{
		IO_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor & 0x01 != 0{
		IO_set_bit(LIGHT_FLOOR_IND2)
	} else{
		IO_clear_bit(LIGHT_FLOOR_IND2)
	}
	return _, nil
}

func Elev_set_door_open_lamp(value int){
	if value{
		IO_set_bit(LIGHT_DOOR_OPEN)
	} else{
		IO_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_set_stop_lamp(value int){
	if value{
		IO_set_bit(LIGHT_STOP)
	} else{
		IO_clear_bit(LIGHT_STOP)
	}
}

func Elev_get_button_signal(button elev_button_type_t, floor int) int, error, int{
	
	if floor < 0 || floor > N_FLOORS || button < 0 || button > N_BUTTONS{
    	return -1, fmt.Errorf("Floor or button has an illegal value"),_

    return _, nil, IO_read_bit(button_channel_matrix[floor][button])
}

func Elev_get_floor_sensor_signal() int{
	if IO_read_bit(SENSOR_FLOOR1){
		return 1
	} else if IO_read_bit(SENSOR_FLOOR2){
		return 2
	} else if IO_read_bit(SENSOR_FLOOR3){
		return 3
	} else if IO_read_bit(SENSOR_FLOOR4){
		return 4
	} else{
		return -1
	}
}

func Get_stop_signal() int{
	return IO_read_bit(STOP)
}

func Get_obstructIOn_signal() int{
	return IO_read_bit(OBSTRUCTION)
}

var Lamp_channel_matrix := [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1,LIGHT_DOWN1,LIGHT_COMMAND1},
	{LIGHT_UP2,LIGHT_DOWN2,LIGHT_COMMAND2},
	{LIGHT_UP3,LIGHT_DOWN3,LIGHT_COMMAND3},
	{LIGHT_UP4,LIGHT_DOWN4,LIGHT_COMMAND4}
}

var Button_channel_matrix := [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1,BUTTON_DOWN1,BUTTON_COMMAND1},
	{BUTTON_UP2,BUTTON_DOWN2,BUTTON_COMMAND2},
	{BUTTON_UP3,BUTTON_DOWN3,BUTTON_COMMAND3},
	{BUTTON_UP4,BUTTON_DOWN4,BUTTON_COMMAND4}
}


