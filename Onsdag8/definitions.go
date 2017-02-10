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

  ON = 1
  OFF = 0
)