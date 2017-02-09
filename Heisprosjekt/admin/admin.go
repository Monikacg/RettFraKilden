package admin

import (
  "fmt"
)

func admin()  {
  orders = create_order_matrix();
  if s := system_state; s == slave {
    //Ask master for orders and copy in those in the orders matrix.
  }
}
