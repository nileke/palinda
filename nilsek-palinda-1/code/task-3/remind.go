package main

import (
	"fmt"
	"time"
)

func Remind(text string, delay time.Duration) {
	for {
		time.Sleep(delay)
		fmt.Println("Klockan är: ", time.Now().Format("15.04"), text)
	}
}

func main() {
	// set duration of delay
	eat := 3 * time.Hour
	work := 8 * time.Hour
	sleep := 24 * time.Hour

	// start threads
	go Remind("Go eat something", eat)
	go Remind("Time to work", work)
	go Remind("Nap time!", sleep)
	select {}
}
