package Network

import "net"
import "fmt"
import "bufio"
import "../Driver"
import "time"

//import "strings"

func Client_main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "localhost:8081")
  for {
    // read in input from stdin
    /*reader := bufio.NewReader(Orders_to_string_client())
      fmt.Print("Klient: ")
      text, _ := reader.ReadString('\n')
    */
    // send to socket
    fmt.Fprintf(conn, Orders_to_string_client())
    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('2')
    Orders_list = String_to_orders_client(message)
    //fmt.Print("Message from server: " + message)
    //time.Sleep(2 * time.Second)
  }
}

func Orders_to_string_client() string {
  /*
     test_inner := [4]int{0, 0, 0, 0}
     test_outer := [4][2]int{
       {0, 0},
       {1, 0},
       {0, 0},
       {0, 0},
     }
  */
  var Orders string = "" //UUUUDDDDCCCC (U = orders button_up | D = orders button_down | C = orders button_command)
  for floor := 0; floor < Driver.N_FLOORS; floor++ {
    if Driver.Order_outer_list[floor][0] /* test_outer[floor][0] */ == 1 {
      Orders = Orders + "1"
    } else {
      Orders = Orders + "0"
    }
  }
  for floor := 0; floor < Driver.N_FLOORS; floor++ {
    if Driver.Order_outer_list[floor][1] /* test_outer[floor][1] */ == 1 {
      Orders = Orders + "1"
    } else {
      Orders = Orders + "0"
    }
  }
  for floor := 0; floor < Driver.N_FLOORS; floor++ {

    if Driver.Order_inner_list[floor] /* test_inner[floor]*/ == 1 {
      Orders = Orders + "1"
    } else {
      Orders = Orders + "0"
    }
  }
  //fmt.Println(Orders)
  return Orders + "2"

}

var Orders_list = [4] /*N_FLOORS*/ [3] /*N_BUTTONS*/ int{
  {0, 0, 0},
  {0, 0, 0},
  {0, 0, 0},
  {0, 0, 0},
}

func String_to_orders_client(Orders1 string) [4][3]int {
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
  //fmt.Println(Orders_list)
  return Orders_list
}

func Order_compare_outer_lists_client() {
  for {
    time.Sleep(1 * time.Second)
    counter := 0
    for floor := 0; floor < 4; floor++ {
      if Driver.Order_outer_list[floor][0] != Orders_list[floor][0] && Driver.Order_outer_list[floor][0] != 1 {
        Driver.Order_outer_list[floor][0] = Orders_list[floor][0]
        // Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_UP, floor, 0)
        counter++

      }
      if Driver.Order_outer_list[floor][1] != Orders_list[floor][1] && Driver.Order_outer_list[floor][1] != 1 {
        Driver.Order_outer_list[floor][1] = Orders_list[floor][1]
        //Driver.Elev_set_button_lamp(Driver.BUTTON_CALL_DOWN, floor, 0)
        counter++
      }
    }
    if counter != 0 {
      Driver.Set_new_order_var()
    }
  }

}
