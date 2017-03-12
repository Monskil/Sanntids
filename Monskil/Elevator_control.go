package driver // where "driver" is the folder that contains io.go, io.c, io.h, channels.go, channels.h and driver.go
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
//import "C"

import (
	. "fmt"
	"runtime"
	"time"
)

func elev_init() int, error{
	var init_success int = io_init()
	
	if !init_success{											
		return -1, fmt.Errorf("Unable to initialize elevator hardware")											
	}														

	for f := 0; f < N_FLOORS; f++ {
		for b elev_button_type := 0; b < N_BUTTONS; b++{
			elev_set_button_lamp(b,f,0)
		}
	}
	elev_set_stop_lamp(0)
	elev_set_door_open_lamp(0)
	elev_set_floor_indicator(0)
	
	return nil,_
}

func elev_set_motor_dir(dir_n elev_motor_direction_t){
	if dir_n = 0{
		io_write_analog(MOTOR , 0)
	}
	else if dir_n > 0{
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	}
	else if dir_n < 0{
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func elev_set_button_lamp(button elev_button_type_t, floor int, value int) int, error {

    if floor < 0 || floor > N_FLOORS || button < 0 || button > N_BUTTONS{
    	return -1, fmt.Errorf("Floor or button has an illegal value")
    if value{
    	io_set_bit(lamp_channel_matrix[floor][button])
    }
    else{
    	io_clear_bit(lamp_channel_matrix[floor][button])
    }
    return _,nil
}

func elev_set_floor_indicator(floor int) int, error{
	
    // Binary encoding. One light must always be on.

	if floor < 0 || floor > N_FLOORS{
		return -1, fmt.Errorf("Floor has an illegal value")
	}
	if floor & 0x02{
		io_set_bit(LIGHT_FLOOR_IND1)
	}
	else{
		io_clear_bit(LIGHT_FLOOR_IND1)
	}
	if floor & 0x01{
		io_set_bit(LIGHT_FLOOR_IND2)
	}
	else{
		io_clear_bit(LIGHT_FLOOR_IND2)
	}
	return _, nil
}

func elev_set_door_open_lamp(value int){
	if value{
		io_set_bit(LIGHT_DOOR_OPEN)
	}
	else{
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func elev_set_stop_lamp(value int){
	if value{
		io_set_bit(LIGHT_STOP)
	}
	else{
		io_clear_bit(LIGHT_STOP)
	}
}

func elev_get_button_signal(button elev_button_type_t, floor int) int, error, int{
	
	if floor < 0 || floor > N_FLOORS || button < 0 || button > N_BUTTONS{
    	return -1, fmt.Errorf("Floor or button has an illegal value"),_

    return _, nil, io_read_bit(button_channel_matrix[floor][button])
}

func elev_get_floor_sensor_signal() int{
	if io_read_bit(SENSOR_FLOOR1){
		return 1
	}
	else if io_read_bit(SENSOR_FLOOR2){
		return 2
	}
	else if io_read_bit(SENSOR_FLOOR3){
		return 3
	}
	else if io_read_bit(SENSOR_FLOOR4){
		return 4
	}
	else{
		return -1
	}
}

func get_stop_signal() int{
	return io_read_bit(STOP)
}

func get_obstruction_signal() int{
	return io_read_bit(OBSTRUCTION)
}

var lamp_channel_matrix := [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1,LIGHT_DOWN1,LIGHT_COMMAND1},
	{LIGHT_UP2,LIGHT_DOWN2,LIGHT_COMMAND2},
	{LIGHT_UP3,LIGHT_DOWN3,LIGHT_COMMAND3},
	{LIGHT_UP4,LIGHT_DOWN4,LIGHT_COMMAND4}
}

var button_channel_matrix := [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1,BUTTON_DOWN1,BUTTON_COMMAND1},
	{BUTTON_UP2,BUTTON_DOWN2,BUTTON_COMMAND2},
	{BUTTON_UP3,BUTTON_DOWN3,BUTTON_COMMAND3},
	{BUTTON_UP4,BUTTON_DOWN4,BUTTON_COMMAND4}
}


