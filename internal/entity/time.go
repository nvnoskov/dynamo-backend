package entity

import "time"

type FlightTime struct {
	*time.Time
}

func (t FlightTime) MarshalJSON() ([]byte, error) {
	return []byte(t.Format("\"2006-01-02T15:04:05Z\"")), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// The time is expected to be a quoted string in RFC 3339 format.
func (t *FlightTime) UnmarshalJSON(data []byte) (err error) {
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse("\"2006-01-02T15:04:05Z\"", string(data))
	*t = FlightTime{&tt}
	return
}
