package admin

import (
  "fmt"
  "/order_matrix"
  "/calculate_order"
)

// Må import peers.go for peerUpdate med mindre flyttes til definitions

func Adm_init(button_inside_chan <-chan int, button_outside_chan <-chan button, floor_sensor_chan <-chan int,
  local_order_chan chan<- order, adm_transmitt_chan chan<- udp, adm_receive_chan chan<- udp, peer_chan <-chan int,
  start_timer_chan chan<- string, time_out_chan <-chan string, interrupt_timer_chan chan<- string)  {

  go admin(button_inside_chan, button_outside_chan, floor_sensor_chan,
		local_order_chan, adm_transmitt_chan, adm_receive_chan, peer_chan,
		start_timer_chan, time_out_chan, interrupt_timer_chan)
}

func admin(button_inside_chan <-chan int, button_outside_chan <-chan button, floor_sensor_chan <-chan int,
  local_order_chan chan<- order, adm_transmitt_chan chan<- udp, adm_receive_chan chan<- udp, peer_chan <-chan int,
  start_timer_chan chan<- string, time_out_chan <-chan string, interrupt_timer_chan chan<- string)  {

  orders = order_matrix.Create_order_matrix()
  properties = lift_properties.Create_lift_prop_list()
  ID = 0; //SETT ID


  for {
    select {
    case bi := <- button_inside_chan:
      //Send melding ut til NW på adm_transmitt_chan
      order_matrix.Add_order(orders, bi, ID, BUTTON_COMMAND)

    case bo := <- button_outside_chan:
      //Send melding ut til NW på adm_transmitt_chan
      order_matrix.Add_order(orders, bi.floor, ID, bi.button_dir)

    case fs := <- floor_sensor_chan:
      lift_properties.Set_last_floor(properties, ID, fs)

      if calculate_order.Should_stop(orders, fs, ID) == true {
        lift_properties.Set_state(properties, ID, DOOR_OPEN)
        //local_order_chan <-  Send "LIGHT", DIRN_STOP, -1, ON
        start_timer_chan <- "DOOR_OPEN"
        order_matrix.Complete_order(orders, fs, ID)
      }
      //Send melding ut til NW på adm_transmitt_chan. Eller etter if?

      case 
    }
  }
}
