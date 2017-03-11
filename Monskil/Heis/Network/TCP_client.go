package Network

//package Network

import (
	"../Driver"
	"fmt"
	"net"
	//"os"
	//"time"
)

func Orders_to_string_1() string {

	/*test_inner := [4]int{0, 0, 0, 0}
	test_outer := [4][2]int{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 0},
	}
	*/
	var Orders string = "" //UUUUDDDDCCCC (U = orders button_up | D = orders button_down | C = orders button_command)
	for floor := 0; floor < Driver.N_FLOORS; floor++ {
		if Driver.Order_outer_list[floor][0] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < Driver.N_FLOORS; floor++ {
		if Driver.Order_outer_list[floor][1] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < Driver.N_FLOORS; floor++ {

		if Driver.Order_inner_list[floor] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	return Orders

}

//var lol bool = Driver.Bursdagskvinn()

func Network_client_main( /*New_order bool*/ ) {
	// connect to this socket
	conn, _ := net.Dial("tcp", "localhost:1201")

	for {
		// send to socket
		//fmt.Println(Orders_to_string_1())

		fmt.Fprintf(conn, Orders_to_string_1())
	}
}

/*func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}*/