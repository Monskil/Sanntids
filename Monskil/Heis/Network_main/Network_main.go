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
	"strings"
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
var received_current_floor int = 0
var received_direction int = 0
var received_is_idle bool = true

var elev_1 = HelloMsg{Message: "0", IP: "000", Current_floor: 0, Direction: 0, Is_idle: false} //THIS ELEV
var elev_2 = HelloMsg{Message: "0", IP: "000", Current_floor: 0, Direction: 0, Is_idle: false}
var elev_3 = HelloMsg{Message: "0", IP: "000", Current_floor: 0, Direction: 0, Is_idle: false}
var num_elevs_online int = 1
var elev_1_ID int = 0
var elev_2_ID int = 0
var elev_3_ID int = 0
var elev_new_ID int = 0
var num_unique_IPs int = 1

//var IP_list = [20]string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"}
var elev_lost string = ""

func Network_main(Order_chan chan bool, full_array_chan chan bool) {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(1201 /*15647*/, id, peerTxEnable)
	go peers.Receiver(1201 /*15647*/, peerUpdateCh)
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	go bcast.Transmitter(8081 /*16569*/, helloTx)
	go bcast.Receiver(8081 /*16569*/, helloRx)

	LocalIP, _ := localip.LocalIP()
	go func() {
		for {
			current_floor1 := Driver.Current_floor
			Dir := Driver.IO_read_bit(Driver.MOTORDIR)
			idle := Driver.Elev_is_idle(Order_chan)
			Message := HelloMsg{Message: Orders_to_string(), IP: LocalIP, Current_floor: current_floor1, Direction: Dir, Is_idle: idle}
			helloTx <- Message
			time.Sleep(5 * time.Millisecond)
		}
	}()

	for {
		select {

		case a := <-helloRx:
			received_msg = String_to_orders(a.Message)
			received_IP = a.IP
			received_current_floor = a.Current_floor
			received_direction = a.Direction
			received_is_idle = a.Is_idle
			Set_ID_from_IP()

			if received_IP == LocalIP {
				elev_1 = a
			} else if (received_IP != LocalIP) && (elev_2.IP == "000") {
				elev_2.IP = received_IP
			} else if received_IP == elev_2.IP {
				elev_2 = a
			} else if (received_IP != LocalIP) && (received_IP != elev_2.IP) {
				elev_3 = a
			}
			if elev_3.IP != "000" {
				num_elevs_online = 3
			} else if elev_2.IP != "000" {
				num_elevs_online = 2
			} else {
				num_elevs_online = 1
			}
			/*for i := 0; i < 20; i++ {
				if IP_list[i] != "0" {
					IP_list[i] = received_IP
					break
				} else {
					full_array_chan <- true
				}
			}*/
		///////////////////////////////////
		case p := <-peerUpdateCh:
			//fmt.Printf("Peer update:\n")
			//fmt.Printf("  Peers:    %q\n", p.Peers)
			//fmt.Printf("  New:      %q\n", p.New)
			//fmt.Printf("  Lost:     %q\n", p.Lost)
			/*if p.New != "" {
				//num_elevs_online = num_elevs_online + 1
			}*/

			if strings.Join(p.Lost, "") != "" {
				num_elevs_online = num_elevs_online - 1
				elev_lost = (strings.Join(p.Lost, ""))
				if _, err := fmt.Sscanf(elev_lost, "peer-129.241.187.%3d", &elev_new_ID); err == nil {
					//fmt.Println("new_id: ", elev_new_ID)
					//fmt.Println(elev_2_ID)
					//fmt.Println(elev_3_ID)
					if elev_new_ID == elev_2_ID {
						elev_2_ID = 0
						elev_new_ID = 0
						elev_2 = HelloMsg{Message: "0", IP: "000", Current_floor: 0, Direction: 0, Is_idle: false}
					} else if elev_new_ID == elev_3_ID {
						elev_3_ID = 0
						elev_new_ID = 0
						elev_3 = HelloMsg{Message: "0", IP: "000", Current_floor: 0, Direction: 0, Is_idle: false}
					}
				}

			}
			////////////////////////////////////////////////
		}

		//fmt.Println(num_elevs_online)
		//fmt.Println(received_msg)
		//fmt.Println(received_IP)
		//fmt.Println("Current floor: ", received_current_floor)
		//fmt.Println("Direction: ", received_direction)
		//fmt.Println("No orders: ", received_is_idle)*/
		//fmt.Println("\n")
		//fmt.Println("Number of elevators online: ", num_elevs_online)
		//fmt.Println(elev_1_ID)
		//fmt.Println(elev_2_ID)
		//fmt.Println(elev_3_ID)

	}

}

