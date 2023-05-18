package timecheck

import (
	"log"
	"time"

	"github.com/beevik/ntp"
)

const (
	ntpServer = "pool.ntp.org"
)

func GetTimeDifference() (time.Duration, error) {
	ntpTime, err := ntp.Time(ntpServer)
	if err != nil {
		return 0, err
	}
	log.Printf("NTP time: %s", ntpTime)

	return time.Since(ntpTime), nil
}
