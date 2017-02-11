package definitions

/* NB! These constants are defined to fit with the functions in elev.h.
The constants have not been directly imported from C simply because we
didn't find a way to when writing this. One obvious improvement of the
code would be to import (some) of these constants in order to fix a
weakness in the system.
*/

const (
  // Number of floors/buttons
  N_FLOORS = 4
  N_BUTTONS = 3
  MAX_N_LIFTS = 3

  // Lift states
  INIT = -1
  IDLE = 0
  MOVING = 1
  DOOR_OPEN = 2

  // Button calls
  BUTTON_CALL_UP = 0
  BUTTON_CALL_DOWN = 1
  BUTTON_COMMAND = 2

  // Motor directions
  DIRN_DOWN = -1
  DIRN_STOP = 0
  DIRN_UP = 1


  ON = 1
  OFF = 0
)

type Button struct { // Brukes på button_outside_chan. Kanskje endre til å bruke DIRN_UP og DIRN_DOWN?
  floor int
  button_dir int
}

type Order struct { // Brukes på local_order_chan
  cat string // "LIGHT"/"DIR"
  order int  // DIRN_DOWN/UP/STOP, BUTTON_CALL_UP/DOWN/COMMAND
  floor int //0-3 (0-N_FLOORS)
  value int // ON/OFF for lys, settes bare for "LIGHT"
} // Floor trengs ikke på doorlight, value trengs ikke på retn.

type Properties_struct struct { // Brukes i get_properties. ENDRE NAVN!
  last_floor int
  dirn int
  state int
}

type Udp struct {
  id int
  message string
  Properties_struct
}

// NB! GÅ gjennom og endre til stor bokstav i koden også.