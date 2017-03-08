package lift_properties

import (
	"fmt"
	"testing"

	. "./../../definitions"
)

func TestFn(t *testing.T) {
	fmt.Println("Create lift properties list: ")
	properties := Create_lift_prop_list()
	fmt.Println(properties)

	dirn := DIRN_STOP
	state := DOOR_OPEN
	last_floor := 1

	fmt.Println("Set last floor = 1 for all lifts: ")
	for lift := 0; lift < MAX_N_LIFTS; lift++ {
		Set_last_floor(properties, lift, last_floor)
	}
	fmt.Println(properties)

	fmt.Println("Set dirn = DIRN_STOP (0) for all lifts: ")
	for lift := 0; lift < MAX_N_LIFTS; lift++ {
		Set_dirn(properties, lift, dirn)
	}
	fmt.Println(properties)

	fmt.Println("Set state = DOOR_OPEN (2) for all lifts: ")
	for lift := 0; lift < MAX_N_LIFTS; lift++ {
		Set_state(properties, lift, state)
	}
	fmt.Println(properties)

}
