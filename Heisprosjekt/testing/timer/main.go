package main

import (
  "fmt"
  "time"
)

func main()  {
  start_timer_chan := make(chan string, 100)
	time_out_chan := make(chan string, 100)
  go timer.timer_init(start_timer_chan, time_out_chan)
}
