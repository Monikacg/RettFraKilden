package network

import (
	"fmt"
	"os"

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

/*
serializing
oppdater listen over heiser som er ilive :)
*/

//Heismodul (alle deler minus denne .go-filen) laget av github.com/klasbo.

type HelloMsg struct {
	Message string
	Iter    int
}

func Network_init(IDInput int, adm_transmit_chan <-chan Udp, adm_receive_chan chan<- Udp, peer_chan chan<- Peer) {
	//vet ikke helt hva som skal være her enda
	fmt.Println(DIRN_STOP)
	network(IDInput, adm_transmit_chan, adm_receive_chan, peer_chan)
}

func network(IDInput int, adm_transmit_chan <-chan Udp, adm_receive_chan chan<- Udp, peer_chan chan<- Peer) {

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

	ID := IDInput

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("Peer-%d-%s-%d", ID, localIP, os.Getpid())

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

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
	go bcast.Transmitter(16569, helloTx) //16569
	go bcast.Receiver(16569, helloRx)    //16569

	fmt.Println("Started")

	/*
		Hvordan bruke denne: Kan det som var helloRx/Tx være hvilken som helst struct?
		New Peer er string, Lost er []string: Hvorfor? begge IP? trenger IP for map til ID

		Hvordan endre det som kommer ut fra helloRx tilbake til noe fornuftig?
	*/

	for {
		select {

		case u := <-adm_transmit_chan:

		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

			if len(p.New) > 0 {
				newID := 2
				//funksjon som mapper fra IP til ID her
				peer_chan <- Peer{"New", newID}
			}
			if len(p.Lost) > 0 {
				lostID := 3
				//funksjon som mapper fra IP til ID her
				peer_chan <- Peer{"Lost", lostID}
			}

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
		//time.Sleep(1 * time.Second)
	}
}
