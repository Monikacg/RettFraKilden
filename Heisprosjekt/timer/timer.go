package timer



import (
  "time"
  //"fmt" Bare for test
)

func timer_init(start_timer_chan <-chan string, time_out_chan chan<- string)  {
  go timer(start_timer_chan, time_out_chan) // Trengs egentlig "go" her? allerede kalt som goroutine fra main
}

func timer(start_timer_chan <-chan string, time_out_chan chan<- string)  {

  for {
    select {
    case start_msg := <- start_timer_chan:
      if start_msg == "DOOR_OPEN" {
        go door_open_timer(time_out_chan)
      }
      if start_msg == "UDP" {
        go udp_timer(time_out_chan)
      }
      // kan ha if/else if/else, else gir error-melding.
    }
  }
}

func door_open_timer(time_out_chan chan<- string)  {
  for {
    select {
    case <- time.After(3*time.Second):
      time_out_chan <- "DOOR_OPEN"
      return
    }
  }
}

func udp_timer(time_out_chan chan<- string)  {
  for {
    select {
    case <- time.After(100*time.Millisecond):
      time_out_chan <- "UDP"
      return
    }
  }
}

// Note: Dette vil enten blokkere sin egen goroutine (som vi vil)
// eller blokkere hele (som vi ikke vil). Test og finn ut
// => implementer overalt hvis det fungerer.