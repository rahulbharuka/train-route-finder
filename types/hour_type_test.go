package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHourType(t *testing.T) {
	t.Run("non-peak-hour", func(t *testing.T) {
		journeyTime, _ := time.Parse("2006-01-02T15:04", "2020-04-12T15:04")
		assert.Equal(t, HTNonPeak, GetHourType(journeyTime))
	})

	t.Run("peak-hour", func(t *testing.T) {
		journeyTime, _ := time.Parse("2006-01-02T15:04", "2020-04-10T20:04")
		assert.Equal(t, HTPeak, GetHourType(journeyTime))
	})
}

func TestConvertToHourType(t *testing.T) {
	t.Run("happy-path", func(t *testing.T) {
		assert.Equal(t, HTPeak, ConvertToHourType("Peak"))
	})

	t.Run("invalid-hourType", func(t *testing.T) {
		assert.Equal(t, HTInvalid, ConvertToHourType("somethinh"))
	})
}
