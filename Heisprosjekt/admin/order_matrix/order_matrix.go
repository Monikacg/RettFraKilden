package order_matrix



import (
  "fmt"
  "definitions"
)

func Create_order_matrix() ([][]int) {

  order_matrix := make([][]int, 2+MAX_N_LIFTS)
  for i := 0; i < MAX_N_LIFTS; i++ {
     order_matrix[i] = make([]int, N_FLOORS)
     for j := 0; j < N_FLOORS; j++ {
       order_matrix[i][j] = -1
     }
   }
   return order_matrix // The 4 number is N_FLOORS from elev.h, can't import. Find a way to do it to improve.
}

func Add_order(orders [][]int, floor, lift, button_call int)  {
  switch button_call {
  case BUTTON_CALL_UP:
    orders[BUTTON_CALL_UP][floor] = 0 // Index [0][]
  case BUTTON_CALL_DOWN:
    orders[BUTTON_CALL_DOWN][floor] = 0 // Index [1][]
  case BUTTON_COMMAND:
    orders[BUTTON_COMMAND+lift][floor] = 0 // Index[2+lift][]
  }
}

func Delete_order(orders [][]int, floor, lift, button_call int)  {
  switch button_call {
  case BUTTON_CALL_UP:
    orders[BUTTON_CALL_UP][floor] = -1 // Index [0][]
  case BUTTON_CALL_DOWN:
    orders[BUTTON_CALL_DOWN][floor] = -1 // Index [1][]
  case BUTTON_COMMAND:
    orders[BUTTON_COMMAND+lift][floor] = -1 // Index[2+lift][]
  }
}

func Assign_order(orders [][]int, floor, lift, button_call int)   { // NB! BØR LEGGE TIL RETURVERDI SOM INDIKERER OM VI FIKK ASSIGNA
  switch button_call {
  case BUTTON_CALL_UP:
    if orders[BUTTON_CALL_UP][floor] == 0 {
      orders[BUTTON_CALL_UP][floor] = lift+1 // Index [0][]
    }
  case BUTTON_CALL_DOWN:
    if orders[BUTTON_CALL_DOWN][floor] == 0 {
      orders[BUTTON_CALL_DOWN][floor] = lift+1 // Index [1][]
    }
  case BUTTON_COMMAND:
    if orders[BUTTON_COMMAND+lift][floor] == 0 {
      orders[BUTTON_COMMAND+lift][floor] = lift+1 // Index[2+lift][]
    }
  }
}

func Deassign_orders(orders [][]int, lift int)  { // Hvis mister nett -> noen andre skal ta over.
  for floor := 0; floor < 4; floor++ { //SETT INN N_FLOORS
    if orders[BUTTON_CALL_UP][floor] == lift+1 {
      orders[BUTTON_CALL_UP][floor] = 0
    }
    if orders[BUTTON_CALL_DOWN][floor] == lift+1 {
      orders[BUTTON_CALL_DOWN][floor] = 0
    }
    if orders[BUTTON_COMMAND+lift][floor] == lift+1 {
      orders[BUTTON_COMMAND+lift][floor] = 0
    }
  }
}

func Complete_order(orders [][]int, floor, lift int)  { // Akkurat nå har med alle som er utenfor
  orders[BUTTON_CALL_UP][floor] = -1 // Index [0][]
  orders[BUTTON_CALL_DOWN][floor] = -1 // Index [1][]
  orders[BUTTON_COMMAND+lift][floor] = -1 // Index[2+lift][]
}

// Trengs en funksjon som tar vekk én assigned order? Trur det! Reavaluate orders hver gang når etasje/knappetrykk.
