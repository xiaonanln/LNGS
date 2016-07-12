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

func MaxInt(i1 int, i2 int) int {
	if i1 >= i2 {
		return i1
	} else {
		return i2
	}
}

func MinInt(i1 int, i2 int) int {
	if i1 <= i2 {
		return i1
	} else {
		return i2
	}
}

func AbsInt(a int) int {
	if a >= 0 {
		return a
	} else {
		return -a
	}
}
