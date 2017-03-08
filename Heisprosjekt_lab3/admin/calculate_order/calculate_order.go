package calculate_order

import (
	"fmt"

	. "../../definitions"
	. "../lift_properties"
)

func CalculateNextOrder(orders [][]int, ID int, properties []int, alive_lifts []int) (int, int) {
	var newDirn, dest int = NOT_VALID, NOT_VALID
	dest = findDestination(orders, ID, properties, alive_lifts) //get destination?
	newDirn = GetNewDirection(dest, Get_last_floor(properties, ID))
	return newDirn, dest
}

func GetNewDirection(dest int, currentFloor int) int {
	if dest == NOT_VALID {
		return NOT_VALID
	}
	if dest-currentFloor > 0 {
		return DIRN_UP
	} else if dest-currentFloor < 0 {
		return DIRN_DOWN
	} else {
		return DIRN_STOP
	}
}

func findDestination(orders [][]int, ID int, properties []int, alive_lifts []int) int {
	fmt.Println("CalcO: findDest")
	dest, destExists := checkForValidDestination(orders, ID)
	//^does not care about where you are. Needs only 1 order in the system (alt, 1 floor). If not, will go from 0 to 3?
	if destExists {
		return dest
	}
	return newDestination(orders, ID, properties, alive_lifts)

}

func checkForValidDestination(orders [][]int, ID int) (int, bool) {
	fmt.Println("CalcO: checkIfValidDest")
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

func newDestination(orders [][]int, ID int, properties []int, alive_lifts []int) int {
	fmt.Println("CalcO: newDest")
	var newDest int = NOT_VALID
	var newDestExists, iAmClosest bool = false, false
	// Sjekk om skal av i samme etasje.
	switch Get_state(properties, ID) { //Needs to know which elevators are alive
	case DOOR_OPEN:

		switch Get_dirn(properties, ID) {
		case DIRN_UP:
			fmt.Println("CalcO: newDest/DIRN_UP")
			if orderCurrentFloorRightDirection(orders, properties, ID) {
				return Get_last_floor(properties, ID)
			}
			newDest, newDestExists = orderAbove(orders, properties, ID)
			if newDestExists {
				return newDest
			}

			if orderCurrentFloorOppositeDirection(orders, properties, ID) {
				return Get_last_floor(properties, ID)
			}

			// None over, changing direction
			newDest, newDestExists = orderBelow(orders, properties, ID)
			if newDestExists {
				return newDest
			}
		case DIRN_DOWN:
			fmt.Println("CalcO: newDest/DIRN_DOWN")
			if orderCurrentFloorRightDirection(orders, properties, ID) {
				return Get_last_floor(properties, ID)
			}
			newDest, newDestExists = orderBelow(orders, properties, ID)
			if newDestExists {
				return newDest
			}

			if orderCurrentFloorOppositeDirection(orders, properties, ID) {
				return Get_last_floor(properties, ID)
			}
			// None over, changing direction
			newDest, newDestExists = orderAbove(orders, properties, ID)
			if newDestExists {
				return newDest
			}
		}

	case IDLE:
		// NB! Dette kan føre til at flere tar samme. Bør endres?
		if orderCurrentFloorAny(orders, properties, ID) {
			return Get_last_floor(properties, ID)
		}
		newDest, iAmClosest = amIClosestToNewOrder(orders, properties, alive_lifts, ID)
		//Sjekker hvilke andre (som er i live) som er IDLE,
		//finner ut hvem som er nærmest. Lavest ID prioritet etter lavest avstand
		//Hvis en nærmere, tar nest nærmest frem til ingen igjen.

		if iAmClosest {
			return newDest
		}
	}
	return NOT_VALID
}

/*
Sjekke ytre knapper for ordre, indre i alle IDLE
Vil bare vær 1 knapp trykket

MÅ TESTES GRUNDIG
*/
func amIClosestToNewOrder(orders [][]int, properties []int, alive_lifts []int, ID int) (int, bool) {
	fmt.Println("CalcO: amIClosest")
	var closestLift, newDest, shortestDistance int = NOT_VALID, NOT_VALID, N_FLOORS + 1
	var lf []int

	for floor := 0; floor < N_FLOORS; floor++ {
		if orders[BUTTON_CALL_UP][floor] == 0 {
			newDest = floor
		}
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			newDest = floor
		}
	}

	//Place all Idle lifts in a slice and iterate over them instead of alive_lifts
	for _, lift := range alive_lifts {
		lf = append(lf, Get_last_floor(properties, lift))
		for floor := 0; floor < N_FLOORS; floor++ {
			if newDest == NOT_VALID && orders[BUTTON_COMMAND+lift][floor] == 0 {
				if lift == ID {
					return floor, true
				}
			}
		}
	}

	// Gives priority to lowest ID. REQUIRES SAME ORDER alive_lifts IN ALL
	// (SORT FROM LOWEST TO HIGHEST?)
	// Use IDLE list made above.
	for _, lift := range alive_lifts {
		if abs(Get_last_floor(properties, lift)-newDest) < shortestDistance {
			shortestDistance = abs(Get_last_floor(properties, lift) - newDest)
			closestLift = lift
		}
	}

	if closestLift == ID {
		return newDest, true
	}
	return NOT_VALID, false
}

