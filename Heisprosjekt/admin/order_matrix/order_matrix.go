package order_matrix



import (
  //"fmt"
  //"definitions"
)

//Til definitions kan bli lagt til
const (
  BUTTON_CALL_UP = 0
  BUTTON_CALL_DOWN = 1
  BUTTON_COMMAND = 2

  N_FLOORS = 4
  N_BUTTONS = 3
  MAX_N_LIFTS = 3

)

func create_order_matrix() ([][]int) {

  order_matrix := make([][]int, 2+MAX_N_LIFTS)
  for i := 0; i < 2+MAX_N_LIFTS; i++ {
     order_matrix[i] = make([]int, N_FLOORS)
     for j := 0; j < N_FLOORS; j++ {
       order_matrix[i][j] = -1
     }
   }
   return order_matrix // The 4 number is N_FLOORS from elev.h, can't import. Find a way to do it to improve.
}

func add_order(orders [][]int, floor, lift, button_call int)  {
  switch button_call {
  case BUTTON_CALL_UP:
    orders[BUTTON_CALL_UP][floor] = 0 // Index [0][]
  case BUTTON_CALL_DOWN:
    orders[BUTTON_CALL_DOWN][floor] = 0 // Index [1][]
  case BUTTON_COMMAND:
    orders[BUTTON_COMMAND+lift][floor] = 0 // Index[2+lift][]
  }
}

func delete_order(orders [][]int, floor, lift, button_call int)  {
  switch button_call {
  case BUTTON_CALL_UP:
    orders[BUTTON_CALL_UP][floor] = -1 // Index [0][]
  case BUTTON_CALL_DOWN:
    orders[BUTTON_CALL_DOWN][floor] = -1 // Index [1][]
  case BUTTON_COMMAND:
    orders[BUTTON_COMMAND+lift][floor] = -1 // Index[2+lift][]
  }
}

func assign_order(orders [][]int, floor, lift, button_call int)   { // NB! BØR LEGGE TIL RETURVERDI SOM INDIKERER OM VI FIKK ASSIGNA
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

func deassign_orders(orders [][]int, lift int)  { // Hvis mister nett -> noen andre skal ta over.
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

func complete_order(orders [][]int, floor, lift int)  { // Akkurat nå har med alle som er utenfor
  orders[BUTTON_CALL_UP][floor] = -1 // Index [0][]
  orders[BUTTON_CALL_DOWN][floor] = -1 // Index [1][]
  orders[BUTTON_COMMAND+lift][floor] = -1 // Index[2+lift][]
}

// Trengs en funksjon som tar vekk én assigned order? Trur det! Reavaluate orders hver gang når etasje/knappetrykk.
