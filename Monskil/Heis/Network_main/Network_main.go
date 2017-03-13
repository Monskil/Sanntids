package Network_main

import (
	"../bcast"
	//"../conn"
	"../Driver"
	"../localip"
	"../peers"
	"flag"
	"fmt"
	"os"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Message       string
	IP            string
	Current_floor int
	Direction     int
	Is_idle       bool
}

var received_msg = [4] /*N_FLOORS*/ [3] /*N_BUTTONS*/ int{
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
	{0, 0, 0},
}
var received_IP = "0"
var received_current_floor int = 0 //Driver.Direction //a.current_floor
var received_direction int = 0     //Driver.Current_floor //a.direction
var received_is_idle bool = true

///////////////////////////////////

var elev_1 = HelloMsg{Message: "0", IP: "0", Current_floor: 0, Direction: 0, Is_idle: true}
var elev_2 = HelloMsg{Message: "0", IP: "0", Current_floor: 0, Direction: 0, Is_idle: true}
var elev_3 = HelloMsg{Message: "0", IP: "0", Current_floor: 0, Direction: 0, Is_idle: true}
var num_elevs_online int = 1

////////////////////////////////

func Network_main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()
	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)

	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	LocalIP, _ := localip.LocalIP()
	// The example message. We just send one of these every second.
	go func() {

		for {
			current_floor1 := Driver.Current_floor
			Dir := Driver.IO_read_bit(Driver.MOTORDIR)
			idle := Driver.Elev_is_idle()
			Message := HelloMsg{Message: Orders_to_string(), IP: LocalIP, Current_floor: current_floor1, Direction: Dir, Is_idle: idle}
			/*Message.Message := Orders_to_string()
			Message.IP := 3 //my_ip //////////*/
			helloTx <- Message
			//fmt.Println("Current floor: ", current_floor)
			//fmt.Println("Direction: ", Dir)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		// case p := <-peerUpdateCh:
		/*fmt.Printf("Peer update:\n")
		fmt.Printf("  Peers:    %q\n", p.Peers)
		fmt.Printf("  New:      %q\n", p.New)
		fmt.Printf("  Lost:     %q\n", p.Lost)
		*/
		case a := <-helloRx:
			received_msg = String_to_orders(a.Message) ////////////
			received_IP = a.IP
			received_current_floor = a.Current_floor
			received_direction = a.Direction
			received_is_idle = a.Is_idle

			if received_IP == LocalIP {
				elev_1 = a
			} else if (received_IP != LocalIP) && (elev_2.IP == "0") {
				elev_2 = a
			} else if (received_IP != LocalIP) && (received_IP != elev_2.IP) {
				elev_3 = a
			} //HUSK Å SETTE ALLE MISTEDE HEISER TIL 0 SOM DE STÅR ØVERST I FILEN OG OPPDATERE NUM_ELEVS
			if elev_3.IP != "0" {
				num_elevs_online = 3
			} else if elev_2.IP != "0" {
				num_elevs_online = 2
			} else {
				num_elevs_online = 1
			}

		}
		//fmt.Println(received_msg)
		//fmt.Println(received_IP)
		//fmt.Println("Current floor: ", received_current_floor)
		//fmt.Println("Direction: ", received_direction)
		//fmt.Println("No orders: ", received_is_idle)*/
		fmt.Println("\n")
		fmt.Println(elev_1.IP)
		fmt.Println(elev_2.IP)
		fmt.Println(elev_3.IP)

	}
}

/*
func shall_me() {

	if  {

	}

}
*/
func Orders_to_string() string {

	var Orders string = ""
	for floor := 0; floor < Driver.N_FLOORS; floor++ {
		if Driver.Order_shared_outer_list[floor][0] == 1 {
			Orders = Orders + "1"
		} else {
			Orders = Orders + "0"
		}
	}
	for floor := 0; floor < Driver.N_FLOORS; floor++ {
		if Driver.Order_shared_outer_list[floor][1] == 1 {
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
	//fmt.Println(Orders)
	return Orders

}

func String_to_orders(Orders1 string) [4][3]int {
	//fmt.Println(Orders1)
	//var Orders int = [12] "000000000000" //UUUUDDDDCCCC (U = orders button_up | D = orders button_down | C = orders button_command)
	var Orders_list = [4] /*N_FLOORS*/ [3] /*N_BUTTONS*/ int{
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	//fmt.Println(Orders1)
	for i := 0; i < 4; i++ {
		if Orders1[i] == byte(49) {
			Orders_list[i][0] = 1
		} else if Orders1[i] == byte(48) {
			Orders_list[i][0] = 0
		} else {
			fmt.Println("Button_Up " + string(i) + " has an illegal value")
		}
	}
	for j := 4; j < 8; j++ {
		if Orders1[j] == byte(49) {
			Orders_list[j-4][1] = 1
		} else if Orders1[j] == byte(48) {
			Orders_list[j-4][1] = 0
		} else {
			fmt.Println("Button_Down " + string(j) + " has an illegal value")
		}
	}
	for k := 8; k < 12; k++ {
		if Orders1[k] == byte(49) {
			Orders_list[k-8][2] = 1
		} else if Orders1[k] == byte(48) {
			Orders_list[k-8][2] = 0
		} else {
			fmt.Println("Button_Command " + string(k) + "has an illegal value")
		}
	}
	//fmt.Println(Orders_list) //får over nettverk
	return Orders_list
}

func Order_compare_outer_list() {
	for {
		time.Sleep(1 * time.Second)
		counter := 0
		localIP, _ := localip.LocalIP()
		for floor := 0; floor < 4; floor++ {
			if (Driver.Order_shared_outer_list[floor][0] != received_msg[floor][0]) && (received_IP != localIP) /*&& Driver.Order_outer_list[floor][0] != 1*/ {
				Driver.Order_shared_outer_list[floor][0] = received_msg[floor][0]
				// Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_UP, floor, 0)
				counter++

			}
			if (Driver.Order_shared_outer_list[floor][1] != received_msg[floor][1]) && (received_IP != localIP) /*&& Driver.Order_outer_list[floor][1] != 1*/ {
				Driver.Order_shared_outer_list[floor][1] = received_msg[floor][1]
				//Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_DOWN, floor, 0)
				counter++
			}
		}
		if counter != 0 {
			Driver.Set_new_order_var()
		}
	}
}
