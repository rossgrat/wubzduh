package main

import (
	"fmt"
	"time"

	wubzduh "github.com/rossgrat/wubzduh/src"
)

const (
	time1 = iota //12:30am Next Day
	time2        //11:30pm Same Day
)

const (
	numTimes = 2
)

var FetchTimer *time.Timer
var PurgeTimer *time.Timer

func FetchThread(timeCase int) {

	currentTime := time.Now()
	var diff time.Duration
	switch timeCase {
	case time1:
		// 12:30am the next day
		currentTimePlusDay := currentTime.AddDate(0, 0, 1)
		targetTime := time.Date(currentTimePlusDay.Year(), currentTimePlusDay.Month(), currentTimePlusDay.Day(), 0, 30, 0, 0, currentTimePlusDay.Location())
		diff = targetTime.Sub(currentTime)
		fmt.Printf("Set Fetch Timer for 12:30am Tomorrow. ETE: %s\n", diff)
	case time2:
		//11:30pm the same day
		targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 30, 0, 0, currentTime.Location())
		diff = targetTime.Sub(currentTime)
		fmt.Printf("Set Fetch Timer for 11:30pm Today. ETE: %s\n", diff)
	}

	FetchTimer = time.NewTimer(diff)
	<-FetchTimer.C
	wubzduh.Fetch(DB)
	go FetchThread((timeCase + 1) % numTimes)
}

func PurgeThread() {
	currentTime := time.Now()

	currentTimePlusDay := currentTime.AddDate(0, 0, 1)
	targetTime := time.Date(currentTimePlusDay.Year(), currentTimePlusDay.Month(), currentTimePlusDay.Day(), 12, 0, 0, 0, currentTimePlusDay.Location())
	diff := targetTime.Sub(currentTime)

	PurgeTimer = time.NewTimer(diff)
	<-PurgeTimer.C
	wubzduh.Purge(DB)
	go PurgeThread()
}
