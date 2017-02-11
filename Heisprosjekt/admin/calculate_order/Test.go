package main

import(
	"fmt"
)

const (
  BUTTON_CALL_UP = 0
  BUTTON_CALL_DOWN = 1
  BUTTON_COMMAND = 2

  N_FLOORS = 4
  N_BUTTONS = 3
  MAX_N_LIFTS = 3

)

func main() {
	order_matrix := make([][]int, 2+MAX_N_LIFTS)
		for i := 0; i < 2+MAX_N_LIFTS; i++ {
			order_matrix[i] = make([]int, N_FLOORS)
			for j := 0; j < N_FLOORS; j++ {
				order_matrix[i][j] = -1
		}
	}

	first := calculate_order(order_matrix, 0)
	fmt.Println("first order")
	fmt.Println(first)
}