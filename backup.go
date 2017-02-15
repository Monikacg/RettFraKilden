package main

import (
  "fmt"
  "strconv"
  "time"
  "os/exec"
  "net"
)


func main()  {
  in := make(chan string)
  go listen(in)

  counter := 0

  fmt.Println("I am backup")

  // Ta bort etter test
  select {
  case message := <- in:
    counter, _ = strconv.Atoi(message)
    fmt.Printf("Backup gets start = %d\n", counter)
  case <- time.After(5*time.Second):
    fmt.Println("No start value in backup")
  }

  for {
    select {
    case message := <- in:
      counter, _ = strconv.Atoi(message)
  case <- time.After(5*time.Second):
    fmt.Println("Assumes primary is dead.")
    start_primary(counter)
    }
  }

}

func listen(in chan<- string)  {
  laddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", 44004))
  conn, _ := net.ListenUDP("udp4", laddr)
  for {
    buf := make([]byte, 1024)
    n, _, _ := conn.ReadFromUDP(buf)
    in <- string(buf[:n])
  }
}

func start_primary(counter int)  {
  // arg := fmt.Sprintf("go", "run", "primary.go", "‰d", counter) Fungerer ikke
  //full_arg := "\"Tell Application 'Terminal' to do script "
  //full_arg += arg + "\""
  s_counter := strconv.Itoa(counter)
  cmd := exec.Command("gnome-terminal","-x","go", "run", "primary.go", s_counter)//"sh","-c",
  err := cmd.Run()
  if err != nil {
    fmt.Printf("%s", err)
  }
}
