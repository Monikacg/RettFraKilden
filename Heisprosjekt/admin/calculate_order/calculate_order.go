package main

/*
Admin kaller funksjonen calculate_order for å finne ut hva heisen skal gjøre nå.
returnerer hvilken etasje heisen skal kjøre til (og retning? eller finner heisen ut av dette på egenhånd?)

Hvordan finne første ordre? Gå gjennom lista med en for-løkke?
Hva med indre ordre? skal disse letes gjennom først? JA
*/
////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type dummy struct {
	button int
	floor int
}

func find_order(orders [][]int, lift int) dummy { // N_FLOORS
	for f := 0; f < 4; f++ {
		if orders[lift + 2][f] == 0 {
			innside := dummy{button: lift+2, floor: f}
			return innside
		} 
	}
	for f := 0; f < 4; f++ {
		for dir := 0; dir < 2; dir++ {
			if orders[dir][f] == 0{
				outside := dummy{button: dir, floor: f}
				return outside
			}
		}
	}
	non := dummy{button: -1, floor: -1}
	return non
}

func calculate_order(orders [][]int, lift int) int {
	order := find_order(orders, lift)
	return order.floor
}