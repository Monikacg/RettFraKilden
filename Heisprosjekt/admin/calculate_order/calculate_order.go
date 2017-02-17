package calculate_order

import (
	. "./../../definitions"
)

/*
Admin kaller funksjonen calculate_order for å finne ut hva heisen skal gjøre nå.
returnerer hvilken etasje heisen skal kjøre til (og retning? eller finner heisen ut av dette på egenhånd?)
Hvordan finne første ordre? Gå gjennom lista med en for-løkke?
Hva med indre ordre? skal disse letes gjennom først? JA
*/
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/*
type dummy struct {
	button int
	floor int
}

func find_order(orders [][]int, lift int) dummy { // N_FLOORS
	for f := 0; f < 4; f++ {
		if orders[lift + 2][f] == 0 {
			innside := dummy{button: lift+2, floor: f}
			return innside
		}
	}
	for f := 0; f < 4; f++ {
		for dir := 0; dir < 2; dir++ {
			if orders[dir][f] == 0{
				outside := dummy{button: dir, floor: f}
				return outside
			}
		}
	}
	non := dummy{button: -1, floor: -1}
	return non
}

func calculate_order(orders [][]int, lift int) int {
	order := find_order(orders, lift)
	return order.floor
}

*/










func Calculate_order(orders[][]int, ID int, properties[]int) (int, int) {
	var new_dirn, dest int
	dest = find_destination(orders, ID, properties) //get destination?
	new_dirn = get_new_direction(ID, properties, dest)
	return new_dirn, dest
}

func find_destination(orders[][]int, ID int, properties[]int) int {
	dest, dest_exists := check_if_valid_destination_exists(orders, ID)
	//^does not take care about where you are. Needs only 1 order in the system (alt, 1 floor). If not, will go from 0 to 3?
	if dest_exists {
		return dest
	}
	return new_destination(orders, ID, properties)

}

func new_destination(orders[][]int, ID int, properties[]int) int {
	var new_dest int = NOT_VALID

	return new_dest
}

func get_new_direction(ID int, properties[]int, dest int) int {
	if dest - properties[3*ID] > 0 { // properties[3*ID] er last_floor. Endre til 2*ID hvis State tas bort
		return DIRN_UP
	} else if dest - properties[3*ID] < 0 {
		return DIRN_DOWN
	} else {
		return DIRN_STOP
	}
}

func check_if_valid_destination_exists(orders[][]int, ID int) (int, bool) {
	for floor := 0; floor < N_FLOORS; floor++ {
		if orders[BUTTON_CALL_UP][floor] == ID+1 {
			return floor, true
		}
		if orders[BUTTON_CALL_DOWN][floor] == ID+1 {
			return floor, true
		}
		if orders[BUTTON_COMMAND+ID][floor] == ID+1 {
			return floor, true
		}
	}
	return NOT_VALID, false
}

func Should_stop(orders[][]int, floor int, ID int) bool {
	if assigned_order_exists(orders, floor, ID) {
		return true
	}
	return false
}

func assigned_order_exists(orders[][]int, floor int, ID int) bool {
	if orders[BUTTON_CALL_UP][floor] == ID+1 {
		return true
	}
	if orders[BUTTON_CALL_DOWN][floor] == ID+1 {
		return true
	}
	if orders[BUTTON_COMMAND+ID][floor] == ID+1 {
		return true
	}
	return false
}
