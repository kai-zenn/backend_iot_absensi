package helper

import "time"

var wibLocation *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*60*60)
	}
	wibLocation = loc
}

func NowWIB() time.Time {
	return time.Now().In(wibLocation)
}

func StartOfDayWIB(t time.Time) time.Time {
	tWIB := t.In(wibLocation)
	return time.Date(tWIB.Year(), tWIB.Month(), tWIB.Day(), 0, 0, 0, 0, wibLocation)
}
