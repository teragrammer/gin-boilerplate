package utilities

import (
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

func TestAddDay(t *testing.T) {
	layout := "2006-01-02 15:04"
	currentDateTime := "2023-06-10 13:45"
	addedDays := "2023-06-12 13:45"

	// Parse the MySQL datetime string into a time.Time object
	parsedTime, err := time.Parse(layout, currentDateTime)
	if err != nil {
		t.Error(err)
		return
	}

	addedDateTime := AddDay(parsedTime, 2)
	parsedAddedDateTime, err := time.Parse(layout, addedDateTime.Format(layout))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, parsedAddedDateTime.Format(layout), addedDays)
}

func TestSubtractDay(t *testing.T) {
	layout := "2006-01-02 15:04"
	currentDateTime := "2023-06-10 13:45"
	subtractDays := "2023-06-08 13:45"

	// Parse the MySQL datetime string into a time.Time object
	parsedTime, err := time.Parse(layout, currentDateTime)
	if err != nil {
		t.Error(err)
		return
	}

	addedDateTime := SubtractDay(parsedTime, 2)
	parsedAddedDateTime, err := time.Parse(layout, addedDateTime.Format(layout))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, parsedAddedDateTime.Format(layout), subtractDays)
}

func TestAddHour(t *testing.T) {
	layout := "2006-01-02 15:04"
	currentDateTime := "2023-06-15 10:45"
	addedHours := "2023-06-15 12:45"

	// Parse the MySQL datetime string into a time.Time object
	parsedTime, err := time.Parse(layout, currentDateTime)
	if err != nil {
		t.Error(err)
		return
	}

	addedDateTime := AddHour(parsedTime, 2)
	parsedAddedDateTime, err := time.Parse(layout, addedDateTime.Format(layout))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, parsedAddedDateTime.Format(layout), addedHours)
}

func TestSubtractHour(t *testing.T) {
	layout := "2006-01-02 15:04"
	currentDateTime := "2023-06-15 10:45"
	subtractedHours := "2023-06-15 08:45"

	// Parse the MySQL datetime string into a time.Time object
	parsedTime, err := time.Parse(layout, currentDateTime)
	if err != nil {
		t.Error(err)
		return
	}

	addedDateTime := SubtractHour(parsedTime, 2)
	parsedAddedDateTime, err := time.Parse(layout, addedDateTime.Format(layout))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, parsedAddedDateTime.Format(layout), subtractedHours)
}

func TestAddMinute(t *testing.T) {
	layout := "2006-01-02 15:04"
	currentDateTime := "2023-06-15 13:10"
	addedMinutes := "2023-06-15 13:15"

	// Parse the MySQL datetime string into a time.Time object
	parsedTime, err := time.Parse(layout, currentDateTime)
	if err != nil {
		t.Error(err)
		return
	}

	addedDateTime := AddMinute(parsedTime, 5)
	parsedAddedDateTime, err := time.Parse(layout, addedDateTime.Format(layout))
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, parsedAddedDateTime.Format(layout), addedMinutes)
}
