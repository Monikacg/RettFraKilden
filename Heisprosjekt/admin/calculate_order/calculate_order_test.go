package calculate_order

import (
	"fmt"
	. "./../order_matrix"
	. "./../lift_properties"
	. "./../../definitions"
	"testing"
)

func TestFn(t *testing.T) {
	orders := Create_order_matrix()
	properties := Create_lift_prop_list()
	alive_lifts := make([]int, MAX_N_LIFTS)
	ID := 0

	new_dirn, new_dest := Calculate_order(orders, ID, properties, alive_lifts)
	fmt.Println("first order")
	fmt.Println(new_dirn, new_dest)

	if Should_stop(orders, properties, 2, ID) {
		fmt.Println("Should really stop")
	}
}
