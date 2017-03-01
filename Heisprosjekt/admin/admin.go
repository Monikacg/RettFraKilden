package admin

import (
	"sort"

	. "./../definitions"
	. "./calculate_order"
	. "./lift_properties"
	. "./order_matrix"
)

// Mulig å slå sammen som interface{}?

func Admin_init(IDInput int, button_chan <-chan Button, floor_sensor_chan <-chan int,
	local_order_chan chan<- Order, adm_transmit_chan chan<- Udp, adm_receive_chan <-chan Udp, peer_chan <-chan Peer,
	start_timer_chan chan<- string, time_out_chan <-chan string) {

	go admin(IDInput, button_chan, floor_sensor_chan,
		local_order_chan, adm_transmit_chan, adm_receive_chan, peer_chan,
		start_timer_chan, time_out_chan)
}

func admin(IDInput int, button_chan <-chan Button, floor_sensor_chan <-chan int,
	local_order_chan chan<- Order, adm_transmit_chan chan<- Udp, adm_receive_chan <-chan Udp, peer_chan <-chan Peer,
	start_timer_chan chan<- string, time_out_chan <-chan string) {

	orders := Create_order_matrix()
	properties := Create_lift_prop_list()
	ID := IDInput //SETT ID (spør nett?)
	// Check om orders/prop list eksisterer noen andre plasser på nettet
	alive_lifts := make([]int, 0, MAX_N_LIFTS)
	alive_lifts = append(alive_lifts, ID)

	for {
		select {
		// Problem med å sende melding om button pressed ut på nettet og deretter melding fra find_new_order?
		// evt legge ved hvilke ordre vi tar hver gang i find_new_order-melding => alle andre kan oppdatere.
		// Husk "problem" med at assign bare tar de som allerede finnes, så
		// må ha en måte å slå sammen her.

		case b := <-button_chan: // INSIDE AND OUTSIDE
			//If order already exists (legg til funksjon i order_matrix), sett
			// legg ny order inn midlertidig plass til får melding fra network om alene
			// Eller: Kan jo sende til NW, få tilbake, så legge til hvis
			// ikke alene. Legge til uansett hvis indre.
			if b.Button_dir == BUTTON_COMMAND {
				Add_order(orders, b.Floor, ID, b.Button_dir) // Inside order
			} else if len(alive_lifts) > 1 {
				Add_order(orders, b.Floor, ID, b.Button_dir) // Outside order. Tas bare når vi vet at andre heiser eksisterer.
			}
			//Send melding ut til NW på adm_transmit_chan
			adm_transmit_chan <- Udp{ID, "ButtonPressed", b.Floor, b.Button_dir}

			//Tanke: Legg inn noe som gjør at det ikke legges til(sendes ut på NW) hvis allerede finnes i orders.
			/*if not in orders {
				adm_transmit_chan <- Udp{ID, "ButtonPressed", b.Floor, b.Button_dir}
			}*/

			if Get_state(properties, ID) == IDLE {
				find_new_order(orders, ID, properties, alive_lifts, start_timer_chan, local_order_chan, adm_transmit_chan)
			}

		case fs := <-floor_sensor_chan:
			switch Get_state(properties, ID) {
			case DOOR_OPEN:
				//Intentionally blank, probably might as well just remove this case, right now for completeness
				// Just needs to break, which it will do without these. Maybe a small sleep on these?
			case IDLE:
				// See DOOR_OPEN
			case MOVING:
				Set_last_floor(properties, ID, fs)

				if Should_stop(orders, properties, fs, ID) == true {
					Assign_orders(orders, fs, ID) // In case du tar en som ikke var assigna.
					//local_order_chan <-  Send "DIRN", DIRN_STOP, NOT_VALID, ON
					local_order_chan <- Order{"DIRN", DIRN_STOP, NOT_VALID, ON}
					Set_state(properties, ID, DOOR_OPEN)
					start_timer_chan <- "DOOR_OPEN"
					Complete_order(orders, fs, ID)
					//Send melding ut til NW på adm_transmit_chan
					// ID, "Stoppet", etasje (DOOR_OPEN)
					adm_transmit_chan <- Udp{ID, "Stopped", fs, NOT_VALID}
				} else {
					//Send melding ut til NW på adm_transmit_chan
					// ID, "kjørte forbi", etasje
					adm_transmit_chan <- Udp{ID, "DrovePast", fs, NOT_VALID}
				}
			}

		case <-time_out_chan:
			//local_order_chan <-  Send "DOOR", NOT_VALID, NOT_VALID, OFF
			local_order_chan <- Order{"DOOR", NOT_VALID, NOT_VALID, OFF}
			find_new_order(orders, ID, properties, alive_lifts, start_timer_chan, local_order_chan, adm_transmit_chan)

		case m := <-adm_receive_chan:
			switch m.ID {
			case ID:
				//Alt for egen heis
				switch m.Type {
				case "ButtonPressed":
					Add_order(orders, m.Floor, m.ID, m.ExtraInfo) // Ta bort den over og la den her stå? bedre her her.
					local_order_chan <- Order{"LIGHT", m.ExtraInfo, m.Floor, ON}
				case "Stopped":
					// Får ingenting tilbake fra andre her.
				case "DrovePast":
					// Ingenting
				case "NewOrder":
					// Gjør alt før, er bare ack her. Skal det i det hele tatt komme tilbake hit?
				case "Idle":
					// Samme som over. Nada.
				}

			default: //Any other lift
				switch m.Type {
				case "ButtonPressed":
					Add_order(orders, m.Floor, m.ID, m.ExtraInfo)
				case "Stopped":
					Assign_orders(orders, m.Floor, m.ID)
					Complete_order(orders, m.Floor, m.ID)
					Set_state(properties, m.ID, DOOR_OPEN)
					Set_last_floor(properties, m.ID, m.Floor)

				case "DrovePast":
					Set_last_floor(properties, m.ID, m.Floor)

				case "NewOrder":
					Set_state(properties, m.ID, MOVING)
					Set_dirn(properties, m.ID, Get_new_direction(m.Floor, Get_last_floor(properties, m.ID)))
				case "Idle":
					Set_state(properties, m.ID, IDLE)
				}
			}

		case peer_msg := <-peer_chan:
			switch peer_msg.Change {
			case "New":
				alive_lifts = append(alive_lifts, peer_msg.ChangedPeer)
				sort.Slice(alive_lifts, func(i, j int) bool { return alive_lifts[i] < alive_lifts[j] })
			case "Lost":
				for i, n := range alive_lifts {
					if n == peer_msg.ChangedPeer {
						alive_lifts = append(alive_lifts[:i], alive_lifts[i+1:]...)
					}
				}
			}
		}
	}
}

