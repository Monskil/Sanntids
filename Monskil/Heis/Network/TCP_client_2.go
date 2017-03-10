package main

//package Network

import (
	//"../Driver"
	"fmt"
	"net"
	"os"
	//"time"
)

func Orders_to_string_1() string {

	test_inner := [4]int{1, 1, 0, 1}
	test_outer := [4][2]int{
		{0, 0},
		{1, 0},
		{0, 1},
		{0, 0},
	}

	var Orders string = "" //UUUUDDDDCCCC (U = orders button_up | D = orders button_down | C = orders button_command)
	for floor := 0; floor < 4; /*Driver.N_FLOORS*/ floor++ {
		if /*Driver.Order_outer_list[floor][0]*/ test_outer[floor][0] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < 4; /*Driver.N_FLOORS*/ floor++ {
		if /*Driver.Order_outer_list[floor][1]*/ test_outer[floor][1] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < 4; /*Driver.N_FLOORS*/ floor++ {

		if /*Driver.Order_inner_list[floor] */ test_inner[floor] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	return Orders
}

func Orders_to_string_2() string {

	test_inner := [4]int{1, 1, 1, 1}
	test_outer := [4][2]int{
		{1, 1},
		{1, 1},
		{0, 0},
		{0, 0},
	}

	var Orders string = "" //UUUUDDDDCCCC (U = orders button_up | D = orders button_down | C = orders button_command)
	for floor := 0; floor < 4; /*Driver.N_FLOORS*/ floor++ {
		if /*Driver.Order_outer_list[floor][0]*/ test_outer[floor][0] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < 4; /*Driver.N_FLOORS*/ floor++ {
		if /*Driver.Order_outer_list[floor][1]*/ test_outer[floor][1] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < 4; /*Driver.N_FLOORS*/ floor++ {

		if /*Driver.Order_inner_list[floor] */ test_inner[floor] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	return Orders

}

func main() {
	// connect to this socket
	conn, _ := net.Dial("tcp", "localhost:1201")
	for {
		fmt.Println(Orders_to_string_2())
		fmt.Fprintf(conn, Orders_to_string_2())
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
