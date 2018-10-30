package main

import (
	"time"

	"github.com/jinzhu/now"
)

func yesterdayFrom() string {
	yesterday := time.Now().AddDate(0, 0, -1)
	startDate := now.New(yesterday).BeginningOfMonth()
	return startDate.Format("01-02-2006")
}

func yesterdayUntil() string {
	yesterday := time.Now().AddDate(0, 0, -1)
	endDate := now.New(yesterday).EndOfMonth()
	return endDate.Format("01-02-2006")
}

func defaultFrom() string {
	startDate := now.New(time.Now()).BeginningOfMonth()
	return startDate.Format("01-02-2006")
}

func defaultUntil() string {
	endDate := time.Now()
	return endDate.Format("01-02-2006")
}
