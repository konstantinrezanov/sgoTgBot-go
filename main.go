package main

import (
	"log"
	"sgoTgBot-go/bot"
	"sgoTgBot-go/conversion"
	"sync"
	"time"
)

var schedule conversion.Schedule
var m sync.Mutex
var wg sync.WaitGroup

func init() {
	schedule.AssignSchedule(&m)
}
func main() {
	wg.Add(1)

	parseTicker := time.NewTicker(15 * time.Minute)
	clearTicker := time.NewTicker(6 * time.Hour)
	done := make(chan bool)

	go bot.StartBot(schedule, &m)
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-parseTicker.C:
				schedule.AssignSchedule(&m)
			}

		}
	}()
	go func() {
		for {
			select {
			case <-done:
				return
			case _ = <-clearTicker.C:
				schedule.ScheduleClear(7, &m)
			}
		}
	}()
	wg.Wait()
}

