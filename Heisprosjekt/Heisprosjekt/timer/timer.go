package timer

//Har sikkert forbedringspotensiale. trur ikke det her klarer flere timer av gangen...
import (
  "time"
)
/*
var (
  door_flag bool = false
  udp_flag bool = false
  door_timer time.Time
  udp_timer time.Time
)

func timer_init(start_timer_chan, time_out_chan chan string)  {
  int_time_out_chan := make(chan int)
  timer(start_timer_chan, time_out_chan, int_time_out_chan)
}

func loopthatjuststartstimerswhenitgetstold(int_time_out_chan chan int)  {
  for {
    select{
      case <- door_timer.After(3*time.Seconds)
      if door_flag {
      }
      time.Sleep(100*time.Millisecond)
    }
  }
}

func timer(start_timer_chan, time_out_chan chan string, int_time_out_chan chan int)  {

  go loopthatjuststartstimerswhenitgetstold(int_time_out_chan)

  for {
    select {
    case t := <- start_timer_chan:
      if t == "DOOR_OPEN" {
        door_timer = time.Now()
        door_flag = true
      }
      if t == "UDP" {
        udp_timer = time.Now()
        udp_flag = true
      }
    case time_out := <- int_time_out_chan:
      if time_out == "DOOR_OPEN" {
        time_out_chan <- "DOOR_OPEN"
      } else if time_out == "UDP" {
        time_out_chan <- "UDP"
      } //Burde si hvilken den fikk timeout på. Men koden mangler mer for å skille mellom forskjellige timere også så langt.
    }
  }
}
*/

/*
Ny tankemåte: en hører og en gjører på timer: begge startes av init, bruker channels til å snakke mellom (fra listen til doing)
Gjører sender til adm, listen får fra adm. Går i liten sirkel.
Selve timer i doing, starter når... med mindre det bare gjør det vanskeliger? e det ikke bare det som gjøres
når det e i 1 funksjon? ... Må tenk når ingen snakke. Kan forsåvidt vær to tråa antar æ, men hjelpe det...
*/
func timer_init(start_timer_chan, time_out_chan chan string)  {
  go timer(start_timer_chan, time_out_chan) // Trengs egentlig "go" her? allerede kalt som goroutine fra main
}

func timer(start_timer_chan, time_out_chan chan string)  {

  for {
    select {
    case start_msg := <- start_timer_chan:
      if start_msg == "DOOR_OPEN" {
        go door_open_timer(time_out_chan)
      } else if start_msg == "UDP" {
        go udp_timer(time_out_chan)
      }
      // kan ha if/else if/else, else gir error-melding.
    }
  }
}

func door_open_timer(time_out_chan chan string)  {
  for {
    select {
    case <- time.After(3*time.Second):
      time_out_chan <- "DOOR_OPEN"
    }
  }
}

func udp_timer(time_out_chan chan string)  {
  for {
    select {
    case <- time.After(100*time.Millisecond):
      time_out_chan <- "UDP"
    }
  }
}

// Note: Dette vil enten blokkere sin egen goroutine (som vi vil)
// eller blokkere hele (som vi ikke vil). Test og finn ut
// => implementer overalt hvis det fungerer.
