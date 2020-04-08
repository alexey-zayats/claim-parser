package fsdump

import (
	"fmt"
	"strings"
	"time"
)

// Time ...
type Time struct {
	time.Time
}

// Now ...
func Now() Time {
	return Time{
		time.Now(),
	}
}

const timeFormat = "2006-01-02T15:04:05"
const timeLayout = "02.01.2006, 15:04:05"

// String ...
func (t Time) String() string {
	return t.Time.String()
}

// UnmarshalJSON ...
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	value := strings.Trim(string(b), "\"")
	if value == "null" {
		t.Time = time.Time{}
		return
	}

	t.Time, err = time.Parse(timeLayout, value)
	return
}

// MarshalJSON ...
func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time.UnixNano() == 0 {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(timeFormat))), nil
}
