package lift_control

import (
  "./driver"
  "./definitions"
)

func lift_control_init(button_inside_chan, button_outside_chan,
  floor_sensor_chan, local_order_chan chan int)  {

  go lift_control(button_inside_chan, button_outside_chan,
    floor_sensor_chan, local_order_chan)
}


func lift_control(button_inside_chan, floor_sensor_chan, local_order_chan chan int,
  button_outside_chan chan button)  {


  //Trenger sikkert select/å pakke inn i funksjoner som er go-rutiner
  inside_button := inside_button_pressed()
  if inside_button != -1 {
    button_inside_chan <- inside_button
  }

  outside_button := outside_button_pressed()
  if outside_button != -1 {
    button_outside_chan <- outside_button
  }

  floor_sensor := driver.Elev_get_floor_sensor_signal()
  if floor_sensor != -1 {
    floor_sensor_chan <- floor_sensor
  }

  //Her har msg fra adm 4 deler: ordercat, order, floor, value (on/off). kanskje endre
  // ordercat siden ordertype ser ut som gjør noe merkelig.
  go func ()  {
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
}

func outside_button_pressed(button_outside_chan chan button) int {
  for floor := 0; floor < N_FLOORS; floor++ {
    if driver.Elev_get_button_signal(BUTTON_CALL_UP, floor) {
      button_outside_chan <- button{"U",floor}
      // Trenger sikkert noe sånn at det ikke blir spammet på hvis
      // noen holder inn knappen, og noe sånn at knapp høyere opp
      // blir tatt hvis lengre inn blir holdt inn.
      // (Antar samme løsning, ble brukt forrige heisprosjekt)
    }
    if driver.Elev_get_button_signal(BUTTON_CALL_DOWN, floor) {
      button_outside_chan <- button{"D",floor}
    }
  }
}

func inside_button_pressed(button_inside_chan chan int) int {
  for floor := 0; floor < N_FLOORS; floor++ {
    if driver.Elev_get_button_signal(BUTTON_COMMAND, floor) {
      button_inside_chan <- floor
      // Trenger sikkert noe sånn at det ikke blir spammet på hvis
      // noen holder inn knappen, og noe sånn at knapp høyere opp
      // blir tatt hvis lengre inn blir holdt inn.
      // (Antar samme løsning, ble brukt forrige heisprosjekt)
      // Endre "floor", se notat.txt for ekstra.
    }
  }
}
