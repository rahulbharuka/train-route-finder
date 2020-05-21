package types

import "time"

// HourType ...
type HourType int

const (
	// HTInvalid ...
	HTInvalid HourType = iota
	// HTPeak ...
	HTPeak
	// HTNonPeak ...
	HTNonPeak
	// HTNight ...
	HTNight
)

// GetHourType returns hour type for given time
func GetHourType(t time.Time) HourType {
	weekday := t.Weekday()
	hrs := t.Hour()

	if hrs >= 22 && hrs < 6 {
		return HTNight
	}

	if (weekday >= 1 && weekday <= 6) && ((hrs >= 6 && hrs < 9) || (hrs >= 18 && hrs < 21)) {
		return HTPeak
	}
	return HTNonPeak
}

// ConvertToHourType converts string to HourType
func ConvertToHourType(str string) HourType {
	switch str {
	case "Peak":
		return HTPeak
	case "NonPeak":
		return HTNonPeak
	case "Night":
		return HTNight
	}
	return HTInvalid
}
