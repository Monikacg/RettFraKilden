package main

import (
	."consoleIO"
	."fmt"
	"os/exec"
	."strings"
	."strconv"
	."time"
	"net"
	
)
var localIP string = "10.22.79.27"

func main(){
	consoleIn := make(chan string)
	consoleOut := make(chan string)
	go InitConsoleIO(consoleIn,consoleOut)

	bossAlive := true 
	bossCount := 0
	startTime := 0
	udpaddr, _ := net.ResolveUDPAddr("udp", ":20011")
	conn, _ := net.ListenUDP("udp", udpaddr)
	
	
	Println("I am backup")
	
	for (bossAlive){
		conn.SetReadDeadline(Now().Add(Second*2))
		
		data := make([]byte, 16)
		
		n, _, err := conn.ReadFromUDP(data[0:])
		
		if err != nil {
			bossAlive = false
		} else {
			bossCount = computeMessage(string(data[0:n]))
			Println("Someone else is counting:", bossCount)
		}
	}
	conn.Close()
	
	
	Println("I am master")
	SpawnNewFriend()
	startTime = bossCount
	
	addr, _ := net.ResolveUDPAddr("udp4", "10.22.79.255:20011")
	UDPconn, err := net.DialUDP("udp4", nil, addr)
	if err != nil{
			Println("# ERROR:conn: ",err.Error())
	}
	
	
	for {		
		msg := "Bosscount:" + Itoa(bossCount)

		_, err := UDPconn.Write([]byte(msg))
		Println("I am counting:", bossCount)
		if err != nil {
			Println("# ERROR:Broadcast:", err.Error())
		}		
		
		if (bossCount == 13 + startTime) {
			break
		}
		bossCount++
		Sleep(Second)
	}
	UDPconn.Close()
	Println("I am done")
}

func computeMessage(msg string) int {
	num := TrimLeft(msg, "Boscunt: ")
	bosscount,err := Atoi(num)
	if (err != nil) { return -1}
	return bosscount
} 

func SpawnNewFriend() {
	Println("Spawning new \"friend\"")

    cmd := exec.Command("mate-terminal", "go", "run", "phoenix.go") // mellom mt og go "-x",
	out, err := cmd.Output()

    if err != nil {
        Println(err.Error())
        Println("error")
        return
    }

    print(string(out))
    
}