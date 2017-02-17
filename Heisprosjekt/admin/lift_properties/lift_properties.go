package lift_properties

import (
  . "./../../definitions"
)

func Create_lift_prop_list() ([]int) {
  prop_list := make([]int, 3*MAX_N_LIFTS)
  for i := 0; i < MAX_N_LIFTS; i++ {
     prop_list[3*i] = NOT_VALID // Last floor (not valid) Kan godt ha -2 som ikke gyldig på begge
     prop_list[3*i+1] = NOT_VALID // Direction (not valid) Ide: Legg til def av not valid?
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

func Get_properties(properties []int, lift int) (Properties_struct) {
  return Properties_struct{Last_floor: properties[3*lift], Dirn: properties[3*lift+1], State: properties[3*lift+2]}
} // Tror ikke den her trengs/skal brukes. Sender hel tabell når slår opp.
// Alternativet er å oppdatere en struct hele tiden i tillegg.

func Get_state(properties []int, lift int) int {
  return properties[3*lift+2]
}