func Set_ID_from_IP() {
	if _, err := fmt.Sscanf(elev_1.IP, "129.241.187.%3d", &elev_1_ID); err == nil {
		//fmt.Println(elev_1_ID)
	}
	if _, err := fmt.Sscanf(elev_2.IP, "129.241.187.%3d", &elev_2_ID); err == nil {
		//fmt.Println(elev_2_ID)
	}
	if _, err := fmt.Sscanf(elev_3.IP, "129.241.187.%3d", &elev_3_ID); err == nil {
		//fmt.Println(elev_2_ID)
	}
}

func Cost_function() {
	for {
		time.Sleep(10 * time.Millisecond)
		var elev_sufficient bool = false
		var elev_1_difference int = 0
		var elev_2_difference int = 0
		var elev_3_difference int = 0

		for floor := 0; floor < Driver.N_FLOORS; floor++ {

			if Driver.Order_shared_outer_list[floor][0] == 1 {

				if num_elevs_online == 1 {
					Driver.Order_outer_list[floor][0] = 1
				} else if num_elevs_online == 2 { //////////////////////////////////////////////// OPP 2 Heiser

					elev_1_difference = floor - elev_1.Current_floor
					elev_2_difference = floor - elev_2.Current_floor
					elev_3_difference = 5

					if (elev_1.Direction == 0 || elev_1.Is_idle == true) && elev_1.Current_floor <= floor {

						elev_sufficient = true
					}
					if (elev_2.Direction == 0 || elev_2.Is_idle == true) && elev_2.Current_floor <= floor {
						if elev_1_difference > elev_2_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_2_difference {
							if elev_1_ID < elev_2_ID {
								elev_sufficient = false
							}
						}
					}
				}
			}
			if elev_sufficient == true {
				Driver.Order_outer_list[floor][0] = 1
				//fmt.Println(elev_sufficient)
				elev_sufficient = false
			} /*else if elev_sufficient == false { ////////////////////////////////////////////////DENNE ELSEN
				Driver.Order_outer_list[floor][0] = 0
			}*/

			if Driver.Order_shared_outer_list[floor][1] == 1 {

				if num_elevs_online == 1 {
					Driver.Order_outer_list[floor][1] = 1
				} else if num_elevs_online == 2 {

					elev_1_difference = floor - elev_1.Current_floor
					elev_2_difference = floor - elev_2.Current_floor
					elev_3_difference = 5

					if (elev_1.Direction == 0 || elev_1.Is_idle == true) && elev_1.Current_floor <= floor {

						elev_sufficient = true
					}
					if (elev_2.Direction == 0 || elev_2.Is_idle == true) && elev_2.Current_floor <= floor {
						if elev_1_difference > elev_2_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_2_difference {
							if elev_1_ID < elev_2_ID {
								elev_sufficient = false
							}
						}
					}
				}
			}
			if elev_sufficient == true {
				Driver.Order_outer_list[floor][1] = 1
				elev_sufficient = false
			} else if num_elevs_online == 3 { /////////////////////////////////////////////////////////////////// OPP 3 Heiser

				elev_1_difference = floor - elev_1.Current_floor
				elev_2_difference = floor - elev_2.Current_floor
				elev_3_difference = floor - elev_3.Current_floor

				if (elev_1.Direction == 0 || elev_1.Is_idle == true) && elev_1.Current_floor <= floor {

					elev_sufficient = true
				}
				if (elev_2.Direction == 0 || elev_2.Is_idle == true) && elev_2.Current_floor <= floor {
					if elev_1_difference > elev_2_difference {
						elev_sufficient = false
					} else if elev_1_difference == elev_2_difference {
						if elev_1_ID < elev_2_ID {
							elev_sufficient = false
						}
					}
				}
				if (elev_3.Direction == 0 || elev_3.Is_idle == true) && elev_3.Current_floor <= floor {
					if elev_1_difference > elev_3_difference {
						elev_sufficient = false
					} else if elev_1_difference == elev_3_difference {
						if elev_1_ID < elev_3_ID {
							elev_sufficient = false
						}
					}
				}
			}

			if elev_sufficient == true {
				Driver.Order_outer_list[floor][0] = 1
				elev_sufficient = false
			} /* else if elev_sufficient == false { ////////////////////////////////////////////////DENNE ELSEN
				Driver.Order_outer_list[floor][0] = 0
			}*/

			if Driver.Order_shared_outer_list[floor][1] == 1 {

				if num_elevs_online == 3 {

					elev_1_difference = floor - elev_1.Current_floor
					elev_2_difference = floor - elev_2.Current_floor
					elev_3_difference = floor - elev_3.Current_floor

					if (elev_1.Direction == 0 || elev_1.Is_idle == true) && elev_1.Current_floor <= floor {

						elev_sufficient = true
					}
					if (elev_2.Direction == 0 || elev_2.Is_idle == true) && elev_2.Current_floor <= floor {
						if elev_1_difference > elev_2_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_2_difference {
							if elev_1_ID < elev_2_ID {
								elev_sufficient = false
							}
						}
					}
					if (elev_3.Direction == 0 || elev_3.Is_idle == true) && elev_3.Current_floor <= floor {
						if elev_1_difference > elev_3_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_3_difference {
							if elev_1_ID < elev_3_ID {
								elev_sufficient = false
							}
						}
					}
				}
			}
			if elev_sufficient == true {
				Driver.Order_outer_list[floor][1] = 1
				elev_sufficient = false
			} /* else if elev_sufficient == false { ////////////////////////////////////////////////DENNE ELSEN
				Driver.Order_outer_list[floor][1] = 0
			}*/
			/////////////////////////////////////////////////////////////////////////////////////////////////SLUTT OPP 3 HEISER
		} /////////////////////////////////////////////////////////////////////////////////////////////////////START NED 2 HEISER
		for floor := Driver.N_FLOORS - 1; floor >= 0; floor-- {

			if Driver.Order_shared_outer_list[floor][0] == 1 {
				if num_elevs_online == 2 {

					elev_1_difference = elev_1.Current_floor - floor
					elev_2_difference = elev_2.Current_floor - floor
					elev_3_difference = 5

					if (elev_1.Direction == 1 || elev_1.Is_idle == true) && elev_1.Current_floor >= floor {

						elev_sufficient = true
					}
					if (elev_2.Direction == 1 || elev_2.Is_idle == true) && elev_2.Current_floor >= floor {
						if elev_1_difference > elev_2_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_2_difference {
							if elev_1_ID < elev_2_ID {
								elev_sufficient = false
							}
						}
					}
				}
			}
			if elev_sufficient == true {
				Driver.Order_outer_list[floor][0] = 1
				elev_sufficient = false
			} /* else if elev_sufficient == false { ////////////////////////////////////////////////DENNE ELSEN
				Driver.Order_outer_list[floor][0] = 0
			}*/

			if Driver.Order_shared_outer_list[floor][1] == 1 {

				if num_elevs_online == 2 {

					elev_1_difference = elev_1.Current_floor - floor
					elev_2_difference = elev_2.Current_floor - floor
					elev_3_difference = 5

					if (elev_1.Direction == 1 || elev_1.Is_idle == true) && elev_1.Current_floor >= floor {
						elev_sufficient = true
					}
					if (elev_2.Direction == 1 || elev_2.Is_idle == true) && elev_2.Current_floor >= floor {
						if elev_1_difference > elev_2_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_2_difference {
							if elev_1_ID < elev_2_ID {
								elev_sufficient = false
							}
						}
					}
				}
			}
			if elev_sufficient == true {
				Driver.Order_outer_list[floor][1] = 1
				elev_sufficient = false
			} /*else if elev_sufficient == false { ////////////////////////////////////////////////DENNE ELSEN
				Driver.Order_outer_list[floor][1] = 0
			}*/
			/////////////////////////////////////////////////////////////////////////////////////////////////////SLUTT NED 2 HEISER
			/////////////////////////////////////////////////////////////////////////////////////////////////////START NED 3 HEISER
			if Driver.Order_shared_outer_list[floor][0] == 1 {

				if num_elevs_online == 3 {

					elev_1_difference = elev_1.Current_floor - floor
					elev_2_difference = elev_2.Current_floor - floor
					elev_3_difference = elev_3.Current_floor - floor

					if (elev_1.Direction == 1 || elev_1.Is_idle == true) && elev_1.Current_floor >= floor {

						elev_sufficient = true
					}
					if (elev_2.Direction == 1 || elev_2.Is_idle == true) && elev_2.Current_floor >= floor {
						if elev_1_difference > elev_2_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_2_difference {
							if elev_1_ID < elev_2_ID {
								elev_sufficient = false
							}
						}
					}
					if (elev_3.Direction == 1 || elev_3.Is_idle == true) && elev_3.Current_floor >= floor {
						if elev_1_difference > elev_3_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_3_difference {
							if elev_1_ID < elev_3_ID {
								elev_sufficient = false
							}
						}
					}
				}
			}
			if elev_sufficient == true {
				Driver.Order_outer_list[floor][0] = 1
				elev_sufficient = false
			} /*else if elev_sufficient == false { ////////////////////////////////////////////////DENNE ELSEN
				Driver.Order_outer_list[floor][0] = 0
			}*/

			if Driver.Order_shared_outer_list[floor][1] == 1 {

				if num_elevs_online == 3 {

					elev_1_difference = elev_1.Current_floor - floor
					elev_2_difference = elev_2.Current_floor - floor
					elev_3_difference = elev_3.Current_floor - floor

					if (elev_1.Direction == 1 || elev_1.Is_idle == true) && elev_1.Current_floor >= floor {
						elev_sufficient = true
					}
					if (elev_2.Direction == 1 || elev_2.Is_idle == true) && elev_2.Current_floor >= floor {
						if elev_1_difference > elev_2_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_2_difference {
							if elev_1_ID < elev_2_ID {
								elev_sufficient = false
							}
						}
					}
					if (elev_3.Direction == 1 || elev_3.Is_idle == true) && elev_3.Current_floor >= floor {
						if elev_1_difference > elev_3_difference {
							elev_sufficient = false
						} else if elev_1_difference == elev_3_difference {
							if elev_1_ID < elev_3_ID {
								elev_sufficient = false
							}
						}
					}
				}
			}
			if elev_sufficient == true {
				Driver.Order_outer_list[floor][1] = 1
				elev_sufficient = false
			} /*else if elev_sufficient == false { ////////////////////////////////////////////////DENNE ELSEN
				Driver.Order_outer_list[floor][1] = 0
			}*/
			/////////////////////////////////////////////////////////////////////////////////////////////////////SLUTT NED 3 HEISER
		}
	}
}

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
		counter := 0
		localIP, _ := localip.LocalIP()
		for floor := 0; floor < 4; floor++ {
			if Driver.Order_shared_outer_list[floor][0] != received_msg[floor][0] && (received_IP != localIP) /*&& Driver.Order_outer_list[floor][0] != 1*/ {
				Driver.Order_shared_outer_list[floor][0] = received_msg[floor][0]
				// Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_UP, floor, 0)
				counter++

			}
			if Driver.Order_shared_outer_list[floor][1] != received_msg[floor][1] && (received_IP != localIP) /*&& Driver.Order_outer_list[floor][1] != 1*/ {
				Driver.Order_shared_outer_list[floor][1] = received_msg[floor][1]
				//Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_DOWN, floor, 0)
				counter++
			}
		}
		if counter != 0 {
			Driver.Set_new_order_var()
		}
		time.Sleep(5 * time.Millisecond)
	}
}
