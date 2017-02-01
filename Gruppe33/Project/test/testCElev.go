package main

import "./driver"
import "time"

func main() {
    driver.Elev_init()
    driver.Elev_set_motor_direction(-1)
    time.Sleep(1 * time.Second)
    driver.Elev_set_motor_direction(0)
}
