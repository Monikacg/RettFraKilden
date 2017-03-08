package calculate_order

import (
	"fmt"
	"testing"

	. "./../../definitions"
	. "./../lift_properties"
	. "./../order_matrix"
)

func TestFn(t *testing.T) {
	orders := Create_order_matrix()
	properties := Create_lift_prop_list()
	aliveLifts := make([]int, MAX_N_LIFTS)
	ID := 0

	new_dirn, new_dest := CalculateNextOrder(orders, ID, properties, aliveLifts)
	fmt.Println("first order")
	fmt.Println(new_dirn, new_dest)

	if ShouldStop(orders, properties, 2, ID) {
		fmt.Println("Should really stop")
	}
}
