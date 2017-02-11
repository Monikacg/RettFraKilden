package lift_properties

import (
  //"definitions"
)

// Til definitions fungerer
const (
  // Lift states
  INIT = -1
  IDLE = 0
  MOVING = 1
  DOOR_OPEN = 2

  type Properties_struct struct { // Brukes i get_properties. ENDRE NAVN!
    last_floor int
    dirn int
    state int
  }
)

func Create_lift_prop_list() ([]int) {
  prop_list := make([]int, 3*MAX_N_LIFTS)
  for i := 0; i < MAX_N_LIFTS; i++ {
     prop_list[3*i] = -1 // Last floor (not valid)
     prop_list[3*i+1] = -2 // Direction (not valid)
     prop_list[3*i+2] = INIT // State (INIT)
   }
  return prop_list
}

func Set_last_floor(properties []int, lift, last_floor int)  {
  properties[3*lift] = last_floor
}

func Set_dirn(properties []int, lift, dirn int)  {
  properties[3*lift+1] = dirn
}

func Set_state(properties []int, lift, state int)  {
  properties[3*lift+2] = state
}

func Get_properties(properties []int, lift int) (properties_struct) {
  return properties_struct{last_floor: properties[3*lift], dirn: properties[3*lift+1], state: properties[3*lift+2]}
}
