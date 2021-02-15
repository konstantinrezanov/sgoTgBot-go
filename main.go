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
	updateTicker:=time.NewTicker(1*time.Minute)
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
	
	go func() {
		for {
			select {
			case <-done:
				return
			case _=<-updateTicker.C:
				checkTime(&parseTicker)
			}
		}
	}()
	wg.Wait()
}

func checkTime(ticker **time.Ticker) {
	if time.Now().Hour()>21 || time.Now().Weekday()==time.Sunday {
		*ticker=time.NewTicker(time.Hour)
		log.Println("Set interval: 1 Hour")
		log.Println(*ticker)
	} else if time.Now().Hour() < 12 {
		*ticker=time.NewTicker(30*time.Minute)
		log.Println("Set interval: 30 Minutes")
		log.Println(*ticker)
	} else {
		*ticker=time.NewTicker(15*time.Minute)
		log.Println("Set interval: 15 minutes")
		log.Println(*ticker)
	}
}
