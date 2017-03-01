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
	button_chan := make(chan Button, 100)
	floor_sensor_chan := make(chan int, 100)

	local_order_chan := make(chan Order, 100)

	adm_transmit_chan := make(chan int, 100)
	adm_receive_chan := make(chan int, 100)
	peer_chan := make(chan int, 100)

	start_timer_chan := make(chan string, 100)
	time_out_chan := make(chan string, 100)

	//Ta inn IP som input i shell, map til ID. Send ID til NETWORK og ADMIN

	Elev_init()

	go Lift_control_init(button_inside_chan, button_outside_chan,
		floor_sensor_chan, local_order_chan)

	go Network_init(IDInput, adm_transmit_chan, adm_receive_chan, peer_chan)

	go Admin_init(IDInput, button_inside_chan, button_outside_chan, floor_sensor_chan,
		local_order_chan, adm_transmit_chan, adm_receive_chan, peer_chan,
		start_timer_chan, time_out_chan)

	go Timer_init(start_timer_chan, time_out_chan)
}
