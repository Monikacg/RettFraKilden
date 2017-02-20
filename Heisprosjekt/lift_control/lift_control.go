package lift_control

import (
  . "./../driver"
  . "./../definitions"
)

func Lift_control_init(button_chan chan<- Button, floor_sensor_chan chan<- int,
   local_order_chan <-chan Order)  {

   go isanythinghappening(button_chan, floor_sensor_chan) // Trenger nytt navn

   //Her har msg fra adm 4 deler: ordercat, order, floor, value (on/off). kanskje endre
   // ordercat siden ordertype ser ut som gjør noe merkelig.
   go checking_for_orders_from_admin(local_order_chan)
}

func checking_for_orders_from_admin(local_order_chan <-chan Order)  { // Nytt navn
  for {
    select{
      case msg := <- local_order_chan:
        if msg.Cat == "LIGHT" {
          if msg.Order == "DOOR" {
            Elev_set_door_open_lamp(msg.Value)
          } else {
            Elev_set_button_lamp(msg.Order, msg.Floor, msg.Value)
          }
        } else if msg.Cat == "DIRN" {
          if msg.Order == DIRN_STOP {
            Elev_set_motor_direction(DIRN_STOP)
            Elev_set_door_open_lamp(ON)
          } else {
            Elev_set_motor_direction(msg.Order)
          }
        }
    }
  }
}

func isanythinghappening(button_chan chan<- Button, floor_sensor_chan chan<- int)  {
  for {
    inside_button_pressed(button_chan)
    outside_button_pressed(button_chan)
    floor_sensor_triggered(floor_sensor_chan)
    time.Sleep(150*time.Millisecond)
  }
}

func outside_button_pressed(button_chan chan<- Button) int {
  for floor := 0; floor < N_FLOORS; floor++ {
    if Elev_get_button_signal(BUTTON_CALL_UP, floor) {
      button_chan <- Button{floor, BUTTON_CALL_UP}
    }
    if Elev_get_button_signal(BUTTON_CALL_DOWN, floor) {
      button_chan <- Button{floor, BUTTON_CALL_DOWN}
    }
  }
}

func inside_button_pressed(button_chan chan<- Button) int {
  for floor := 0; floor < N_FLOORS; floor++ {
    if Elev_get_button_signal(BUTTON_COMMAND, floor) {
      button_chan <- Button{floor, BUTTON_COMMAND}
      // Endre "floor", se notat.txt for ekstra. (evt floor_number)
    }
  }
}

func floor_sensor_triggered(floor_sensor_chan chan<- int)  {
  floor := Elev_get_floor_sensor_signal()
  if floor != -1 {
    floor_sensor_chan <- floor
  }
}

// NB: Wg.Wait() kan være mulighet/nødvendig en eller annen plass.
