package def

import "fmt"

const (

	N_ELEV				int = 1
	N_FLOORS			int = 4
	N_BUTTONS			int = 3
	BUTTON_CALL_UP		int = 0
	BUTTON_CALL_DOWN	int = 1
	BUTTON_COMMAND		int = 2
	DIR_UP				int = 1
	DIR_DOWN			int = -1
	DIR_STOP			int = 0
	ON 					int = 1
	OFF					int = 0

	//states
	IDLE 		int = 0
	MOVING 		int = 1
	DOOR_OPEN 	int = 2

	LOCAL_LISTEN_PORT		int = 20005
	BROADCAST_LISTEN_PORT	int = 30005
	MESSAGE_SIZE	 		int = 1024

	NOTHING 		int = 0
	ADD_ORDERS 		int = 1
	REMOVE_ORDERS 	int = 2
)



type MSG struct{
	MsgType			int
	State 			int
	PrevFloor 		int
	Dir   			int 	//never 0.
	ExUpOrders 		[N_FLOORS]int
	ExDownOrders	[N_FLOORS]int
	InOrders		[N_FLOORS]int
}

type Order struct{
	Floor 	int
	Button 	int
}

type Udp_message struct {
	Raddr  string //if receiving raddr=senders address, if sending raddr should be set to "broadcast" or an ip:port
	Data   []byte //TODO: implement another encoding, strings are meh
	Length int    //length of received data, in #bytes // N/A for sending
}

var Msg = MSG{}

var buff = make([]byte, 1024)
var Udp_msg = Udp_message{"broadcast", buff, 1024}

// func (msg MSG) String() string{
// 	return fmt.Sprintf()
// }

func PrintMsg() {
	fmt.Println()

	for i:=0;i<N_FLOORS;i++ {
		defer fmt.Println(Msg.ExDownOrders[i], " " ,Msg.ExUpOrders[i], " ", Msg.InOrders[i])
	}
	switch Msg.State {
	case IDLE:
		fmt.Println("State: IDLE")
	case MOVING:
		fmt.Println("State: MOVING")
	case DOOR_OPEN:
		fmt.Println("State: DOOR_OPEN")
	default:
		fmt.Println("Invalid state: ", Msg.State)
	}

	fmt.Println("Floor: ", Msg.PrevFloor)

	switch Msg.Dir {
		case DIR_UP:
			fmt.Println("DIR UP")
		case DIR_DOWN:
			fmt.Println("DIR DOWN")
		default:
			fmt.Println("DIR 0, ERROR!!!")
	}
}
