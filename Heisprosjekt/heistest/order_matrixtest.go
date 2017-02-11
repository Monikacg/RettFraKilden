package main

import (
  "fmt"
  "../admin/order_matrix"
  "../definitions"
)


func main()  {
  // Create order matrix test
  orders := order_matrix.Create_order_matrix()
  fmt.Println("Create order matrix test: ")
  fmt.Println(orders)


  // Add order test: adds orders on all possible.
  //func add_order(orders [][]int, floor, lift, button_call int)
  for i := 0; i < N_FLOORS; i++ {
    for j := 0; j < MAX_N_LIFTS; j++ {
      for k := 0; k < 3; k++ {
        order_matrix.Add_order(orders, i, j, k)
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
        order_matrix.Delete_order(orders, 1, 2, 0)
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
        order_matrix.Assign_order(orders, i, j, k)
      }
    }
  }
  fmt.Println("Assign order test: ")
  fmt.Println(orders)

  /*
  // Deassign order test
  // func deassign_orders(orders [][]int, lift int)
  for i := 0; i < N_FLOORS; i++ {
    for j := 0; j < MAX_N_LIFTS; j++ {
      for k := 0; k < 3; k++ {
        order_matrix.Deassign_orders(orders, 0)
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
        order_matrix.Complete_order(orders, 1, 2)
      }
    }
  }
  fmt.Println("Complete order test: ")
  fmt.Println(orders)


}
