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

  // Lift states
  INIT = -1
  IDLE = 0
  MOVING = 1
  DOOR_OPEN = 2

  // Lift directions
  BUTTON_CALL_UP = 0
  BUTTON_CALL_DOWN = 1
  BUTTON_COMMAND = 2

  // Motor directions
  DIRN_DOWN = -1
  DIRN_STOP = 0
  DIRN_UP = 1

  type button struct { // Brukes på button_outside_chan
    direction string
    floor int
  }

  type order struct { // Brukes på local_order_chan
    cat string // "LIGHT"/"DIR"
    order string  // DIRN_DOWN/UP/STOP, BUTTON_CALL_UP/DOWN/COMMAND
    floor int //0-3 (0-N_FLOORS)
    value int // ON/OFF for lys, settes bare for "LIGHT"
  } // Floor trengs ikke på doorlight, value trengs ikke på retn.

  ON = 1
  OFF = 0
)
