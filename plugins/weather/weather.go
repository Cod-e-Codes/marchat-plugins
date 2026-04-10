package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Cod-e-Codes/marchat/plugin/sdk"
)

// WeatherPlugin provides weather information using wttr.in
type WeatherPlugin struct {
	*sdk.BasePlugin
	config sdk.Config
}

// NewWeatherPlugin creates a new weather plugin
func NewWeatherPlugin() *WeatherPlugin {
	return &WeatherPlugin{
		BasePlugin: sdk.NewBasePlugin("weather"),
	}
}

// Init initializes the weather plugin
func (p *WeatherPlugin) Init(config sdk.Config) error {
	p.config = config
	log.Printf("Weather plugin initialized")
	return nil
}

// OnMessage handles incoming messages
func (p *WeatherPlugin) OnMessage(msg sdk.Message) ([]sdk.Message, error) {
	lower := strings.ToLower(msg.Content)
	if strings.HasPrefix(lower, "weather:") || strings.HasPrefix(lower, "weather ") {
		location := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(lower, "weather:"), "weather "))
		if location == "" {
			location = "Charlotte,NC"
		}

		weather, err := p.getWeather(location, false)
		if err != nil {
			return []sdk.Message{{
				Sender:    "WeatherBot",
				Content:   fmt.Sprintf("Failed to get weather: %v", err),
				CreatedAt: time.Now(),
			}}, nil
		}

		return []sdk.Message{{
			Sender:    "WeatherBot",
			Content:   weather,
			CreatedAt: time.Now(),
		}}, nil
	}
	return nil, nil
}

// Commands returns the commands this plugin provides
func (p *WeatherPlugin) Commands() []sdk.PluginCommand {
	return []sdk.PluginCommand{
		{
			Name:        "weather",
			Description: "Get current weather for a location",
			Usage:       ":weather [location]",
			AdminOnly:   false,
		},
		{
			Name:        "forecast",
			Description: "Get weather forecast for a location",
			Usage:       ":forecast [location]",
			AdminOnly:   false,
		},
	}
}

// getWeather fetches weather data from wttr.in
func (p *WeatherPlugin) getWeather(location string, forecast bool) (string, error) {
	url := fmt.Sprintf("https://wttr.in/%s?format=j1", location)

	log.Printf("Fetching weather from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("weather service returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to parse weather data: %w", err)
	}

	currentCondition, ok := data["current_condition"].([]interface{})
	if !ok || len(currentCondition) == 0 {
		return "", fmt.Errorf("no current condition data")
	}

	current := currentCondition[0].(map[string]interface{})

	temp := current["temp_C"].(string)
	feelsLike := current["FeelsLikeC"].(string)
	weatherDesc := current["weatherDesc"].([]interface{})[0].(map[string]interface{})["value"].(string)
	humidity := current["humidity"].(string)
	windSpeed := current["windspeedKmph"].(string)

	nearestArea := data["nearest_area"].([]interface{})[0].(map[string]interface{})
	areaName := nearestArea["areaName"].([]interface{})[0].(map[string]interface{})["value"].(string)
	country := nearestArea["country"].([]interface{})[0].(map[string]interface{})["value"].(string)

	result := fmt.Sprintf("🌤️  Weather for %s, %s\n", areaName, country)
	result += fmt.Sprintf("Condition: %s\n", weatherDesc)
	result += fmt.Sprintf("Temperature: %s°C (Feels like %s°C)\n", temp, feelsLike)
	result += fmt.Sprintf("Humidity: %s%%\n", humidity)
	result += fmt.Sprintf("Wind Speed: %s km/h", windSpeed)

	if forecast {
		weather, ok := data["weather"].([]interface{})
		if ok && len(weather) > 0 {
			result += "\n\n📅 3-Day Forecast:\n"
			for i, day := range weather {
				if i >= 3 {
					break
				}
				dayData := day.(map[string]interface{})
				date := dayData["date"].(string)
				maxTemp := dayData["maxtempC"].(string)
				minTemp := dayData["mintempC"].(string)
				hourly := dayData["hourly"].([]interface{})
				if len(hourly) > 0 {
					desc := hourly[0].(map[string]interface{})["weatherDesc"].([]interface{})[0].(map[string]interface{})["value"].(string)
					result += fmt.Sprintf("%s: %s, High: %s°C, Low: %s°C\n", date, desc, maxTemp, minTemp)
				}
			}
		}
	}

	return result, nil
}

func main() {
	plugin := NewWeatherPlugin()
	if err := sdk.RunStdio(plugin, plugin.handleCommand); err != nil {
		log.Fatalf("plugin exited: %v", err)
	}
}

func (p *WeatherPlugin) handleCommand(command string, args []string) sdk.PluginResponse {
	location := "Charlotte,NC"
	if len(args) > 0 {
		location = strings.Join(args, " ")
	}

	var weather string
	var err error

	switch command {
	case "weather":
		weather, err = p.getWeather(location, false)
	case "forecast":
		weather, err = p.getWeather(location, true)
	default:
		return sdk.PluginResponse{
			Type:    "command",
			Success: false,
			Error:   "unknown command",
		}
	}

	if err != nil {
		return sdk.PluginResponse{
			Type:    "command",
			Success: false,
			Error:   fmt.Sprintf("failed to get weather: %v", err),
		}
	}

	weatherMsg := sdk.Message{
		Sender:    "WeatherBot",
		Content:   weather,
		CreatedAt: time.Now(),
	}

	responseData, err := json.Marshal(weatherMsg)
	if err != nil {
		return sdk.PluginResponse{
			Type:    "command",
			Success: false,
			Error:   err.Error(),
		}
	}
	return sdk.PluginResponse{
		Type:    "message",
		Success: true,
		Data:    responseData,
	}
}
