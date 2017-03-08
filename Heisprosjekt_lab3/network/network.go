package network

import (
	"fmt"
	"strconv"
	"time"

	. "../definitions"
	"./bcast"
	"./peers"
	"sort"
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
/*
func Network_init(IDInput int, adminTChan <-chan Udp, adminRChan chan<- Udp, backupChan chan<- Backup/*, peerChangeChan chan<- Peer) {
	//vet ikke helt hva som skal være her enda
	//fmt.Println(DIRN_STOP)
	network(IDInput, adminTChan, adminRChan, backupChan/*, peerChangeChan)
}
*/

type Ack struct {
	Message Udp
	Counter int
}


type CurrPeers struct {
	CurrPeers []int
}

/*
func sendAcks(helloTx chan<- Udp, currentPeersChan <-chan CurrPeers, sendToAckChan <-chan Udp, receivedAckChan <-chan Udp, adminRChan chan<- Udp) {
	var sendQueue []int
	var peers []int
	const timeout = 100 * time.Millisecond

	var a []Ack

	ackTimer := time.NewTimer(timeout)
	for {
		select {
		case cP := <- currentPeersChan:
			peers = cP.CurrPeers
		case sentMsg := <- sendToAckChan:
			a = append(a, Ack{Message: sentMsg, Counter: 0})
			for i := 0; i < 5; i++ {
				helloTx <- sentMsg
				time.Sleep(1*time.Millisecond)
			}
			ackTimer = time.NewTimer(timeout)
		case recvMsg := <- receivedAckChan:

		case <-ackTimer.C:
				ackTimer = time.NewTimer(timeout)

				//Slutt å send melding rundt. Alle som har nok acker, sender du til admin. Så slettes fra a
				//adminRChan <-
				// Ta bort fra a.
		}
	}


}
*/

func Network(IDInput int, adminTChan <-chan Udp, adminRChan chan<- Udp, backupTChan <-chan BackUp, backupRChan chan<- BackUp, peerChangeChan chan<- Peer) {


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
	// for test:
	//currentPeers = append(currentPeers, 1)

	const timeout = 200 * time.Millisecond

	currentPeersChan := make(chan CurrPeers, 100) // ENDRE! LAG STRUCT FOR currentPeers.
	sendToAckChan := make(chan Udp, 100)
	//receivedAckChan := make(chan Udp, 100)
	//outPutCh := make(chan Udp, 100)



	/*
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id := fmt.Sprintf("Peer-%d-%s-%d", ownID, localIP, os.Getpid())
	*/
	id := strconv.Itoa(ownID)

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

	helloTx := make(chan Udp)
	helloRx := make(chan Udp)
	go bcast.Transmitter(16570, id, helloTx)
	go bcast.Receiver(16570, id, false, helloRx)

	backupTx := make(chan BackUpNW)
	backupRx := make(chan BackUpNW)
	go bcast.Transmitter(16571, id, backupTx)
	go bcast.Receiver(16571, id, false, backupRx)

	//localAckCh := make(chan Udp, 100)


	//outPutCh := make(chan Udp, 100)

	lastMessage := Udp{NOT_VALID, "Test", NOT_VALID, NOT_VALID}

	fmt.Println("Started")

	/*go sendAcks(helloTx, currentPeersChan, sendToAckChan, receivedAckChan, adminRChan)
	c := CurrPeers{}
	c.CurrPeers = currentPeers
	currentPeersChan <- c
	/*
		New Peer er string, Lost er []string: Hvorfor? begge IP? trenger IP for map til ID
	*/

	for {
		select {
		case backupToSend := <-backupTChan:
			fmt.Println("backupToSend", backupToSend.SenderID, backupToSend.Info)

		//Legg til case for tatt imot backup en plass

		case u := <-adminTChan:
			// Når en får beskjed fra admin om noe: legg inn i sende-struct,
			// create den siste lista, legg til sin egen ID som en som vet om,
			// send via bcast Transmitter. (helloTx-kanalen so far)

			// NB!!! Send rett tilbake hvis ingen andre på nett.


			if lastMessage != u {
				fmt.Println("NW: len(currentPeers)", len(currentPeers))
				if len(currentPeers) == 1 {
					//TA BORT YTRE KNAPPER
					//FUNGERER IKKE:  TAR IKKE BORT YTRE KNAPPER; TAR IKKE BORT NOEN.
					if ((u.Type == "ButtonPressed" && u.ExtraInfo == DIRN_UP) || (u.Type == "ButtonPressed" && u.ExtraInfo == DIRN_DOWN)) {
						fmt.Println("Ikke i single elevator mode")
					} else {
						adminRChan <- u
					}
					fmt.Println("NW: Melding", u)
				} else {

					sendToAckChan <- u


					lastMessage = u
				}
			}

		case recv := <-helloRx:
			//If vår ID, send to ack (receivedAckChan)
			fmt.Println("hei", recv)
			// Else, check om den er lik den siste meldingen vi fikk fra den heisen. If so,
			// gjør ingenting. If not, send ut + send til admin.


		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			if len(p.New) > 0 {
				newID, _ := strconv.Atoi(p.New)
				if newID != ownID {
					newID := 2 //TA BORT
					//funksjon som mapper fra IP til ID her NB NB NB mangler
					//legg til i currentPeers
					currentPeers = append(currentPeers, newID)
					sort.Slice(currentPeers, func(i, j int) bool { return currentPeers[i] < currentPeers[j] })
					peerChangeChan <- Peer{"New", newID}

					c := CurrPeers{}
					c.CurrPeers = currentPeers
					currentPeersChan <- c
				}

			}
			if len(p.Lost) > 0 {
				lostID := 3
				//funksjon som mapper fra IP til ID her NB NB NB mangler
				//ta bort fra currentPeers
				// Iterer over LOST, gjør for hver.
				for i, n := range currentPeers {
					if n == lostID {
						currentPeers = append(currentPeers[:i], currentPeers[i+1:]...)
					}
				}
				c := CurrPeers{}
				c.CurrPeers = currentPeers
				currentPeersChan <- c
				peerChangeChan <- Peer{"Lost", lostID}
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
