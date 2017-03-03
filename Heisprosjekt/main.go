package main

import (
	. "./admin"
	. "./driver"
	. "./lift_control"
	. "./network"
	. "./timer"
)

func main() {
	// Alle skal kanskje ikke ha int, kanskje endre til struct på noen (U/I 1/2/3)
	buttonChan := make(chan Button, 100)
	floorSensorChan := make(chan int, 100)

	localOrderChan := make(chan Order, 100)

	adminTChan := make(chan Udp, 100)
	adminRChan := make(chan Udp, 100)
	backupChan := make(chan Backup, 100)
	peerChangeChan := make(chan int, 100)

	startTimerChan := make(chan string, 100)
	timeOutChan := make(chan string, 100)

	//Ta inn IP som input i shell, map til ID. Send ID til NETWORK og ADMIN

	Elev_init()

	go Lift_control_init(buttonChan,
		floorSensorChan, localOrderChan)

	go Network_init(IDInput, adminTChan, adminRChan, backupChan, peerChangeChan)

	go Admin_init(IDInput, buttonChan, floorSensorChan,
		localOrderChan, adminTChan, adminRChan, backupChan, peerChangeChan,
		startTimerChan, timeOutChan)

	go Timer_init(startTimerChan, timeOutChan)
}
