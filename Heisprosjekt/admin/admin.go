package admin

import (
  "fmt"
  . "./../definitions"
  . "order_matrix"
  . "calculate_order"
  . "lift_properties"
)

// Mulig å slå sammen som interface{}?

func Adm_init(button_inside_chan <-chan int, button_outside_chan <-chan button, floor_sensor_chan <-chan int,
  local_order_chan chan<- order, adm_transmitt_chan chan<- udp, adm_receive_chan chan<- udp, peer_chan <-chan int,
  start_timer_chan chan<- string, time_out_chan <-chan string)  {

  go admin(button_inside_chan, button_outside_chan, floor_sensor_chan,
		local_order_chan, adm_transmitt_chan, adm_receive_chan, peer_chan,
		start_timer_chan, time_out_chan)
}

func admin(button_inside_chan <-chan int, button_outside_chan <-chan button, floor_sensor_chan <-chan int,
  local_order_chan chan<- order, adm_transmitt_chan chan<- udp, adm_receive_chan chan<- udp, peer_chan <-chan int,
  start_timer_chan chan<- string, time_out_chan <-chan string)  {

  orders = Create_order_matrix()
  properties = Create_lift_prop_list()
  ID = 0; //SETT ID (spør nett?)
  // Check om orders/prop list eksisterer noen andre plasser på nettet


  for {
    select {
    case bi := <- button_inside_chan:
      Add_order(orders, bi, ID, BUTTON_COMMAND)
      //Send melding ut til NW på adm_transmitt_chan
      if Get_state(properties, ID) == IDLE {
        find_new_order(orders, ID, properties, start_timer_chan, local_order_chan, adm_transmitt_chan)
      } // Problem med å sende melding om button pressed ut på nettet og deretter melding fra find_new_order?
      // evt legge ved hvilke ordre vi tar hver gang i find_new_order-melding => alle andre kan oppdatere.
      // Husk "problem" med at assign bare tar de som allerede finnes, så
      // må ha en måte å slå sammen her.

    case bo := <- button_outside_chan:
      Add_order(orders, bo.floor, ID, bo.button_dir)
      //Send melding ut til NW på adm_transmitt_chan
      if Get_state(properties, ID) == IDLE {
        find_new_order(orders, ID, properties, start_timer_chan, local_order_chan, adm_transmitt_chan)
      }

    case fs := <- floor_sensor_chan:
      Set_last_floor(properties, ID, fs)

      if Get_state(properties, ID) == MOVING {

        if Should_stop(orders, fs, ID) == true {
          //local_order_chan <-  Send "LIGHT", DIRN_STOP, -1, ON
          Set_state(properties, ID, DOOR_OPEN)
          start_timer_chan <- "DOOR_OPEN"
          Complete_order(orders, fs, ID)
          //Send melding ut til NW på adm_transmitt_chan
          // ID, "Stoppet", etasje (DOOR_OPEN)
        } else {
            //Send melding ut til NW på adm_transmitt_chan
            // ID, "kjørte forbi", etasje
        }
      }

    case time_out := <- time_out_chan: //Only for DOOR_OPEN timer atm.
      //local_order_chan <-  Send "LIGHT", DIRN_STOP, -1, OFF
      find_new_order(orders, ID, properties, start_timer_chan, local_order_chan, adm_transmitt_chan)

    case m := <- adm_receive_chan:

    }
  }
}


func find_new_order(orders [][]int, ID int, properties []int, start_timer_chan chan<- string,
  local_order_chan chan<- order, adm_transmitt_chan chan<- udp)  {

  new_dirn, dest := Calculate_order(orders, ID, properties)
  // Should change name on both module called from and function itself.
  // Default dest and new_dirn returned has to be undefined (-2,-1)
  if new_dirn == DIRN_STOP {
    Assign_orders(orders, dest, ID) //NB! Nå lagt til ALLE på den etasjen,
    // noe som er en forenkling som vi kunne gjøre. IKKE TESTET ENNÅ
    //local_order_chan <-  Send "LIGHT", DIRN_STOP, -1, ON
    Set_state(properties, ID, DOOR_OPEN)
    start_timer_chan <- "DOOR_OPEN"
    Complete_order(orders, dest, ID)
    //Send melding ut til NW på adm_transmitt_chan
    // ID, "Stoppet", etasje (DOOR_OPEN)
  } else if new_dirn == DIRN_DOWN || new_dirn == DIRN_UP {
    Assign_orders(orders, dest, ID)
    //local_order_chan <-  Send "DIRN", DIRN_UP/DOWN, -1, ON
    Set_state(properties, ID, MOVING)
    Set_dirn(properties, ID, new_dirn)
    //Send melding ut til NW på adm_transmitt_chan
    // ID, "Moving, desting (new order)", etasje
  } else { // new_dirn == -2
    Set_state(properties, ID, IDLE)
    //Send melding ut til NW på adm_transmitt_chan
    // ID, "IDLE", etasje
  }
}
