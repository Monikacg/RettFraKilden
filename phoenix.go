package main

import (
  "fmt"
  "strconv"
  "time"
  //"os"
  "os/exec"
  "net"
  //"strings"
)

func main()  {
    primary_alive := true
    var primary_counter int

    addr, _ := net.ResolveUDPAddr("udp4", ":22002")
    conn, _ := net.ListenUDP("udp4",addr)

    fmt.Println("I am backup")

    for primary_alive {
        conn.SetReadDeadline(time.Now().Add(3*time.Second))

        buf := make([]byte, 1024)
        n,_,err := conn.ReadFromUDP(buf)

        if err != nil {
            primary_alive = false
            fmt.Println("read err: ", n, err)
        } else {
            buf_s := string(buf[:n])
            fmt.Println(buf_s)
            primary_counter, err = strconv.Atoi(buf_s)
            fmt.Println(err, primary_counter)
        }
    }
    conn.Close()


    fmt.Println("I am master")

    start_new()
    addr, err := net.ResolveUDPAddr("udp4", "localhost:22002")
    fmt.Println("resolve:", err)
    dial, err := net.DialUDP("udp4",nil,addr)
    fmt.Println("dial:", err)

    for {
        message := strconv.Itoa(primary_counter)
        _, err := dial.Write([]byte(message))
        fmt.Println(primary_counter)
        if err != nil {
            fmt.Println("Hei", err)
        }

        primary_counter++
        time.Sleep(700*time.Millisecond)
    }
    dial.Close()
}
/*
func convert( b []byte) string {
    s := make([]string, len(b))
    for i := range b {
        s[i] = strconv.Itoa(int(b[i]))
    }
    return strings.Join(s,"")
}*/

func start_new()  {
    cmd := exec.Command("gnome-terminal","-x","go", "run", "phoenix.go") //"sh","-c",
    err := cmd.Run()
    if err != nil {
        fmt.Printf("%s", err)
    }
}
