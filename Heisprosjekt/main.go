package main

import (
	"fmt"
	"./driver"
	"./lift_control"
	"./network"
	"./admin"
)

func main() {
	// Alle skal kanskje ikke ha int, kanskje endre til struct på noen (U/I 1/2/3)
	button_inside_chan := make(chan int, 100)
	button_outside_chan := make(chan button, 100) //En kan være n+1?
	floor_sensor_chan := make(chan int, 100)

	local_order_chan := make(chan order, 100)

	adm_transmitt_chan := make(chan int, 100)
	adm_receive_chan := make(chan int, 100)
	peer_chan := make(chan int, 100)

	start_timer_chan := make(chan string, 100)
	time_out_chan := make(chan string, 100)
	interrupt_timer_chan := make(chan string, 100)

	driver.Elev_init()

	go lift_control.Lift_control_init(button_inside_chan, button_outside_chan,
		floor_sensor_chan, local_order_chan)

	go network.Network_init(adm_transmitt_chan, adm_receive_chan, peer_chan)

	go admin.Adm_init(button_inside_chan, button_outside_chan, floor_sensor_chan,
		local_order_chan, adm_transmitt_chan, adm_receive_chan, peer_chan,
		start_timer_chan, time_out_chan, interrupt_timer_chan)

	go timer.Timer_init(start_timer_chan, time_out_chan, interrupt_timer_chan)
}
