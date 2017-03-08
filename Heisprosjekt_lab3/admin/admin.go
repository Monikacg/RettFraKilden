package admin

import (
	"fmt"
	"time"

	. "../definitions"
	. "./calculate_order"
	. "./lift_properties"
	. "./order_matrix"
	"sort"
)

// Mulig å slå sammen som interface{}?
/*
func Admin_init(IDInput int, buttonChan <-chan Button, floorSensorChan <-chan int,
	localOrderChan chan<- Order, adminTChan chan<- Udp, adminRChan <-chan Udp, backupChan <-chan Backup/*, peerChangeChan <-chan Peer,
	startTimerChan chan<- string, timeOutChan <-chan string) {

	go admin(IDInput, buttonChan, floorSensorChan,
		localOrderChan, adminTChan, adminRChan, backupChan/*, peerChangeChan,
		startTimerChan, timeOutChan)
}
*/

/*
type ActiveLift struct {
	LiftID int
}
type ActiveLifts []ActiveLift

func (slice ActiveLifts) Len() int {
	return len(slice)
}

func (slice ActiveLifts) Less(i, j int) bool {
	return slice[i].LiftID < slice[j].LiftID
}

func (slice ActiveLifts) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
*/

func Admin(IDInput int, buttonChan <-chan Button, floorSensorChan <-chan int,
	localOrderChan chan<- Order, adminTChan chan<- Udp, adminRChan <-chan Udp, backupTChan chan<- BackUp, backupRChan <-chan BackUp,
	peerChangeChan <-chan Peer, startTimerChan chan<- string, timeOutChan <-chan string) {

	orders := Create_order_matrix()
	properties := Create_lift_prop_list()
	ID := IDInput //SETT ID (spør nett?)
	// Check om orders/prop list eksisterer noen andre plasser på nettet
	aliveLifts := make([]int, 0, MAX_N_LIFTS)
	aliveLifts = append(aliveLifts, ID)
	//For test
	//aliveLifts = append(aliveLifts, 1)

	//lastButtonPressed := Button{NOT_VALID, NOT_VALID} //Går uten dette
	//bi := 0

	// Enter init state (også default, så trengs ikke)
	//Set_state(properties, ID, IDLE)

	//Spør nett om noen har orders og properties. If so, sett orders og properties lik de på nettet. If not, forsett med det samme.

	// Tror det er uavhengig av hvilken state det står i for det her: Vil bare vite om vi står i en etasje.

	//Kan bruke PEERS for å få inn ALLE
searchingForBackupLoop:
	for {
		select {
		//Legg inn en måte å få inn andre som er på nett?
		case backup := <-backupRChan:
			orders = backup.Orders
			properties = backup.Properties
			break searchingForBackupLoop

		case <-time.After(5 * time.Second):
			break searchingForBackupLoop
		}
	}

initLoop:
	for {
		select {
		//Legg inn en måte å få inn andre som er på nett?
		case f := <-floorSensorChan:
			fmt.Println("Adm: initLoop, floor Sensor")
			Set_last_floor(properties, ID, f)
			localOrderChan <- Order{"FLOOR_LIGHT", NOT_VALID, f, ON}
			Assign_orders(orders, f, ID)
			localOrderChan <- Order{"DIRN", DIRN_STOP, NOT_VALID, ON} // Bør sikkert endre navn for å gjør den her enkler å forstå
			startTimerChan <- "DOOR_OPEN"
			Set_state(properties, ID, DOOR_OPEN)
			Set_dirn(properties, ID, DIRN_DOWN)
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
		// Problem med å sende melding om button pressed ut på nettet og deretter melding fra findNewOrder?
		// evt legge ved hvilke ordre vi tar hver gang i findNewOrder-melding => alle andre kan oppdatere.
		// Husk "problem" med at assign bare tar de som allerede finnes, så
		// må ha en måte å slå sammen her.

		case b := <-buttonChan: // INSIDE AND OUTSIDE
			//If order already exists (legg til funksjon i order_matrix), sett
			// legg ny order inn midlertidig plass til får melding fra network om alene
			// Eller: Kan jo sende til NW, få tilbake, så legge til hvis
			// ikke alene. Legge til uansett hvis indre.

			/*if b.Button_dir == BUTTON_COMMAND {
				Add_order(orders, b.Floor, ID, b.Button_dir) // Inside order
			} else if len(aliveLifts) > 1 {
				Add_order(orders, b.Floor, ID, b.Button_dir) // Outside order. Tas bare når vi vet at andre heiser eksisterer.
			}*/

			//Send melding ut til NW på adminTChan
			adminTChan <- Udp{ID, "ButtonPressed", b.Floor, b.Button_dir}
			/*if b != lastButtonPressed {
				adminTChan <- Udp{ID, "ButtonPressed", b.Floor, b.Button_dir}
				bi++

				if bi >= 5 {
					lastButtonPressed = Button{NOT_VALID, NOT_VALID}
				}
			}*/

			//Tanke: Legg inn noe som gjør at det ikke legges til(sendes ut på NW) hvis allerede finnes i orders.
			/*if not in orders {
				adminTChan <- Udp{ID, "ButtonPressed", b.Floor, b.Button_dir}
			}*/

		case fs := <-floorSensorChan:
			switch Get_state(properties, ID) {
			case DOOR_OPEN:
				//Intentionally blank, probably might as well just remove this case, right now for completeness
				// Just needs to break, which it will do without these. Maybe a small sleep on these?
			case IDLE:
				// See DOOR_OPEN
			case MOVING:
				if fs != Get_last_floor(properties, ID) {
					Set_last_floor(properties, ID, fs)
					localOrderChan <- Order{"FLOOR_LIGHT", NOT_VALID, fs, ON}
					fmt.Println("Adm: Verdier på vei inn i Should_stop: (orders, properties, fs, ID)")
					fmt.Println("Adm: ", orders, properties, fs, ID)
					if ShouldStop(orders, properties, fs, ID) == true { // LEGG INN ALLTID STOPP ØVERST OG NEDERST
						fmt.Println("Adm: Should_stop")
						Assign_orders(orders, fs, ID) // In case du tar en som ikke var assigna.
						//localOrderChan <-  Send "DIRN", DIRN_STOP, NOT_VALID, ON
						localOrderChan <- Order{"DIRN", DIRN_STOP, NOT_VALID, ON}
						Set_state(properties, ID, DOOR_OPEN)
						startTimerChan <- "DOOR_OPEN"
						Complete_order(orders, fs, ID)
						/*
							localOrderChan <- Order{"LIGHT", BUTTON_COMMAND, fs, OFF}
							localOrderChan <- Order{"LIGHT", BUTTON_CALL_UP, fs, OFF}
							localOrderChan <- Order{"LIGHT", BUTTON_CALL_DOWN, fs, OFF}
						*/

						//Send melding ut til NW på adminTChan
						// ID, "Stoppet", etasje (DOOR_OPEN)
						adminTChan <- Udp{ID, "Stopped", fs, NOT_VALID}
					} else {
						fmt.Println("Adm: Should_stop NOT")
						//Send melding ut til NW på adminTChan
						// ID, "kjørte forbi", etasje
						adminTChan <- Udp{ID, "DrovePast", fs, NOT_VALID}
						fmt.Println("Adm: Under teit beskjed")
					}
				}
			}

		case <-timeOutChan:
			//localOrderChan <-  Send "DOOR", NOT_VALID, NOT_VALID, OFF
			fmt.Println("Adm: Fikk timeout")
			localOrderChan <- Order{"DOOR", NOT_VALID, NOT_VALID, OFF}

			//TURN OFF LIGHTS!
			localOrderChan <- Order{"LIGHT", BUTTON_COMMAND, Get_last_floor(properties, ID), OFF}
			localOrderChan <- Order{"LIGHT", BUTTON_CALL_UP, Get_last_floor(properties, ID), OFF}
			localOrderChan <- Order{"LIGHT", BUTTON_CALL_DOWN, Get_last_floor(properties, ID), OFF}

			findNewOrder(orders, ID, properties, aliveLifts, startTimerChan, localOrderChan, adminTChan)

		case m := <-adminRChan:
			switch m.ID {
			case ID:
				//Alt for egen heis
				switch m.Type {
				case "ButtonPressed":
					fmt.Println("Adm: Får tilbake fra network")
					Add_order(orders, m.Floor, m.ID, m.ExtraInfo) // Ta bort den over og la den her stå? bedre her her.
					localOrderChan <- Order{"LIGHT", m.ExtraInfo, m.Floor, ON}

					if Get_state(properties, ID) == IDLE {
						fmt.Println("Adm: State == IDLE når knapp trykket på")
						findNewOrder(orders, ID, properties, aliveLifts, startTimerChan, localOrderChan, adminTChan)
					}
				case "Stopped":
					// Får ingenting tilbake fra andre her.
				//case "DrovePast":
				// Ingenting
				//fmt.Println("Adm: DrovePast kommer rundt")
				case "NewOrder":
					// Gjør alt før, er bare ack her. Skal det i det hele tatt komme tilbake hit?
				case "Idle":
					// Samme som over. Nada.
				}

			default: //Any other lift
				switch m.Type {
				case "ButtonPressed":
					Add_order(orders, m.Floor, m.ID, m.ExtraInfo)
					if m.ExtraInfo == BUTTON_CALL_UP || m.ExtraInfo == BUTTON_CALL_DOWN {
						localOrderChan <- Order{"LIGHT", m.ExtraInfo, m.Floor, ON}
					}
				case "Stopped":
					Assign_orders(orders, m.Floor, m.ID)
					Complete_order(orders, m.Floor, m.ID)
					localOrderChan <- Order{"LIGHT", BUTTON_CALL_UP, m.Floor, OFF}
					localOrderChan <- Order{"LIGHT", BUTTON_CALL_DOWN, m.Floor, OFF}
					Set_state(properties, m.ID, DOOR_OPEN)
					Set_last_floor(properties, m.ID, m.Floor)

				case "DrovePast":
					Set_last_floor(properties, m.ID, m.Floor)
					Set_state(properties, m.ID, MOVING)
					fmt.Println("Adm: Men ikke hit")

				case "NewOrder":
					Set_state(properties, m.ID, MOVING)
					Set_dirn(properties, m.ID, GetNewDirection(m.Floor, Get_last_floor(properties, m.ID)))
				case "Idle":
					Set_state(properties, m.ID, IDLE)
				}
			}

		case backupMsg := <-backupRChan:
			switch backupMsg.Info {
			case "VarAlene":
				// Legg inn alle INDRE ordre for backupMsg.SenderID
				// Ta inn properties for backupMsg.SenderID
			case "IkkeAlene":
				// Skriv over alt i orders minus egne indre ordre.
				// Behold egne properties, skriv over resten.
			}

		case peerMsg := <-peerChangeChan:
			switch peerMsg.Change {
			case "New":
				if len(aliveLifts) == 1 {
					//send "var alene", ID, orders, properties
					backupTChan <- BackUp{"VarAlene", ID, orders, properties}
				} else {
					// send "ikke alene", ID, orders, properties
					backupTChan <- BackUp{"IkkeAlene", ID, orders, properties}
				}
				aliveLifts = append(aliveLifts, peerMsg.ChangedPeer)
				sort.Slice(aliveLifts, func(i, j int) bool { return aliveLifts[i] < aliveLifts[j] }) //MÅ FIKSES. NB NB NB

			case "Lost":
				for i, n := range aliveLifts {
					if n == peerMsg.ChangedPeer {
						lostPeer := n
						aliveLifts = append(aliveLifts[:i], aliveLifts[i+1:]...)
						DeassignOuterOrders(orders, lostPeer)
						if Get_state(properties, ID) == IDLE {
							fmt.Println("Adm: State == IDLE, en annen heis er død => kan være nye ordre")
							findNewOrder(orders, ID, properties, aliveLifts, startTimerChan, localOrderChan, adminTChan)
						}
					}
				}
			}

		}
	}
}

