package main

import (
    "fmt"
    "net"
    "time"
    "bufio"
    "os"
    //"strconv"
)

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}

func main() {
    ServerAddr,err := net.ResolveUDPAddr("udp","10.22.77.137:30000")
    CheckError(err)

    LocalAddr, err := net.ResolveUDPAddr("udp", "10.22.71.255")
    CheckError(err)

    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err)

    defer Conn.Close()
    //i := 0
    for {
        //msg := strconv.Itoa(i)
        //i++
        // read in input from stdin
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Text to send: ")
        msg, _ := reader.ReadString('\n')
        buf := []byte(msg)
        _,err := Conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Second * 1)
    }
}
