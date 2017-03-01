package main

import (
	"fmt"
	"sync"
	"time"

	. "./definitions"
	. "./driver"
)

func main() {
	var wg sync.WaitGroup
	button_chan := make(chan Button, 100)
	floor_sensor_chan := make(chan int, 100)

	local_order_chan := make(chan Order, 100)
	wg.Add(1)
	go Lift_control_init(button_chan, floor_sensor_chan, local_order_chan)
	wg.Wait()

}

func Lift_control_init(button_chan chan<- Button, floor_sensor_chan chan<- int,
	local_order_chan <-chan Order) {

	fmt.Println("In init")
	go isanythinghappening(button_chan, floor_sensor_chan) // Trenger nytt navn

	//Her har msg fra adm 4 deler: ordercat, order, floor, value (on/off). kanskje endre
	// ordercat siden ordertype ser ut som gjør noe merkelig.
	go checking_for_orders_from_admin(local_order_chan)
}

func checking_for_orders_from_admin(local_order_chan <-chan Order) { // Nytt navn
	for {
		fmt.Println("In for loop checking orders")
		select {
		case msg := <-local_order_chan:
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

func isanythinghappening(button_chan chan<- Button, floor_sensor_chan chan<- int) {
	for {
		fmt.Println("In isanythinghappening")
		inside_button_pressed(button_chan)
		outside_button_pressed(button_chan)
		floor_sensor_triggered(floor_sensor_chan)
		time.Sleep(1000 * time.Millisecond)
	}
}

func outside_button_pressed(button_chan chan<- Button) {
	fmt.Println("In outside")
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elev_get_button_signal(BUTTON_CALL_UP, floor) == 1 {
			button_chan <- Button{floor, BUTTON_CALL_UP}
		}
		if Elev_get_button_signal(BUTTON_CALL_DOWN, floor) == 1 {
			button_chan <- Button{floor, BUTTON_CALL_DOWN}
		}
	}
}

func inside_button_pressed(button_chan chan<- Button) {
	fmt.Println("In inside")
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elev_get_button_signal(BUTTON_COMMAND, floor) == 1 {
			button_chan <- Button{floor, BUTTON_COMMAND}
			// Endre "floor", se notat.txt for ekstra. (evt floor_number)
		}
	}
}

func floor_sensor_triggered(floor_sensor_chan chan<- int) {
	fmt.Println("In floor sensor")
	floor := Elev_get_floor_sensor_signal()
	if floor != -1 {
		floor_sensor_chan <- floor
	}
}

// NB: Wg.Wait() kan være mulighet/nødvendig en eller annen plass.
