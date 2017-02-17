package network

import(
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)
/*
serializing
oppdater listen over heiser som er ilive :)
*/

type Udp struct {
	id int
	message string
	properties_struct
}

func Network_init(adm_transmitt_chan <-chan udp, adm_receive_chan chan<- udp, peer_chan chan<- int){
	//vet ikke helt hva som skal være her enda
}

func main() {

	//GITT PÅ FORHÅND
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	/*
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}
	*/

	peerUpdateCh := make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	go bcast.Transmitter(16569, helloTx) //16569
	go bcast.Receiver(16569, helloRx) //16569

	//Test som er gitt
	go func() {
		helloMsg := HelloMsg{"Hello from " + id, 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
		time.Sleep(1 * time.Second)
	}
	/////////////////////////////////////////////////////////////////////////////////////////////

}
