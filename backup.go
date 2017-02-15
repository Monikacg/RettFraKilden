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
    }
    case <- time.After(5*time.Second):
      fmt.Println("Assumes primary is dead.")
      start_primary()
  }

}

func listen(in <-chan string)  {
  laddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":", 44004))
  conn, _ := net.ListenUDP("udp4", laddr)
  for {
    buf := make([]byte, 1024)
    n, _, _ := conn.ReadFromUDP(buf)
    in <- string(buf[:n])
  }
}

func start_primary()  {
  arg := fmt.Sprintf("go run primary.go ‰d", counter)
  full_arg := "\"Tell Application 'Terminal' to do script "
  full_arg += arg + "\""
  cmd := exec.Command("osascript","-e",full_arg)
  _ := cmd.Run()
}
