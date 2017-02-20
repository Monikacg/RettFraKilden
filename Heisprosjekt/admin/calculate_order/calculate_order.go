package calculate_order

import (
	. "./../../definitions"
	. "./../lift_properties"
)


func Calculate_order(orders[][]int, ID int, properties[]int, alive_lifts []int) (int, int) {
	var new_dirn, dest int = NOT_VALID, NOT_VALID
	dest = find_destination(orders, ID, properties, alive_lifts) //get destination?
	new_dirn = get_new_direction(ID, properties, dest)
	return new_dirn, dest
}

func get_new_direction(ID int, properties[]int, dest int) int {
	if dest - Get_last_floor(properties, ID) > 0 { // properties[3*ID] er last_floor. Endre til 2*ID hvis State tas bort
		return DIRN_UP
	} else if dest - Get_last_floor(properties, ID) < 0 {
		return DIRN_DOWN
	} else {
		return DIRN_STOP
	}
}

func find_destination(orders[][]int, ID int, properties[]int, alive_lifts []int) int {
	dest, dest_exists := check_if_valid_destination_exists(orders, ID)
	//^does not take care about where you are. Needs only 1 order in the system (alt, 1 floor). If not, will go from 0 to 3?
	if dest_exists {
		return dest
	}
	return new_destination(orders, ID, properties, alive_lifts)

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


func new_destination(orders[][]int, ID int, properties[]int, alive_lifts []int) int {
	var new_dest int = NOT_VALID
	var new_dest_exists, i_am_closest bool = false, false
 // Sjekk om skal av i samme etasje.
	switch Get_state(properties, ID) { //Needs to know which elevators are alive
	case DOOR_OPEN:

		switch Get_dirn(properties, ID) {
		case DIRN_UP:
			if order_at_current_floor_moving(orders, properties, ID) {
				return Get_last_floor(properties, ID)
			}
			new_dest, new_dest_exists = order_above(orders, properties, ID)
			if new_dest_exists {
				return new_dest
			}
			// None over, changing direction
			new_dest, new_dest_exists = order_below(orders, properties, ID)
			if new_dest_exists {
				return new_dest
			}
		case DIRN_DOWN:
			if order_at_current_floor_moving(orders, properties, ID) {
				return Get_last_floor(properties, ID)
			}
			new_dest, new_dest_exists = order_below(orders, properties, ID)
			if new_dest_exists {
				return new_dest
			}
			// None over, changing direction
			new_dest, new_dest_exists = order_above(orders, properties, ID)
			if new_dest_exists {
				return new_dest
			}
		}

	case IDLE:
		// NB! Dette kan føre til at flere tar samme. Bør endres?
		if order_at_current_floor_idle(orders, properties, ID) {
			return Get_last_floor(properties, ID)
		}
		new_dest, i_am_closest = am_i_closest_to_new_order(orders, properties, alive_lifts, ID)
		//Sjekker hvilke andre (som er i live) som er IDLE,
		//finner ut hvem som er nærmest. Lavest ID prioritet etter lavest avstand
		//Hvis en nærmere, tar nest nærmest frem til ingen igjen.

		if i_am_closest {
			return new_dest
		}
	}
	return NOT_VALID
}

/*
Sjekke ytre knapper for ordre, indre i alle IDLE
Vil bare vær 1 knapp trykket

MÅ TESTES GRUNDIG
*/
func am_i_closest_to_new_order(orders[][]int, properties[]int, alive_lifts []int, ID int) (int, bool) {
	var closest_lift, new_dest, shortest_distance int = NOT_VALID, NOT_VALID, N_FLOORS+1
	var lf []int

	for floor := 0; floor < N_FLOORS; floor++ {
		if orders[BUTTON_CALL_UP][floor] == 0 {
			new_dest = floor
		}
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			new_dest = floor
		}
	}

	for _, lift := range alive_lifts {
		lf = append(lf, Get_last_floor(properties, lift))
		for floor := 0; floor < N_FLOORS; floor++ {
			if new_dest == NOT_VALID && orders[BUTTON_COMMAND+lift][floor] == 0 {
				if lift == ID {
					return floor, true
				} else {
					return NOT_VALID, false
				}
			}
		}
	}

	// Gives priority to lowest ID. REQUIRES SAME ORDER alive_lifts IN ALL
	// (SORT FROM LOWEST TO HIGHEST?)
	for _, lift := range alive_lifts {
		if abs(Get_last_floor(properties, lift)-new_dest) < shortest_distance {
			closest_lift = lift
		}
	}

	if closest_lift == ID {
		return new_dest, true
	}
	return NOT_VALID, false
}

