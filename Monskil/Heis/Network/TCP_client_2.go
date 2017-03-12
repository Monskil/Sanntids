package Network

import (
	"../Driver"
	"fmt"
	"net"
	//"os"
	"time"
)

func Orders_to_string_2() string {

	test_inner := [4]int{0, 0, 0, 0}
	test_outer := [4][2]int{
		{0, 0},
		{0, 0},
		{0, 0},
		{0, 1},
	}

	var Orders string = "" //UUUUDDDDCCCC (U = orders button_up | D = orders button_down | C = orders button_command)
	for floor := 0; floor < Driver.N_FLOORS; floor++ {
		if /*Driver.Order_outer_list[floor][0] */ test_outer[floor][0] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < Driver.N_FLOORS; floor++ {
		if /*Driver.Order_outer_list[floor][1] */ test_outer[floor][1] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < Driver.N_FLOORS; floor++ {

		if /*Driver.Order_inner_list[floor] */ test_inner[floor] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	//fmt.Println(Orders)
	return Orders

}

//var lol bool = Driver.Bursdagskvinn()
//LABPLASS 01 = 140
//LABPLASS 06 = 146
//LABPLASS 14 = 142
//LABPLASS 16 = 147
func Network_client_main_2( /*New_order bool*/ ) {
	// connect to this socket
	conn, err := net.Dial("tcp" /*,"129.241.187.142:1201" */, "localhost:8081")
	//fmt.Println(conn)

	/*var monvar net.Conn = (&{{0xc82005a150}})

	if conn != monvar {
		fmt.Println("ikke konnekta")
	}*/

	for {
		if err != nil {
			//fmt.Println("error")
		} else {
			// send to socket
			//fmt.Println(Orders_to_string_1())
			//time.Sleep(500 * time.Millisecond)
			//fmt.Println(Orders_to_string_1())
			//fmt.Println("Client")
			time.Sleep(300 * time.Millisecond)
			fmt.Fprintf(conn, Orders_to_string_2())
		}
	}
}

/*func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}*/
