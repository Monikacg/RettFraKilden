package main

import (
	"sync"

	. "./admin"
	. "./definitions"
	. "./driver"
	. "./lift_control"
	. "./network"
	. "./timer"
)

func main() {
	// Alle skal kanskje ikke ha int, kanskje endre til struct på noen (U/I 1/2/3)
	buttonChan := make(chan Button)
	floorSensorChan := make(chan int)

	localOrderChan := make(chan Order)

	adminTChan := make(chan Udp, 100) // Må være asynkron
	adminRChan := make(chan Udp, 100) // ----"-------"---
	backupTChan := make(chan BackUp, 100)
	backupRChan := make(chan BackUp, 100)
	peerChangeChan := make(chan Peer, 100)

	startTimerChan := make(chan string)
	timeOutChan := make(chan string)

	//Ta inn IP som input i shell, map til ID. Send ID til NETWORK og ADMIN
	// Counterpoint: Kan ta ID som direkte input?

	IDInput := 0

	Elev_init()
	var wg sync.WaitGroup
	wg.Add(1)

	go Lift_control(buttonChan, floorSensorChan, localOrderChan)

	go Network(IDInput, adminTChan, adminRChan, backupTChan, backupRChan, peerChangeChan)

	go Admin(IDInput, buttonChan, floorSensorChan,
		localOrderChan, adminTChan, adminRChan, backupTChan, backupRChan, peerChangeChan, startTimerChan, timeOutChan)

	go Timer(startTimerChan, timeOutChan)
	wg.Wait()
}