func findNewOrder(orders [][]int, ID int, properties []int, aliveLifts []int, startTimerChan chan<- string,
	localOrderChan chan<- Order, adminTChan chan<- Udp) {
	fmt.Println("Adm: Inne i findNewOrder. Orders, properties: ", orders, properties)

	newDirn, dest := CalculateNextOrder(orders, ID, properties, aliveLifts)
	// Should change name on both module called from and function itself.
	// Default dest and newDirn returned has to be undefined (-2,-2)
	fmt.Println("Adm: Got new direction", newDirn, dest)
	if newDirn == DIRN_STOP {
		fmt.Println("Adm: I DIRN_STOP for findNewOrder")
		Assign_orders(orders, dest, ID) //NB! Nå lagt til ALLE på den etasjen,
		// noe som er en forenkling som vi kunne gjøre. IKKE TESTET ENNÅ
		//localOrderChan <-  Send "DOOR", NOT_VALID, NOT_VALID, ON
		localOrderChan <- Order{"DOOR", NOT_VALID, NOT_VALID, ON}
		Set_state(properties, ID, DOOR_OPEN)
		startTimerChan <- "DOOR_OPEN"
		Complete_order(orders, dest, ID)
		/*
			//TURN OFF LIGHTS!
			localOrderChan <- Order{"LIGHT", BUTTON_COMMAND, dest, OFF}
			localOrderChan <- Order{"LIGHT", BUTTON_CALL_UP, dest, OFF}
			localOrderChan <- Order{"LIGHT", BUTTON_CALL_DOWN, dest, OFF}
		*/

		//Send melding ut til NW på adminTChan
		// ID, "Stoppet", etasje (DOOR_OPEN)
		adminTChan <- Udp{ID, "Stopped", dest, NOT_VALID}
	} else if newDirn == DIRN_DOWN || newDirn == DIRN_UP {
		fmt.Println("Adm: I DIRN_DOWN/DIRN_UP for findNewOrder")
		Assign_orders(orders, dest, ID)
		//localOrderChan <-  Send "DIRN", DIRN_UP/DOWN, NOT_VALID, NOT_VALID
		localOrderChan <- Order{"DIRN", newDirn, NOT_VALID, NOT_VALID}
		Set_state(properties, ID, MOVING)
		Set_dirn(properties, ID, newDirn)
		//Send melding ut til NW på adminTChan
		// ID, "Moving, desting (new order)", etasje
		adminTChan <- Udp{ID, "NewOrder", dest, NOT_VALID}
	} else { // newDirn == -2 (NOT_VALID)
		fmt.Println("Adm: I IDLE for findNewOrder")
		Set_state(properties, ID, IDLE)
		//Send melding ut til NW på adminTChan
		// ID, "IDLE", etasje
		adminTChan <- Udp{ID, "Idle", dest, NOT_VALID}
	}
	fmt.Println("Adm: På vei ut av findNewOrder. Orders, properties: ", orders, properties)
}
