package calculate_order

import (
	"fmt"
	. "./../order_matrix"
	. "./../lift_properties"
	"testing"
)

func TestFn(t *testing.T) {
	orders := Create_order_matrix()
	properties := Create_lift_prop_list()
	ID := 0

	new_dirn, new_dest := Calculate_order(orders,ID, properties)
	fmt.Println("first order")
	fmt.Println(new_dirn, new_dest)

	if Should_stop(orders, 2, ID) {
		fmt.Println("Should really stop")
	}
}
