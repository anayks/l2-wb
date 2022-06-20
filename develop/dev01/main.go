package main

import (
	"time"

	ntp "github.com/beevik/ntp"
)

func GetNowTime() time.Time { // настоящее время
	options := ntp.QueryOptions{Timeout: 5 * time.Second, TTL: 25}
	response, err := ntp.QueryWithOptions("0.beevik-ntp.pool.ntp.org", options)
	if err != nil {
		panic(err)
	}
	return response.Time
}

func GetExactTime() time.Time { // точное время
	options := ntp.QueryOptions{Timeout: 5 * time.Second, TTL: 25}
	response, err := ntp.QueryWithOptions("0.beevik-ntp.pool.ntp.org", options)
	if err != nil {
		panic(err)
	}
	return time.Now().Add(response.ClockOffset)
}
