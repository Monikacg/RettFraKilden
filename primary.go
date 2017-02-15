package main

import (
  "fmt"
  "strconv"
  "time"
  "os"
  "os/exec"
  "net"
)


func main()  {
  out := make(chan string)
  go transmitUDP(out)
  start_backup()

  counter := 0

  fmt.Println("I am primary")

  if len(os.Args) > 1 {
    counter, _ = strconv.Atoi(os.Args[1])
    fmt.Printf("Master started with counter = %d\n", counter)
  }

  for {
    fmt.Println(counter)
    out <- strconv.Itoa(counter)
    counter++
    time.Sleep(500*time.Millisecond)
  }

}

func start_backup()  {
  cmd := "\"Tell Application 'Terminal' to do script go run backup.go\""
  //out, err := exec.Command("osascript","-e",cmd).Output() //do script go run backup.go
  //if err != nil {
    //fmt.Printf("%s", err)
  //}
  //fmt.Printf("%s", out)
  err2 := exec.Command("osascript","-e",cmd).Start()
  
  if err2 != nil {
    fmt.Printf("%s", err2)
  }
}

func transmitUDP(out <-chan string)  {
  laddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", 22002))//10.22.70.209
  baddr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", 44004))
  conn, _ := net.ListenUDP("udp4", laddr)

  for {
    select {
    case message := <- out:
      conn.WriteToUDP([]byte(message), baddr)

    }
  }
}
