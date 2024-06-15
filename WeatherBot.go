package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Структуры для данных о погоде
type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

// Функция для получения погоды
func getWeather(city string, apiKey string) (WeatherData, error) {
	var weatherData WeatherData
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return weatherData, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weatherData, fmt.Errorf("failed to get weather data: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		return weatherData, err
	}

	return weatherData, nil
}

// Основная функция
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <TELEGRAM_BOT_TOKEN> <WEATHER_API_KEY>")
		return
	}

	botToken := os.Args[1]
	weatherApiKey := os.Args[2]

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			city := update.Message.Text
			weather, err := getWeather(city, weatherApiKey)
			var msgText string
			if err != nil {
				msgText = fmt.Sprintf("Error: %s", err)
			} else {
				msgText = fmt.Sprintf("Weather in %s: %.2f°C, %s", weather.Name, weather.Main.Temp, weather.Weather[0].Description)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			bot.Send(msg)
		}
	}
}
