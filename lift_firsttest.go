package main

import (
	"fmt"
	"sync"
	"time"

	. "./definitions"
	. "./driver"
)

func main() {
	Elev_init()
	var wg sync.WaitGroup
	button_chan := make(chan Button, 100)
	floor_sensor_chan := make(chan int, 100)

	local_order_chan := make(chan Order, 100)
	wg.Add(1)
	go Lift_control_init(button_chan, floor_sensor_chan, local_order_chan)
	go sending_stuff(local_order_chan)
	//go read_floor_sensor(floor_sensor_chan, local_order_chan, button_chan)
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

/*func read_floor_sensor(floor_sensor_chan <-chan int, local_order_chan chan<- Order, button_chan <-chan Button) {
	var saved_floor int
	for {
		select {
		case f := <-floor_sensor_chan:
			fmt.Println("Floor sensor registrered at floor: ", f)
			local_order_chan <- Order{"FLOOR_LIGHT", NOT_VALID, f, NOT_VALID}
			saved_floor = f
		case b := <-button_chan:
			fmt.Println("Button registrered at floor: (Floor, Button_dir)", b.Floor, b.Button_dir)
			local_order_chan <- Order{"LIGHT", b.Button_dir, b.Floor, ON}
			if (saved_floor - b.Floor) == 0 {
				local_order_chan <- Order{"DIRN", DIRN_STOP, b.Floor, NOT_VALID}
			} else if (saved_floor - b.Floor) < 0 {
				local_order_chan <- Order{"DIRN", DIRN_UP, b.Floor, NOT_VALID}
			} else {
				local_order_chan <- Order{"DIRN", DIRN_DOWN, b.Floor, NOT_VALID}
			}
		}
	}
}*/

func sending_stuff(local_order_chan chan<- Order) {
	time.Sleep(3 * time.Second)

	for i := 0; i < N_FLOORS; i++ {
		fmt.Println("LIGHT, BUTTON_CALL_UP, ", i, ", ON")
		local_order_chan <- Order{"LIGHT", BUTTON_CALL_UP, i, ON}
		fmt.Println("LIGHT, BUTTON_CALL_DOWN, ", i, ", ON")
		local_order_chan <- Order{"LIGHT", BUTTON_CALL_DOWN, i, ON}
		fmt.Println("LIGHT, BUTTON_COMMAND, ", i, ", ON")
		local_order_chan <- Order{"LIGHT", BUTTON_COMMAND, i, ON}
	}
	time.Sleep(5 * time.Second)
	//if Elev_get_button_signal(BUTTON_COMMAND, floor) == 1 {
	//button_chan <- Button{floor, BUTTON_COMMAND}
	// Endre "floor", se notat.txt for ekstra. (evt floor_number)
	//}

	for i := 0; i < N_FLOORS; i++ {
		fmt.Println("LIGHT, BUTTON_CALL_UP, ", i, ", OFF")
		local_order_chan <- Order{"LIGHT", BUTTON_CALL_UP, i, OFF}
		fmt.Println("LIGHT, BUTTON_CALL_DOWN, ", i, ", OFF")
		local_order_chan <- Order{"LIGHT", BUTTON_CALL_DOWN, i, OFF}
		fmt.Println("LIGHT, BUTTON_COMMAND, ", i, ", OFF")
		local_order_chan <- Order{"LIGHT", BUTTON_COMMAND, i, OFF}
	}
	time.Sleep(5 * time.Second)

	fmt.Println("LIGHT, DOOR, NOT_VALID, OFF (Set open door lamp OFF)")
	local_order_chan <- Order{"DOOR", 0, NOT_VALID, OFF}
	time.Sleep(5 * time.Second)

	fmt.Println("LIGHT, DOOR, NOT_VALID, ON (Set open door lamp ON)")
	local_order_chan <- Order{"DOOR", 0, NOT_VALID, ON}
	time.Sleep(5 * time.Second)

	fmt.Println("LIGHT, DOOR, NOT_VALID, OFF (Set open door lamp OFF)")
	local_order_chan <- Order{"DOOR", 0, NOT_VALID, ON}
	time.Sleep(5 * time.Second)

	fmt.Println("LIGHT, BUTTON_CALL_UP, 0, ON (Set button lamp ON 1st floor)")
	local_order_chan <- Order{"LIGHT", BUTTON_CALL_UP, 0, ON}
	time.Sleep(5 * time.Second)

	fmt.Println("LIGHT, BUTTON_CALL_UP, 0, OFF (Set button lamp OFF 1st floor (same one))")
	local_order_chan <- Order{"LIGHT", BUTTON_CALL_UP, 0, OFF}
	time.Sleep(5 * time.Second)

	fmt.Println("DIRN, DIRN_UP, NOT_VALID, NOT_VALID ()")
	local_order_chan <- Order{"DIRN", DIRN_UP, NOT_VALID, NOT_VALID}
	time.Sleep(3 * time.Second)

	fmt.Println("DIRN, DIRN_STOP, NOT_VALID, NOT_VALID ()")
	local_order_chan <- Order{"DIRN", DIRN_STOP, NOT_VALID, NOT_VALID}
	time.Sleep(3 * time.Second)

	fmt.Println("LIGHT, DOOR, NOT_VALID, OFF (Set open door lamp OFF)")
	local_order_chan <- Order{"DOOR", 0, NOT_VALID, OFF}
	time.Sleep(5 * time.Second)

}

/*
type Order struct { // Brukes på local_order_chan
	Cat   string // "LIGHT"/"DIR"
	Order int    // DIRN_DOWN/UP/STOP, BUTTON_CALL_UP/DOWN/COMMAND
	Floor int    //0-3 (0-N_FLOORS)
	Value int    // ON/OFF for lys, settes bare for "LIGHT"
} // Floor trengs ikke på doorlight, value trengs ikke på retn.
*/

func checking_for_orders_from_admin(local_order_chan <-chan Order) { // Nytt navn
	for {
		//fmt.Println("In for loop checking orders")
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
		//fmt.Println("In isanythinghappening")
		button_pressed(button_chan)
		//outside_button_pressed(button_chan)
		floor_sensor_triggered(floor_sensor_chan)
		time.Sleep(70 * time.Millisecond) // Least needed to be safe no one is faster to click on button. Crazy people.
	}
}

func button_pressed(button_chan chan<- Button) {
	//fmt.Println("In outside")
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elev_get_button_signal(BUTTON_COMMAND, floor) == 1 {
			button_chan <- Button{floor, BUTTON_COMMAND}
			// Endre "floor", se notat.txt for ekstra. (evt floor_number)
		}
		if Elev_get_button_signal(BUTTON_CALL_UP, floor) == 1 {
			button_chan <- Button{floor, BUTTON_CALL_UP}
		}
		if Elev_get_button_signal(BUTTON_CALL_DOWN, floor) == 1 {
			button_chan <- Button{floor, BUTTON_CALL_DOWN}
		}
	}
}

/*
func inside_button_pressed(button_chan chan<- Button) {
	//fmt.Println("In inside")
	for floor := 0; floor < N_FLOORS; floor++ {
		if Elev_get_button_signal(BUTTON_COMMAND, floor) == 1 {
			button_chan <- Button{floor, BUTTON_COMMAND}
			// Endre "floor", se notat.txt for ekstra. (evt floor_number)
		}
	}
} */

func floor_sensor_triggered(floor_sensor_chan chan<- int) {
	//fmt.Println("In floor sensor")
	floor := Elev_get_floor_sensor_signal()
	if floor != -1 {
		floor_sensor_chan <- floor
	}
}

// NB: Wg.Wait() kan være mulighet/nødvendig en eller annen plass.
