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

func SubtractDay(datetime time.Time, day int) time.Time {
	// Current date and time
	// time.Now()
	return datetime.AddDate(0, 0, -day)
}

func AddHour(datetime time.Time, hour int64) time.Time {
	// Current date and time
	// time.Now()
	return datetime.Add(time.Duration(hour) * time.Hour)
}

func SubtractHour(datetime time.Time, hour int64) time.Time {
	// Current date and time
	// time.Now()
	return datetime.Add(-time.Duration(hour) * time.Hour)
}

func AddMinute(datetime time.Time, minute int64) time.Time {
	// Current date and time
	// time.Now()
	return datetime.Add(time.Duration(minute) * time.Minute)
}

func SubtractMinute(datetime time.Time, minute int64) time.Time {
	// Current date and time
	// time.Now()
	return datetime.Add(-time.Duration(minute) * time.Minute)
}

// FromDateTimeToDate takes a date-time string in the format "2006-01-02 15:04:05"
// and returns the date part as a string in the format "2006-01-02".
// sample usage: FromDateTimeToDate("your-date-time")
// it converts to UTC format
func FromDateTimeToDate(datetime string, timeZone string) (string, error) {
	// Parse the time string into a time.Time object
	loc, _ := time.LoadLocation(timeZone) // Use "Local" for local time
	localTime, err := time.ParseInLocation("2006-01-02 15:04:05", datetime, loc)
	if err != nil {
		return time.Now().String(), err
	}

	// Format the time.Time object to a string with the desired date format
	dateStr := localTime.UTC().Format("2006-01-02")
	return dateStr, nil
}
