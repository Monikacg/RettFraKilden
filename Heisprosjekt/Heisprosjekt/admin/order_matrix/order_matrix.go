package order_matrix
/*
#cgo CFLAGS: -std=gnu11
#cgo LDFLAGS: -lcomedi -lm
#include "elev.h"
*/
import "C"

func create_order_matrix()  {
  orders := [5][4]int{} // uses 3 elevators, could use down+up+n elevators, room for improvement
  // The 4 number is N_FLOORS from elev.h, can't import. Find a way to do it to improve.
}

func add_order_to_matrix(floor, elev, direction)  {

}
