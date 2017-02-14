package consoleIO

import (
	"bufio"
	."fmt"
	."strings"
	"os"
)

func InitConsoleIO(in chan string, ut chan string){
	q := make(chan int)
	go PrintConsole(in)
	//go ConsoleInput(ut)
	<-ut
	<-q
}

func ConsoleInput(inn chan string){
	for { 
		//Println("Skriv en beskjed: ")
		msg := GetInput("Skriv en beskjed: ")
		inn <- msg	
	}
}

func PrintConsole(ut chan string) {
	var msg string
	for {
		msg = <-ut
		Println("[CONSOLE] ",msg)
	}	
}

func GetInput(prompt string) string {
	// saving a string in input according to prompt 
	reader := bufio.NewReader(os.Stdin)
	Print(prompt)
	input,_ := reader.ReadString('\n')
	
	return TrimRight(input, "\n")	
}