package lift_control

import (
	"fmt"
	"sync"
	"time"

	. "../definitions"
	. "../driver"
)

func Lift_control(buttonChan chan<- Button, floorSensorChan chan<- int,
	localOrderChan <-chan Order) {
	var wgL sync.WaitGroup
	wgL.Add(1)

	fmt.Println("LC_init")
	go isanythinghappening(buttonChan, floorSensorChan) // Trenger nytt navn

	go checking_for_orders_from_admin(localOrderChan)
	wgL.Wait()
}

func checking_for_orders_from_admin(localOrderChan <-chan Order) { // Nytt navn
	for {
		select {
		case msg := <-localOrderChan:
			switch msg.Category {
			case "LIGHT":
				Elev_set_button_lamp(msg.Order, msg.Floor, msg.Value)
			case "DOOR":
				Elev_set_door_open_lamp(msg.Value)
			case "DIRN":
				if msg.Order == DIRN_STOP {
					Elev_set_motor_direction(DIRN_STOP)
					Elev_set_door_open_lamp(ON)
				} else {
					Elev_set_motor_direction(msg.Order)
				}
			case "FLOOR_LIGHT":
				Elev_set_floor_indicator(msg.Floor)
				fmt.Println("Lift: Floor light set", msg.Floor)
			}
		}
	}
}

func isanythinghappening(buttonChan chan<- Button, floorSensorChan chan<- int) {
	for {
		//fmt.Println("isanythinghappening")
		buttonPressed(buttonChan)
		floor_sensor_triggered(floorSensorChan)
		time.Sleep(70 * time.Millisecond)
	}
}

func buttonPressed(buttonChan chan<- Button) {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elev_get_button_signal(BUTTON_COMMAND, floor) == 1 {
			buttonChan <- Button{floor, BUTTON_COMMAND}
			// Endre "floor", se notat.txt for ekstra. (evt floor_number)
		}
		if Elev_get_button_signal(BUTTON_CALL_UP, floor) == 1 {
			buttonChan <- Button{floor, BUTTON_CALL_UP}
		}
		if Elev_get_button_signal(BUTTON_CALL_DOWN, floor) == 1 {
			buttonChan <- Button{floor, BUTTON_CALL_DOWN}
		}
	}
}

func floor_sensor_triggered(floorSensorChan chan<- int) {
	floor := Elev_get_floor_sensor_signal()
	if floor != -1 {
		floorSensorChan <- floor
	}
}
