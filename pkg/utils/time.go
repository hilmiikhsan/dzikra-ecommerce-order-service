package utils

import (
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/constants"
)

const TimeFormat = "2006-01-02T15:04:05.000Z"

func ParseTime(value string) (time.Time, error) {
	return time.Parse(TimeFormat, value)
}

func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

func FormatToWIB(t time.Time) string {
	// WIB (UTC+7)
	wibZone := time.FixedZone("WIB", 7*3600)
	return t.In(wibZone).Format("2006-01-02T15:04:05-07:00")
}

func FormatTimeJakarta() time.Time {
	loc, err := time.LoadLocation(constants.TimeLocationAsiaJakarta)
	if err != nil {
		return time.Now().UTC()
	}
	return time.Now().In(loc)
}
