package lift_control

import (
  "./driver"
  "github.com/Arktoz/Sanntid/Heisprosjekt/definitions.go"
)

func lift_control_init(button_inside_chan chan<- int, floor_sensor_chan chan<- int,
  button_outside_chan chan<- button, local_order_chan <-chan order)  {

  go lift_control(button_inside_chan, button_outside_chan,
    floor_sensor_chan, local_order_chan)
}


func lift_control(button_inside_chan chan<- int, floor_sensor_chan chan<- int,
  button_outside_chan chan<- button, local_order_chan <-chan order)  {

  go isanythinghappening(button_inside_chan, floor_sensor_chan,
    button_outside_chan) // Trenger nytt navn

  //Her har msg fra adm 4 deler: ordercat, order, floor, value (on/off). kanskje endre
  // ordercat siden ordertype ser ut som gjør noe merkelig.
  go checking_for_orders_from_admin(local_order_chan)
}

func checking_for_orders_from_admin(local_order_chan <-chan order)  { // Nytt navn
  for {
    select{
      case msg := <- local_order_chan:
        if msg.cat == "LIGHT" {
          if msg.order == "DOOR" {
            driver.Elev_set_door_open_lamp(msg.value)
          } else {
            driver.Elev_set_button_lamp(msg.order, msg.floor, msg.value)
          }
        } else if msg.cat == "DIR" {
          if msg.order == "DIRN_STOP" {
            driver.Elev_set_motor_direction(DIRN_STOP)
            driver.Elev_set_door_open_lamp(ON)
          } else {
            driver.Elev_set_motor_direction(msg.order)
          }
        }
    }
  }
}

func isanythinghappening(button_inside_chan chan<- int, floor_sensor_chan chan<- int,
  button_outside_chan chan<- button)  {
  for {
    outside_button_pressed(button_outside_chan)
    inside_button_pressed(button_inside_chan)
    floor_sensor_triggered(floor_sensor_chan)
    time.Sleep(150*time.Millisecond)
  }
}

func outside_button_pressed(button_outside_chan chan<- button) int {
  for floor := 0; floor < N_FLOORS; floor++ {
    if driver.Elev_get_button_signal(BUTTON_CALL_UP, floor) {
      button_outside_chan <- button{"U",floor}
    }
    if driver.Elev_get_button_signal(BUTTON_CALL_DOWN, floor) {
      button_outside_chan <- button{"D",floor}
    }
  }
}

func inside_button_pressed(button_inside_chan chan<- int) int {
  for floor := 0; floor < N_FLOORS; floor++ {
    if driver.Elev_get_button_signal(BUTTON_COMMAND, floor) {
      button_inside_chan <- floor
      // Endre "floor", se notat.txt for ekstra. (evt floor_number)
    }
  }
}

func floor_sensor_triggered(floor_sensor_chan chan<- int)  {
  floor := driver.Elev_get_floor_sensor_signal()
  if floor != -1 {
    floor_sensor_chan <- floor
  }
}
