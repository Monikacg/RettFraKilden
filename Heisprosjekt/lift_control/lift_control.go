package lift_control

import (
	"time"

	. "./../definitions"
	. "./../driver"
)

func Lift_control_init(buttonChan chan<- Button, floorSensorChan chan<- int,
	localOrderChan <-chan Order) {

	go isanythinghappening(buttonChan, floorSensorChan) // Trenger nytt navn

	//Her har msg fra adm 4 deler: ordercat, order, floor, value (on/off). kanskje endre
	// ordercat siden ordertype ser ut som gjør noe merkelig.
	go checking_for_orders_from_admin(localOrderChan)
}

func checking_for_orders_from_admin(localOrderChan <-chan Order) { // Nytt navn
	for {
		select {
		case msg := <-localOrderChan:
			/*
				Tidligere brukt. Ta bort om det under fungerer
					if msg.Category == "LIGHT" {
						if msg.Order == "DOOR" {
							Elev_set_door_open_lamp(msg.Value)
						} else {
							Elev_set_button_lamp(msg.Order, msg.Floor, msg.Value)
						}
					} else if msg.Category == "DIRN" {
						if msg.Order == DIRN_STOP {
							Elev_set_motor_direction(DIRN_STOP)
							Elev_set_door_open_lamp(ON)
						} else {
							Elev_set_motor_direction(msg.Order)
						}
					}
			*/

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
			}
		}
	}
}

func isanythinghappening(buttonChan chan<- Button, floorSensorChan chan<- int) {
	for {
		inside_button_pressed(buttonChan)
		outside_button_pressed(buttonChan)
		floor_sensor_triggered(floorSensorChan)
		time.Sleep(150 * time.Millisecond)
	}
}

func outside_button_pressed(buttonChan chan<- Button) int {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elev_get_button_signal(BUTTON_CALL_UP, floor) {
			buttonChan <- Button{floor, BUTTON_CALL_UP}
		}
		if Elev_get_button_signal(BUTTON_CALL_DOWN, floor) {
			buttonChan <- Button{floor, BUTTON_CALL_DOWN}
		}
	}
}

func inside_button_pressed(buttonChan chan<- Button) int {
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elev_get_button_signal(BUTTON_COMMAND, floor) {
			buttonChan <- Button{floor, BUTTON_COMMAND}
			// Endre "floor", se notat.txt for ekstra. (evt floor_number)
		}
	}
}

func floor_sensor_triggered(floorSensorChan chan<- int) {
	floor := Elev_get_floor_sensor_signal()
	if floor != -1 {
		floorSensorChan <- floor
	}
}

// NB: Wg.Wait() kan være mulighet/nødvendig en eller annen plass.
