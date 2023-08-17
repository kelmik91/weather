package weather

import "testing"

func TestWeather(t *testing.T) {

	latitude := 55.6302
	longitude := 37.6045
	forecastDays := "1"

	t.Run("weather", func(t *testing.T) {
		message, err := Weather(latitude, longitude, forecastDays)
		if err != nil {
			t.Error(err.Error())
		}
		if message == "" {
			t.Error("empty response")
		}
	})
}
