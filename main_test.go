package weather

import "testing"

func TestWeather(t *testing.T) {

	latitude := 55.6302
	longitude := 37.6045
	var forecastDays uint8 = 1

	t.Run("weather", func(t *testing.T) {
		weatherConfig := Config{
			latitude:     latitude,
			longitude:    longitude,
			forecastDays: forecastDays,
			tomorrow:     false,
		}
		message, err := Weather(weatherConfig)
		if err != nil {
			t.Error(err.Error())
		}
		if message == "" {
			t.Error("empty response")
		}
	})

	t.Run("weather", func(t *testing.T) {
		weatherConfig := Config{
			latitude:     latitude,
			longitude:    longitude,
			forecastDays: 2,
			tomorrow:     true,
		}
		message, err := Weather(weatherConfig)
		if err != nil {
			t.Error(err.Error())
		}
		if message == "" {
			t.Error("empty response")
		}
	})
}