func find_new_order(orders [][]int, ID int, properties []int, alive_lifts []int, start_timer_chan chan<- string,
	local_order_chan chan<- Order, adm_transmit_chan chan<- Udp) {

	new_dirn, dest := Calculate_order(orders, ID, properties, alive_lifts)
	// Should change name on both module called from and function itself.
	// Default dest and new_dirn returned has to be undefined (-2,-2)
	if new_dirn == DIRN_STOP {
		Assign_orders(orders, dest, ID) //NB! Nå lagt til ALLE på den etasjen,
		// noe som er en forenkling som vi kunne gjøre. IKKE TESTET ENNÅ
		//local_order_chan <-  Send "DOOR", NOT_VALID, NOT_VALID, ON
		local_order_chan <- Order{"DOOR", NOT_VALID, NOT_VALID, ON}
		Set_state(properties, ID, DOOR_OPEN)
		start_timer_chan <- "DOOR_OPEN"
		Complete_order(orders, dest, ID)
		//Send melding ut til NW på adm_transmit_chan
		// ID, "Stoppet", etasje (DOOR_OPEN)
		adm_transmit_chan <- Udp{ID, "Stopped", dest, NOT_VALID}
	} else if new_dirn == DIRN_DOWN || new_dirn == DIRN_UP {
		Assign_orders(orders, dest, ID)
		//local_order_chan <-  Send "DIRN", DIRN_UP/DOWN, NOT_VALID, NOT_VALID
		local_order_chan <- Order{"DIRN", new_dirn, NOT_VALID, NOT_VALID}
		Set_state(properties, ID, MOVING)
		Set_dirn(properties, ID, new_dirn)
		//Send melding ut til NW på adm_transmit_chan
		// ID, "Moving, desting (new order)", etasje
		adm_transmit_chan <- Udp{ID, "NewOrder", dest, NOT_VALID}
	} else { // new_dirn == -2 (NOT_VALID)
		Set_state(properties, ID, IDLE)
		//Send melding ut til NW på adm_transmit_chan
		// ID, "IDLE", etasje
		adm_transmit_chan <- Udp{ID, "Idle", dest, NOT_VALID}
	}
}
