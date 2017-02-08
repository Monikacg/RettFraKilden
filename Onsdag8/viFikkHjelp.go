func Print(msg chan int, numChan chan int){
    for{
        select{
            case m := <-msg:
                fmt.Println(m)

            case order := <- buttonPressed:
                Queue.add(order)
                SetLight
                newChannel <- 1

        }

    }
}

go Print(...)