func abs(value int) int {
	if value < 0 {
		return value * (-1)
	}
	return value
}

// NB! Nå gir den prioritet til de som går ned i høyere etasje over å
// gå ned og hente ny. Endre hvis FAT krever annet.
func orderAbove(orders [][]int, properties []int, ID int) (int, bool) {
	fmt.Println("CalcO: orderAbove")
	floor_start := Get_last_floor(properties, ID) + 1
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

func orderBelow(orders [][]int, properties []int, ID int) (int, bool) {
	fmt.Println("CalcO: orderBelow")
	floor_start := Get_last_floor(properties, ID) - 1
	if floor_start < 0 {
		return NOT_VALID, false
	}
	for floor := floor_start; floor >= 0; floor-- {
		fmt.Println("CalcO: orderBelow/floorloop: floor_start, floor: ", floor_start, floor)
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return floor, true
		}
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			return floor, true
		}
	}
	for floor := floor_start; floor >= 0; floor-- {
		if orders[BUTTON_CALL_UP][floor] == 0 {
			return floor, true
		}
	}
	return NOT_VALID, false
}

//Endre navn sikkert
func orderCurrentFloorRightDirection(orders [][]int, properties []int, ID int) bool {
	fmt.Println("CalcO: orderCurrentFloorMoving")
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

func orderCurrentFloorOppositeDirection(orders [][]int, properties []int, ID int) bool {
	fmt.Println("CalcO: orderCurrentFloorWrongWay")
	floor := Get_last_floor(properties, ID)

	switch Get_dirn(properties, ID) {
	case DIRN_UP:
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			return true
		}

	case DIRN_DOWN:
		if orders[BUTTON_CALL_UP][floor] == 0 {
			return true
		}
	}
	return false
}

func orderCurrentFloorAny(orders [][]int, properties []int, ID int) bool {
	fmt.Println("CalcO: orderCurrentFloorIdle")
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

func ShouldStop(orders [][]int, properties []int, floor int, ID int) bool {
	fmt.Println("CalcO: ShouldStop, floor", floor)

	//Test at fungerer
	if floor == 0 && Get_dirn(properties, ID) == DIRN_DOWN {
		return true
	}
	if floor == N_FLOORS && Get_dirn(properties, ID) == DIRN_UP {
		return true
	}
	if assignedOrderExists(orders, floor, ID) {
		return true
	}
	if unassignedOrderExists(orders, properties, floor, ID) { // En vi skal stoppe på. Feasible unassigned?
		return true
	}
	return false
}

func assignedOrderExists(orders [][]int, floor int, ID int) bool {
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
func unassignedOrderExists(orders [][]int, properties []int, floor int, ID int) bool {
	switch Get_dirn(properties, ID) {
	case DIRN_UP:
		if orders[BUTTON_CALL_UP][floor] == 0 {
			return true
		}
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return true
		}
		/*
			if floor == N_FLOORS { // Trengs den her egentlig? Vil du egentlig kjøre til 4., komme inn i funksjonen her
				if orders[BUTTON_CALL_DOWN][floor] == 0 { // og likevel komme ned hit? Står inntil videre
					return true
				}
			}
		*/
	case DIRN_DOWN:
		if orders[BUTTON_CALL_DOWN][floor] == 0 {
			return true
		}
		if orders[BUTTON_COMMAND+ID][floor] == 0 {
			return true
		}
		/*
			if floor == 0 { // Se over
				if orders[BUTTON_CALL_UP][floor] == 0 {
					return true
				}
			}
		*/
	}
	return false
}
