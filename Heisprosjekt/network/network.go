package network

import (
	"fmt"
	"os"
	"time"

	. "./../definitions"
	"./bcast"
	"./localIP"
	"./peers"
)

/*
import(
	"./bcast"
	"./localip"
	"./peers"
	"flag"
	"fmt"
	"os"
	"time"
)
*/

// All functions used by this file (found in the peers, localip, conn and bcast folders)
// were made by github.com/klasbo.

func Network_init(IDInput int, adminTChan <-chan Udp, adminRChan chan<- Udp, backupChan chan<- Backup, peerChangeChan chan<- Peer) {
	//vet ikke helt hva som skal være her enda
	fmt.Println(DIRN_STOP)
	network(IDInput, adminTChan, adminRChan, backupChan, peerChangeChan)
}

func network(IDInput int, adminTChan <-chan Udp, adminRChan chan<- Udp, backupChan chan<- Backup, peerChangeChan chan<- Peer) {

	/*
		//GITT PÅ FORHÅND
		var id string
		flag.StringVar(&id, "id", "", "id of this peer")
		flag.Parse()


			if id == "" {
				localIP, err := localip.LocalIP()
				if err != nil {
					fmt.Println(err)
					localIP = "DISCONNECTED"
				}
				id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
			}
	*/

	/*
	   Ting som må sendes:
	   ID, Type (type button pressed/arrived+stopped/arrived+drove past/
	 moving,dest (new order)/IDLE), floor, typeOfButton
	*/

	ownID := IDInput

	currentPeers := make([]int, 0, MAX_N_LIFTS)
	currentPeers = append(currentPeers, ownID)

	const timeout = 50 * time.Millisecond

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("Peer-%d-%s-%d", ownID, localIP, os.Getpid())

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15640, id, peerTxEnable) //15647
	go peers.Receiver(15640, peerUpdateCh)        //15647

	/*

		helloTx := make(chan HelloMsg)
		helloRx := make(chan HelloMsg)
		go bcast.transmiter(16569, helloTx) //16569
		go bcast.Receiver(16569, helloRx)    //16569

		//Test som er gitt
		go func() {
			helloMsg := HelloMsg{"Hello from " + id, 0}
			for {
				helloMsg.Iter++
				helloTx <- helloMsg
				time.Sleep(1 * time.Second)
			}
		}()

	*/

	/*
		Lag en struct som er Udp-strukten+ack-felt som sier hvem som vet om beskjeden, ok ide?
	*/

	helloTx := make(chan OverNetwork)
	helloRx := make(chan OverNetwork)
	go bcast.Transmitter(16570, helloTx)
	go bcast.Receiver(16570, helloRx)

	backupTx := make(chan Backup)
	backupRx := make(chan Backup)
	go bcast.Transmitter(16571, backupTx)
	go bcast.Receiver(16571, backupRx)

	localAckCh := make(chan OverNetwork)

	outPutCh := make(chan OverNetwork)

	fmt.Println("Started")

	/*
		New Peer er string, Lost er []string: Hvorfor? begge IP? trenger IP for map til ID
	*/

	for {
		select {

		case u := <-adminTChan:
			// Når en får beskjed fra admin om noe: legg inn i sende-struct,
			// create den siste lista, legg til sin egen ID som en som vet om,
			// send via bcast Transmitter. (helloTx-kanalen so far)

			// NB!!! Send rett tilbake hvis ingen andre på nett.

			if len(currentPeers) == 1 {
				adminRChan <- u
			} else {
				O := OverNetwork{Udp: u}
				O.WhoKnowsAboutThis = make([]int, 0, MAX_N_LIFTS+1)
				O.WhoKnowsAboutThis = append(O.WhoKnowsAboutThis, O.Udp.ID)

				go func(O OverNetwork, localAckCh chan OverNetwork, helloTx chan<- OverNetwork,
					helloRx <-chan OverNetwork, adminRChan chan<- Udp) {

					resendTimer := time.NewTimer(timeout)
					localO := O
					helloTx <- localO

				ackLoop:
					for {
						select {
						//Kan legge inn timer som sender ut etter en stund hvis ingenting annet skjer.
						case v := <-helloRx:
							//Hvis ferdig, send til adm
							if v.Udp == localO.Udp {
								localO = v
								resendTimer = time.NewTimer(timeout)
								if v.WhoKnowsAboutThis[len(v.WhoKnowsAboutThis)-1] == -1 {
									//send it out one last time.
									helloTx <- v
									// Done with order, send to Admin.
									adminRChan <- localO.Udp
									break ackLoop
								} else {
									//Vår egen. Send ut igjen, samme som før
									helloTx <- v
								}
							} else {
								if v.Udp.ID == ownID {
									// Skal til en annen go-routine (ack)
									localAckCh <- v
								} else {
									// Melding kommer fra annen heis.
									outPutCh <- v
								}
							}
						case ack := <-localAckCh:
							// Samme som over.
							if ack.Udp == localO.Udp {
								localO = ack
								resendTimer = time.NewTimer(timeout)
								if ack.WhoKnowsAboutThis[len(ack.WhoKnowsAboutThis)-1] == -1 {
									//send it out one last time.
									helloTx <- ack
									// Done with order, send to Admin.
									adminRChan <- localO.Udp
									break ackLoop
								} else {
									//Vår egen. Send ut igjen, samme som før
									helloTx <- ack
								}
							} else {
								// Skal til en annen go-routine (ack) (vet at ID == ownID)
								localAckCh <- ack
							}

						case <-resendTimer.C:
							helloTx <- localO
							resendTimer = time.NewTimer(timeout)
						}
					}
				}(O, localAckCh, helloTx, helloRx, adminRChan)
			}

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			if len(p.New) > 0 {
				newID := 2
				//funksjon som mapper fra IP til ID her NB NB NB mangler
				//legg til i currentPeers
				currentPeers = append(currentPeers, newID)
				sort.Slice(currentPeers, func(i, j int) bool { return currentPeers[i] < currentPeers[j] })
				peerChangeChan <- Peer{"New", newID}
			}
			if len(p.Lost) > 0 {
				lostID := 3
				//funksjon som mapper fra IP til ID her NB NB NB mangler
				//ta bort fra currentPeers
				for i, n := range currentPeers {
					if n == lostID {
						currentPeers = append(currentPeers[:i], currentPeers[i+1:]...)
					}
				}
				peerChangeChan <- Peer{"Lost", lostID}
			}

		case recv := <-helloRx:
			if recv.Udp.ID == ownID {
				localAckCh <- recv
			} else {
				if recv.WhoKnowsAboutThis[len(recv.WhoKnowsAboutThis)-1] == -1 {
					// Done with order, send to Admin.
					adminRChan <- recv.Udp
				} else {
					thisLiftKnows := false
					for _, n := range recv.WhoKnowsAboutThis {
						if n == ownID {
							thisLiftKnows = true
						}
					}
					if thisLiftKnows == false {
						recv.WhoKnowsAboutThis = append(recv.WhoKnowsAboutThis, ownID)
					}

					if len(recv.WhoKnowsAboutThis) == len(currentPeers) {
						//Shows it's done
						recv.WhoKnowsAboutThis = append(recv.WhoKnowsAboutThis, -1)
						//send it out
						helloTx <- recv
						// Done with order, send to Admin.
						adminRChan <- recv.Udp
					} else {
						//send it out
						helloTx <- recv
					}

				}
			}
			//if ID lik ownID, send på localAckCh. If not:
			// append vår ID hvis den ikke er i lista. Hvis vi er siste,
			// legg til -1 på slutten, send ut og send til admin.

		case out := <-outPutCh:
			// append vår ID hvis den ikke er i lista. Hvis vi er siste,
			// legg til -1 på slutten, send ut og send til admin.
			//In effect samme som over minus første test.
			if out.WhoKnowsAboutThis[len(out.WhoKnowsAboutThis)-1] == -1 {
				// Done with order, send to Admin.
				adminRChan <- out.Udp
			} else {
				thisLiftKnows := false
				for _, n := range out.WhoKnowsAboutThis {
					if n == ownID {
						thisLiftKnows = true
					}
				}
				if thisLiftKnows == false {
					out.WhoKnowsAboutThis = append(out.WhoKnowsAboutThis, ownID)
				}

				if len(out.WhoKnowsAboutThis) == len(currentPeers) {
					//Shows it's done
					out.WhoKnowsAboutThis = append(out.WhoKnowsAboutThis, -1)
					//send it out
					helloTx <- out
					// Done with order, send to Admin.
					adminRChan <- out.Udp
				} else {
					//send it out
					helloTx <- out
				}
			}

			//fmt.Printf("Received: %#v\n", a)
			// Når WhoKnowsAboutThis bare er -1, skal resten av structen sendes til admin
			// Ellers skal vi legge til vår egen ID i WhoKnowsAboutThis (hvis den ikke er der),
			// så sende meldingen (ellers) uendret ut igjen.

			// Husk å send på nytt
		}
		//time.Sleep(1 * time.Second)
	}
}