func abs(value int) int {
	if value < 0 {
		return value *(-1)
	}
	return value
}


// NB! Nå gir den prioritet til de som går ned i høyere etasje over å
// gå ned og hente ny. Endre hvis FAT krever annet.
func order_above(orders [][]int, properties []int, ID int) (int, bool) {
	floor_start := Get_last_floor(properties, ID)+1
	if floor_start >= N_FLOORS {
		return NOT_VALID, false
	}

	for floor := floor_start; floor < N_FLOORS; floor++ {
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return floor, true
		}
		if orders[BUTTON_CALL_UP][floor] == 0 {
			return floor, true
		}
	}
	for floor := floor_start; floor < N_FLOORS; floor++ {
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			return floor, true
		}
	}
	return NOT_VALID, false
}

func order_below(orders [][]int, properties []int, ID int) (int, bool) {
	floor_start := Get_last_floor(properties, ID)-1
	if floor_start < 0 {
		return NOT_VALID, false
	}
	for floor := floor_start; floor < N_FLOORS; floor++ {
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return floor, true
		}
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			return floor, true
		}
	}
	for floor := floor_start; floor < N_FLOORS; floor++ {
		if orders[BUTTON_CALL_UP][floor] == 0 {
			return floor, true
		}
	}
	return NOT_VALID, false
}

//Endre navn sikkert
func order_at_current_floor_moving(orders [][]int, properties []int, ID int) bool {
	floor := Get_last_floor(properties, ID)

	switch Get_dirn(properties, ID) {
	case DIRN_UP:
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return true
		}
		if orders[BUTTON_CALL_UP][floor] == 0 {
			return true
		}
	case DIRN_DOWN:
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return true
		}
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			return true
		}
	}
	return false
}


func order_at_current_floor_idle(orders [][]int, properties []int, ID int) bool {
	floor := Get_last_floor(properties, ID)

	if orders[BUTTON_COMMAND+ID][floor] == 0 {
		return true
	}
	if orders[BUTTON_CALL_UP][floor] == 0 {
		return true
	}
	if orders[BUTTON_CALL_DOWN][floor] == 0 {
		return true
	}
	return false
}





func Should_stop(orders[][]int, properties []int, floor int, ID int) bool {
	if assigned_order_exists(orders, floor, ID) {
		return true
	}
	if unassigned_order_exists(orders, properties, floor, ID) { // En vi skal stoppe på. Feasible unassigned?
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
// Vurder om bør bruk listenotasjon istedenfor å ta inn fra lift_properties
func unassigned_order_exists(orders[][]int, properties []int, floor int, ID int) bool {
	switch Get_dirn(properties, ID) {
	case DIRN_UP:
		if orders[BUTTON_CALL_UP][floor] == 0 {
			return true
		}
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return true
		}
		if floor == N_FLOORS { // Trengs den her egentlig? Vil du egentlig kjøre til 4., komme inn i funksjonen her
			if orders[BUTTON_CALL_DOWN][floor] == 0 { // og likevel komme ned hit? Står inntil videre
				return true
			}
		}
	case DIRN_DOWN:
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			return true
		}
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return true
		}
		if floor == 0 { // Se over
			if orders[BUTTON_CALL_UP][floor] == 0 {
				return true
			}
		}
	}
	return false
}
