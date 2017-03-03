package admin

import (
	"sort"
	"time"

	. "./../definitions"
	. "./calculate_order"
	. "./lift_properties"
	. "./order_matrix"
)

// Mulig å slå sammen som interface{}?

func Admin_init(IDInput int, buttonChan <-chan Button, floorSensorChan <-chan int,
	localOrderChan chan<- Order, adminTChan chan<- Udp, adminRChan <-chan Udp, backupChan <-chan Backup, peerChangeChan <-chan Peer,
	startTimerChan chan<- string, timeOutChan <-chan string) {

	go admin(IDInput, buttonChan, floorSensorChan,
		localOrderChan, adminTChan, adminRChan, backupChan, peerChangeChan,
		startTimerChan, timeOutChan)
}

func admin(IDInput int, buttonChan <-chan Button, floorSensorChan <-chan int,
	localOrderChan chan<- Order, adminTChan chan<- Udp, adminRChan <-chan Udp, backupChan <-chan Backup, peerChangeChan <-chan Peer,
	startTimerChan chan<- string, timeOutChan <-chan string) {

	orders := Create_order_matrix()
	properties := Create_lift_prop_list()
	ID := IDInput //SETT ID (spør nett?)
	// Check om orders/prop list eksisterer noen andre plasser på nettet
	alive_lifts := make([]int, 0, MAX_N_LIFTS)
	alive_lifts = append(alive_lifts, ID)

	// Enter init state
	Set_state(properties, ID, IDLE)

	//Spør nett om noen har orders og properties. If so, sett orders og properties lik de på nettet. If not, forsett med det samme.

	// Tror det er uavhengig av hvilken state det står i for det her: Vil bare vite om vi står i en etasje.
initLoop:
	for {
		select {
		case f := <-floorSensorChan:
			Set_last_floor(properties, ID, f)
			localOrderChan <- Order{"FLOOR_LIGHT", NOT_VALID, f, ON}
			Assign_orders(orders, f, ID)
			localOrderChan <- Order{"DIRN", DIRN_STOP, NOT_VALID, ON} // Bør sikkert endre navn for å gjør den her enkler å forstå
			startTimerChan <- "DOOR_OPEN"
			Set_state(properties, ID, DOOR_OPEN)
			Complete_order(orders, f, ID)
			adminTChan <- Udp{ID, "Stopped", f, NOT_VALID}
			break initLoop

		case <-time.After(3 * time.Second):
			Set_state(properties, ID, MOVING)
			Set_dirn(properties, ID, DIRN_DOWN)
			localOrderChan <- Order{"DIRN", DIRN_DOWN, NOT_VALID, NOT_VALID}
			break initLoop
		}
	}
	// Exit init state.

	for {
		select {
		// Problem med å sende melding om button pressed ut på nettet og deretter melding fra find_new_order?
		// evt legge ved hvilke ordre vi tar hver gang i find_new_order-melding => alle andre kan oppdatere.
		// Husk "problem" med at assign bare tar de som allerede finnes, så
		// må ha en måte å slå sammen her.

		case b := <-buttonChan: // INSIDE AND OUTSIDE
			//If order already exists (legg til funksjon i order_matrix), sett
			// legg ny order inn midlertidig plass til får melding fra network om alene
			// Eller: Kan jo sende til NW, få tilbake, så legge til hvis
			// ikke alene. Legge til uansett hvis indre.
			if b.Button_dir == BUTTON_COMMAND {
				Add_order(orders, b.Floor, ID, b.Button_dir) // Inside order
			} else if len(alive_lifts) > 1 {
				Add_order(orders, b.Floor, ID, b.Button_dir) // Outside order. Tas bare når vi vet at andre heiser eksisterer.
			}
			//Send melding ut til NW på adminTChan
			adminTChan <- Udp{ID, "ButtonPressed", b.Floor, b.Button_dir}

			//Tanke: Legg inn noe som gjør at det ikke legges til(sendes ut på NW) hvis allerede finnes i orders.
			/*if not in orders {
				adminTChan <- Udp{ID, "ButtonPressed", b.Floor, b.Button_dir}
			}*/

			if Get_state(properties, ID) == IDLE {
				find_new_order(orders, ID, properties, alive_lifts, startTimerChan, localOrderChan, adminTChan)
			}

		case fs := <-floorSensorChan:
			switch Get_state(properties, ID) {
			case DOOR_OPEN:
				//Intentionally blank, probably might as well just remove this case, right now for completeness
				// Just needs to break, which it will do without these. Maybe a small sleep on these?
			case IDLE:
				// See DOOR_OPEN
			case MOVING:
				Set_last_floor(properties, ID, fs)
				localOrderChan <- Order{"FLOOR_LIGHT", NOT_VALID, fs, ON}

				if Should_stop(orders, properties, fs, ID) == true {
					Assign_orders(orders, fs, ID) // In case du tar en som ikke var assigna.
					//localOrderChan <-  Send "DIRN", DIRN_STOP, NOT_VALID, ON
					localOrderChan <- Order{"DIRN", DIRN_STOP, NOT_VALID, ON}
					Set_state(properties, ID, DOOR_OPEN)
					startTimerChan <- "DOOR_OPEN"
					Complete_order(orders, fs, ID)
					//Send melding ut til NW på adminTChan
					// ID, "Stoppet", etasje (DOOR_OPEN)
					adminTChan <- Udp{ID, "Stopped", fs, NOT_VALID}
				} else {
					//Send melding ut til NW på adminTChan
					// ID, "kjørte forbi", etasje
					adminTChan <- Udp{ID, "DrovePast", fs, NOT_VALID}
				}
			}

		case <-timeOutChan:
			//localOrderChan <-  Send "DOOR", NOT_VALID, NOT_VALID, OFF
			localOrderChan <- Order{"DOOR", NOT_VALID, NOT_VALID, OFF}
			find_new_order(orders, ID, properties, alive_lifts, startTimerChan, localOrderChan, adminTChan)

		case m := <-adminRChan:
			switch m.ID {
			case ID:
				//Alt for egen heis
				switch m.Type {
				case "ButtonPressed":
					Add_order(orders, m.Floor, m.ID, m.ExtraInfo) // Ta bort den over og la den her stå? bedre her her.
					localOrderChan <- Order{"LIGHT", m.ExtraInfo, m.Floor, ON}
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

		case peer_msg := <-peerChangeChan:
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

func find_new_order(orders [][]int, ID int, properties []int, alive_lifts []int, startTimerChan chan<- string,
	localOrderChan chan<- Order, adminTChan chan<- Udp) {

	new_dirn, dest := Calculate_order(orders, ID, properties, alive_lifts)
	// Should change name on both module called from and function itself.
	// Default dest and new_dirn returned has to be undefined (-2,-2)
	if new_dirn == DIRN_STOP {
		Assign_orders(orders, dest, ID) //NB! Nå lagt til ALLE på den etasjen,
		// noe som er en forenkling som vi kunne gjøre. IKKE TESTET ENNÅ
		//localOrderChan <-  Send "DOOR", NOT_VALID, NOT_VALID, ON
		localOrderChan <- Order{"DOOR", NOT_VALID, NOT_VALID, ON}
		Set_state(properties, ID, DOOR_OPEN)
		startTimerChan <- "DOOR_OPEN"
		Complete_order(orders, dest, ID)
		//Send melding ut til NW på adminTChan
		// ID, "Stoppet", etasje (DOOR_OPEN)
		adminTChan <- Udp{ID, "Stopped", dest, NOT_VALID}
	} else if new_dirn == DIRN_DOWN || new_dirn == DIRN_UP {
		Assign_orders(orders, dest, ID)
		//localOrderChan <-  Send "DIRN", DIRN_UP/DOWN, NOT_VALID, NOT_VALID
		localOrderChan <- Order{"DIRN", new_dirn, NOT_VALID, NOT_VALID}
		Set_state(properties, ID, MOVING)
		Set_dirn(properties, ID, new_dirn)
		//Send melding ut til NW på adminTChan
		// ID, "Moving, desting (new order)", etasje
		adminTChan <- Udp{ID, "NewOrder", dest, NOT_VALID}
	} else { // new_dirn == -2 (NOT_VALID)
		Set_state(properties, ID, IDLE)
		//Send melding ut til NW på adminTChan
		// ID, "IDLE", etasje
		adminTChan <- Udp{ID, "Idle", dest, NOT_VALID}
	}
}
