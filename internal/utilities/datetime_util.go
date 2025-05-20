package utilities

import "time"

func DatetimeUTC(datetime string, layout string, timeZone string) (time.Time, error) {
	// Parse the time string into a time.Time object
	loc, _ := time.LoadLocation(timeZone) // Use "Local" for local time
	localTime, err := time.ParseInLocation(layout, datetime, loc)
	if err != nil {
		return time.Now(), err
	}

	// Convert local time to UTC
	return localTime.UTC(), nil
}

func AddDay(datetime time.Time, day int) time.Time {
	// Current date and time
	// time.Now()
	return datetime.AddDate(0, 0, day)
}
