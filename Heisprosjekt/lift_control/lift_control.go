package lift_control

import (
  "./driver"
)

func lift_control_init(button_inside_chan, button_outside_chan,
  floor_sensor_chan, local_order_chan chan int)  {

  go lift_control(button_inside_chan, button_outside_chan,
    floor_sensor_chan, local_order_chan)
}


func lift_control(button_inside_chan, button_outside_chan,
  floor_sensor_chan, local_order_chan chan int)  {


  //Trenger sikkert select/å pakke inn i funksjoner som er go-rutiner
  inside_button := inside_button_pressed()
  if inside_button {
    button_inside_chan <- inside_button
  }

  outside_button := outside_button_pressed()
  if outside_button {
    button_outside_chan <- outside_button
  }

  floor_sensor := floor_sensor_triggered()
  if floor_sensor {
    floor_sensor_chan <- floor_sensor
  }

  //Her har msg fra adm tre deler: info, type, value
  go func ()  {
    for {
      case msg := <- local_order_chan:
        if msg.info == light {
          if msg.type == door {
            driver.Elev_set_door_open_lamp(msg.value)
          } else {
            driver.Elev_set_button_lamp(msg.value)
          }
        } else if msg.info == dir {
          if msg.type == stop {
            driver.Elev_set_motor_direction(DIRN_STOP)
            driver.Elev_set_door_open_lamp(msg.value)
          } else {
            driver.Elev_set_motor_direction(msg.value)
          }
        }
      }
    }
  }
}
