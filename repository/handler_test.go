package repository

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

// setup rail network
func setup() {
	os.Setenv("STATION_MAP_FILE", "../StationMap.csv")
	os.Setenv("TRAINLINE_COST_FILE", "../trainline_cost.csv")
	os.Setenv("INTERCHANGE_COST_FILE", "../interchange_cost.csv")
	os.Setenv("MAX_ROUTES", "3")
	RailNetworkInit()
}

func TestFindRoutes(t *testing.T) {
	h := GetHandler()

	t.Run("invalid-src", func(t *testing.T) {
		routes, err := h.FindRoutes("Wonderland", "Bugis", time.Time{}, false)
		assert.EqualError(t, ErrInvalidRequest, err.Error())
		assert.Nil(t, routes)
	})

	t.Run("simple-routes", func(t *testing.T) {
		expectedRoutes := []*Route{
			&Route{
				Heading: "Number of stops to destination: 7",
				Steps:   "Take CC line from Holland Village to Botanic Gardens. Change from CC line to DT line. Take DT line from Botanic Gardens to Bugis.",
			},
			&Route{
				Heading: "Number of stops to destination: 8",
				Steps:   "Take CC line from Holland Village to Botanic Gardens. Change from CC line to DT line. Take DT line from Botanic Gardens to Little India. Change from DT line to NE line. Take NE line from Little India to Dhoby Ghaut. Change from NE line to NS line. Take NS line from Dhoby Ghaut to City Hall. Change from NS line to EW line. Take EW line from City Hall to Bugis.",
			},
			&Route{
				Heading: "Number of stops to destination: 9",
				Steps:   "Take CC line from Holland Village to Caldecott. Change from CC line to TE line. Take TE line from Caldecott to Stevens. Change from TE line to DT line. Take DT line from Stevens to Bugis.",
			},
		}

		routes, err := h.FindRoutes("Holland Village", "Bugis", time.Time{}, false)
		assert.NoError(t, err)
		for i, route := range routes {
			assert.Equal(t, expectedRoutes[i].Heading, route.Heading)
			assert.Equal(t, expectedRoutes[i].Steps, route.Steps)
		}
	})

	t.Run("realtime-routes", func(t *testing.T) {
		journeyTime, _ := time.Parse("2006-01-02T15:04", "2019-01-31T19:00")

		expectedRoutes := []*Route{
			&Route{
				Heading: "Expected Travel time: 150",
				Steps:   "Take EW line from Boon Lay to Buona Vista. Change from EW line to CC line. Take CC line from Buona Vista to Botanic Gardens. Change from CC line to DT line. Take DT line from Botanic Gardens to Little India.",
			},
			&Route{
				Heading: "Expected Travel time: 173",
				Steps:   "Take EW line from Boon Lay to Outram Park. Change from EW line to NE line. Take NE line from Outram Park to Little India.",
			},
			&Route{
				Heading: "Expected Travel time: 185",
				Steps:   "Take EW line from Boon Lay to Buona Vista. Change from EW line to CC line. Take CC line from Buona Vista to Caldecott. Change from CC line to TE line. Take TE line from Caldecott to Stevens. Change from TE line to DT line. Take DT line from Stevens to Little India.",
			},
		}

		routes, err := h.FindRoutes("Boon Lay", "Little India", journeyTime, true)
		assert.NoError(t, err)
		for i, route := range routes {
			assert.Equal(t, expectedRoutes[i].Heading, route.Heading)
			assert.Equal(t, expectedRoutes[i].Steps, route.Steps)
		}
	})
}
