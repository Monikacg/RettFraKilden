package lift_properties

import (
	. "./../../definitions"
)

func Create_lift_prop_list() []int {
	prop_list := make([]int, 3*MAX_N_LIFTS)
	for i := 0; i < MAX_N_LIFTS; i++ {
		prop_list[3*i] = NOT_VALID   // Last floor
		prop_list[3*i+1] = NOT_VALID // Direction
		prop_list[3*i+2] = INIT      // State
	}
	return prop_list
}

func Set_last_floor(properties []int, lift, last_floor int) {
	properties[3*lift] = last_floor
}

func Set_dirn(properties []int, lift, dirn int) {
	properties[3*lift+1] = dirn
}

func Set_state(properties []int, lift, state int) {
	properties[3*lift+2] = state
}

/* Slett hvis ikke blir brukt
func Get_properties(properties []int, lift int) Properties_struct {
	return Properties_struct{Last_floor: properties[3*lift], Dirn: properties[3*lift+1], State: properties[3*lift+2]}
} // Tror ikke den her trengs/skal brukes. Sender hel tabell når slår opp.
*/

func Get_last_floor(properties []int, lift int) int { // I calculate_order
	return properties[3*lift]
}

func Get_dirn(properties []int, lift int) int { // I calculate_order
	return properties[3*lift+1]
}

func Get_state(properties []int, lift int) int { // I admin, calculate_order
	return properties[3*lift+2]
}
