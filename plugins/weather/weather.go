package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	// Check for weather queries in messages
	lower := strings.ToLower(msg.Content)
	if strings.HasPrefix(lower, "weather:") || strings.HasPrefix(lower, "weather ") {
		location := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(lower, "weather:"), "weather "))
		if location == "" {
			location = "Charlotte,NC" // Default location
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
	// Build URL
	url := fmt.Sprintf("https://wttr.in/%s?format=j1", location)

	log.Printf("Fetching weather from: %s", url)

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("weather service returned status %d", resp.StatusCode)
	}

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to parse weather data: %w", err)
	}

	// Extract current conditions
	currentCondition, ok := data["current_condition"].([]interface{})
	if !ok || len(currentCondition) == 0 {
		return "", fmt.Errorf("no current condition data")
	}

	current := currentCondition[0].(map[string]interface{})

	// Format weather information
	temp := current["temp_C"].(string)
	feelsLike := current["FeelsLikeC"].(string)
	weatherDesc := current["weatherDesc"].([]interface{})[0].(map[string]interface{})["value"].(string)
	humidity := current["humidity"].(string)
	windSpeed := current["windspeedKmph"].(string)

	// Get location info
	nearestArea := data["nearest_area"].([]interface{})[0].(map[string]interface{})
	areaName := nearestArea["areaName"].([]interface{})[0].(map[string]interface{})["value"].(string)
	country := nearestArea["country"].([]interface{})[0].(map[string]interface{})["value"].(string)

	result := fmt.Sprintf("🌤️  Weather for %s, %s\n", areaName, country)
	result += fmt.Sprintf("Condition: %s\n", weatherDesc)
	result += fmt.Sprintf("Temperature: %s°C (Feels like %s°C)\n", temp, feelsLike)
	result += fmt.Sprintf("Humidity: %s%%\n", humidity)
	result += fmt.Sprintf("Wind Speed: %s km/h", windSpeed)

	// Add forecast if requested
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

// main function for the plugin
func main() {
	plugin := NewWeatherPlugin()

	// Set up JSON communication
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	// Log to stderr
	log.SetOutput(os.Stderr)

	for {
		var req sdk.PluginRequest
		if err := decoder.Decode(&req); err != nil {
			log.Printf("Failed to decode request: %v", err)
			break
		}

		response := plugin.handleRequest(req)

		if err := encoder.Encode(response); err != nil {
			log.Printf("Failed to encode response: %v", err)
			break
		}
	}
}

// handleRequest handles incoming requests
func (p *WeatherPlugin) handleRequest(req sdk.PluginRequest) sdk.PluginResponse {
	switch req.Type {
	case "init":
		var initData map[string]interface{}
		if err := json.Unmarshal(req.Data, &initData); err != nil {
			return sdk.PluginResponse{
				Type:    "init",
				Success: false,
				Error:   fmt.Sprintf("failed to parse init data: %v", err),
			}
		}

		if configData, ok := initData["config"].(map[string]interface{}); ok {
			config := sdk.Config{
				PluginDir: configData["plugin_dir"].(string),
				DataDir:   configData["data_dir"].(string),
				Settings:  make(map[string]string),
			}
			if settings, ok := configData["settings"].(map[string]interface{}); ok {
				for k, v := range settings {
					if str, ok := v.(string); ok {
						config.Settings[k] = str
					}
				}
			}

			if err := p.Init(config); err != nil {
				return sdk.PluginResponse{
					Type:    "init",
					Success: false,
					Error:   fmt.Sprintf("failed to initialize plugin: %v", err),
				}
			}
		}

		return sdk.PluginResponse{
			Type:    "init",
			Success: true,
		}

	case "message":
		var msg sdk.Message
		if err := json.Unmarshal(req.Data, &msg); err != nil {
			return sdk.PluginResponse{
				Type:    "message",
				Success: false,
				Error:   fmt.Sprintf("failed to parse message: %v", err),
			}
		}

		responses, err := p.OnMessage(msg)
		if err != nil {
			return sdk.PluginResponse{
				Type:    "message",
				Success: false,
				Error:   fmt.Sprintf("failed to process message: %v", err),
			}
		}

		if len(responses) > 0 {
			responseData, _ := json.Marshal(responses[0])
			return sdk.PluginResponse{
				Type:    "message",
				Success: true,
				Data:    responseData,
			}
		}

		return sdk.PluginResponse{
			Type:    "message",
			Success: true,
		}

	case "command":
		var args []string
		if err := json.Unmarshal(req.Data, &args); err != nil {
			return sdk.PluginResponse{
				Type:    "command",
				Success: false,
				Error:   fmt.Sprintf("failed to parse command args: %v", err),
			}
		}

		location := "Charlotte,NC" // Default
		if len(args) > 0 {
			location = strings.Join(args, " ")
		}

		var weather string
		var err error

		switch req.Command {
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

		responseData, _ := json.Marshal(weatherMsg)
		return sdk.PluginResponse{
			Type:    "message",
			Success: true,
			Data:    responseData,
		}

	case "shutdown":
		return sdk.PluginResponse{
			Type:    "shutdown",
			Success: true,
		}

	default:
		return sdk.PluginResponse{
			Type:    req.Type,
			Success: false,
			Error:   "unknown request type",
		}
	}
}
