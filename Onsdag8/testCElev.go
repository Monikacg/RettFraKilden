package main

import "./driver"
import "time"

func main() {
    //tester at heisen kjører
    driver.Elev_init()
    driver.Elev_set_motor_direction(1)
    time.Sleep(1 * time.Second)
    driver.Elev_set_motor_direction(0)

    //tester alle lysene :-)
    for i := 0; i < 4; i++ {
        driver.Elev_set_button_lamp(0, i, 0)
        driver.Elev_set_button_lamp(1, i, 0)
        driver.Elev_set_button_lamp(2, i, 0)
        driver.Elev_set_floor_indicator(i)
        time.Sleep(1*time.Second)
    }
    driver.Elev_set_door_open_lamp(0)
}
