package order_matrix

import (
	"fmt"
	"testing"

	. "./../../definitions"
)

func TestFn(t *testing.T) {
	// Create order matrix test
	orders := Create_order_matrix()
	fmt.Println("Create order matrix test: ")
	fmt.Println(orders)

	// Add order test: adds orders on all possible.
	//func add_order(orders [][]int, floor, lift, button_call int)
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < MAX_N_LIFTS; j++ {
			for k := 0; k < 3; k++ {
				Add_order(orders, i, j, k)
			}
		}
	}
	fmt.Println("Add order test: ")
	fmt.Println(orders)

	/*
	  // Delete order test
	  //func delete_order(orders [][]int, floor, lift, button_call int)
	  for i := 0; i < N_FLOORS; i++ {
	    for j := 0; j < MAX_N_LIFTS; j++ {
	      for k := 0; k < 3; k++ {
	        Delete_order(orders, 1, 2, 0)
	      }
	    }
	  }
	  fmt.Println("Delete order test: ")
	  fmt.Println(orders)
	*/

	// Assign order test
	//func delete_order(orders [][]int, floor, lift, button_call int)
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < MAX_N_LIFTS; j++ {
			for k := 0; k < 3; k++ {
				Assign_order(orders, i, j, k)
			}
		}
	}
	fmt.Println("Assign order test: ")
	fmt.Println(orders)

	/*
	  // Assign orders test
	  //func delete_order(orders [][]int, floor, lift, button_call int)
	  for i := 0; i < N_FLOORS; i++ {
	    for j := 0; j < MAX_N_LIFTS; j++ {
	      Assign_orders(orders, i, j)
	    }
	  }
	  fmt.Println("Assign orders test: ")
	  fmt.Println(orders)
	*/

	/*
	  // Deassign order test
	  // func deassign_orders(orders [][]int, lift int)
	  for i := 0; i < N_FLOORS; i++ {
	    for j := 0; j < MAX_N_LIFTS; j++ {
	      for k := 0; k < 3; k++ {
	        Deassign_orders(orders, 0)
	      }
	    }
	  }
	  fmt.Println("Deassign order test: ")
	  fmt.Println(orders)
	*/

	// Complete order test
	// func complete_order(orders [][]int, floor, lift int)
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < MAX_N_LIFTS; j++ {
			for k := 0; k < 3; k++ {
				Complete_order(orders, 1, 2)
			}
		}
	}
	fmt.Println("Complete order test: ")
	fmt.Println(orders)

}
