package main

import "time"

const (
	TIMEZONE_OFFSET = -28800
)

func IsSameDay(t1, t2 int64) bool {
	t1 = t1 + TIMEZONE_OFFSET
	t2 = t2 + TIMEZONE_OFFSET
	return (t1 / 86400) == (t2 / 86400)
}

func GetTime() int64 {
	return time.Now().Unix()
}
