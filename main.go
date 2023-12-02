package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type weatherJson struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationtimeMs     float64 `json:"generationtime_ms"`
	UtcOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`
	CurrentWeather       struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		Winddirection float64 `json:"winddirection"`
		Weathercode   int     `json:"weathercode"`
		IsDay         int     `json:"is_day"`
		Time          int     `json:"time"`
	} `json:"current_weather"`
	DailyUnits struct {
		Time             string `json:"time"`
		Weathercode      string `json:"weathercode"`
		Temperature2MMax string `json:"temperature_2m_max"`
		Temperature2MMin string `json:"temperature_2m_min"`
		Sunrise          string `json:"sunrise"`
		Sunset           string `json:"sunset"`
	} `json:"daily_units"`
	Daily struct {
		Time             []int     `json:"time"`
		Weathercode      []int     `json:"weathercode"`
		Temperature2MMax []float64 `json:"temperature_2m_max"`
		Temperature2MMin []float64 `json:"temperature_2m_min"`
		Sunrise          []int     `json:"sunrise"`
		Sunset           []int     `json:"sunset"`
	} `json:"daily"`
}

type Config struct {
	latitude     float64
	longitude    float64
	forecastDays uint8
	tomorrow     bool
}

var weatherCode = map[int]string{
	0:  "☀️Чистое небо ☀️",
	1:  "☀️ В основном ясно ☀️",
	2:  "🌤 Переменная облачность 🌤",
	3:  "🌥 Пасмурная погода 🌥",
	45: "🌫Туман🌫",
	48: "🌫❄️ Туман с изморозью ❄️🌫",
	51: "🌧 Моросящий дождь: Легкий 🌧",
	53: "🌧🌧 Моросящий дождь: Умеренный 🌧🌧",
	55: "🌧🌧🌧 Моросящий дождь: Густой и интенсивный 🌧🌧🌧",
	56: "🌧❄️ Ледяной моросящий дождь: Легкий ❄️🌧",
	57: "🌧❄️🌧 Ледяной моросящий дождь:  Густой и интенсивный 🌧❄️🌧",
	61: "🌧 Дождь: Небольшой 🌧",
	63: "🌧🌧 Дождь: Умеренный 🌧🌧",
	65: "🌧🌧🌧 Дождь: Высокая интенсивность 🌧🌧🌧",
	66: "🌨❄️ Ледяной дождь: Легкий ❄️🌨",
	67: "🌨❄️🌨 Ледяной дождь: Высокая интенсивность 🌨❄️🌨",
	71: "❄️ Снег: Небольшой ❄️",
	73: "❄️❄️ Снег: Умеренный ❄️❄️",
	75: "❄️❄️❄️ Снег: Высокая интенсивность ❄️❄️❄️",
	77: "❄️ Снежные крупинки ❄️",
	80: "🌩 Ливневые дожди: Легкий 🌩",
	81: "🌩🌩 Ливневые дожди: Умеренный 🌩🌩",
	82: "🌩🌩🌩 Ливневые дожди: Сильные 🌩🌩🌩",
	85: "🌨 Снегопад: Легкий 🌨",
	86: "🌨🌨Снегопад: Тяжелый 🌨🌨",
	95: "⚡ Гроза: Слабая или умеренная ⚡",
	96: "⚡🧊 Гроза со слабым градом 🧊⚡",
	99: "⚡🧊🧊⚡Гроза с сильным градом ⚡🧊🧊⚡",
}

func Weather(config Config) (string, error) {
	if config.forecastDays == 0 {
		return "", errors.New("forecastDays is empty")
	}

	builder := &strings.Builder{}
	builder.WriteString("https://api.open-meteo.com/v1/forecast?")

	builder.WriteString("latitude=")
	builder.WriteString(fmt.Sprint(config.latitude))

	builder.WriteString("&longitude=")
	builder.WriteString(fmt.Sprint(config.longitude))

	builder.WriteString("&daily=weathercode,temperature_2m_max,temperature_2m_min,sunrise,sunset")

	builder.WriteString("&current_weather=true")
	builder.WriteString("&windspeed_unit=ms")
	builder.WriteString("&timeformat=unixtime")
	builder.WriteString("&timezone=Europe%2FMoscow")

	builder.WriteString("&forecast_days=")
	builder.WriteString(fmt.Sprint(config.forecastDays))

	get, err := http.Get(builder.String())
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(get.Body)

	getResp, _ := io.ReadAll(get.Body)
	weather := weatherJson{}
	err = json.Unmarshal(getResp, &weather)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	loc, _ := time.LoadLocation(weather.Timezone)
	sunrise := time.Unix(int64(weather.Daily.Sunrise[0]), 0)
	sunset := time.Unix(int64(weather.Daily.Sunset[0]), 0)

	var message string
	message += fmt.Sprint(time.Now().In(loc).Format("02.01.2006")) + " \n"
	message += "Текущая температура: " + fmt.Sprint(weather.CurrentWeather.Temperature) + "°C \n"
	message += "Скорость ветра: " + fmt.Sprint(weather.CurrentWeather.Windspeed) + " m/s \n"
	message += "Погода: " + fmt.Sprint(weatherCode[weather.Daily.Weathercode[0]]) + " \n"
	message += "Максимальная температура: " + fmt.Sprint(weather.Daily.Temperature2MMax[0]) + "°C \n"
	message += "Минимальная температура: " + fmt.Sprint(weather.Daily.Temperature2MMin[0]) + "°C \n"
	message += "Расвет: " + fmt.Sprint(sunrise.Format("15:04")) + " \n"
	message += "Закат: " + fmt.Sprint(sunset.Format("15:04")) + " \n"
	message += "Световой день: " + fmt.Sprint(sunset.Sub(sunrise))

	if (time.Now().In(loc).Hour() >= 18 && config.forecastDays >= 2) || (config.tomorrow && config.forecastDays >= 2) {

		sunriseTomorrow := time.Unix(int64(weather.Daily.Sunrise[1]), 0)
		sunsetTomorrow := time.Unix(int64(weather.Daily.Sunset[1]), 0)

		message += "\n"
		message += "\n"
		message += "Завтра " + fmt.Sprint(time.Now().In(loc).AddDate(0, 0, 1).Format("02.01.2006")) + " ожидается\n"
		message += "Погода: " + fmt.Sprint(weatherCode[weather.Daily.Weathercode[1]]) + " \n"
		message += "Максимальная температура: " + fmt.Sprint(weather.Daily.Temperature2MMax[1]) + "°C \n"
		message += "Минимальная температура: " + fmt.Sprint(weather.Daily.Temperature2MMin[1]) + "°C \n"
		message += "Расвет: " + fmt.Sprint(sunriseTomorrow.Format("15:04")) + " \n"
		message += "Закат: " + fmt.Sprint(sunsetTomorrow.Format("15:04")) + " \n"
		message += "Световой день: " + fmt.Sprint(sunsetTomorrow.Sub(sunriseTomorrow))
	}

	return message, nil
}
