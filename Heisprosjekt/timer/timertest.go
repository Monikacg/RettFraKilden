package main

// MAIN FUNCTION USED TO TEST TIMER

import (
  "fmt"
  "time"
  //"github.com/Arktoz/Sanntid/Heisprosjekt/timer"
)

func main()  {
  start_timer_chan := make(chan string, 100)
	time_out_chan := make(chan string, 100)
  go timer_init(start_timer_chan, time_out_chan)

  go listen(time_out_chan)

  for {
    start_timer_chan <- "DOOR_OPEN"
    fmt.Println("Starting timer for DOOR_OPEN")

    start_timer_chan <- "UDP"
    fmt.Println("Starting timer for UDP")

    time.Sleep(10*time.Second)
  }
}

func listen(time_out_chan <-chan string)  {
  for {
    select {
    case m := <- time_out_chan:
      if m == "DOOR_OPEN" {
        fmt.Println("DOOR_OPEN")
      }
      if m == "UDP" {
        fmt.Println("UDP")
      }
    }
  }
}
