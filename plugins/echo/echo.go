package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Cod-e-Codes/marchat/plugin/sdk"
)

// EchoPlugin is a simple echo plugin
type EchoPlugin struct {
	*sdk.BasePlugin
	config sdk.Config
}

// NewEchoPlugin creates a new echo plugin
func NewEchoPlugin() *EchoPlugin {
	return &EchoPlugin{
		BasePlugin: sdk.NewBasePlugin("echo"),
	}
}

// Init initializes the echo plugin
func (p *EchoPlugin) Init(config sdk.Config) error {
	p.config = config
	return nil
}

// OnMessage handles incoming messages
func (p *EchoPlugin) OnMessage(msg sdk.Message) ([]sdk.Message, error) {
	// Echo messages that start with "echo:"
	if len(msg.Content) > 5 && msg.Content[:5] == "echo:" {
		echoMsg := sdk.Message{
			Sender:    "EchoBot",
			Content:   msg.Content[5:], // Remove "echo:" prefix
			CreatedAt: time.Now(),
		}
		return []sdk.Message{echoMsg}, nil
	}
	return nil, nil
}

// Commands returns the commands this plugin provides
func (p *EchoPlugin) Commands() []sdk.PluginCommand {
	return []sdk.PluginCommand{
		{
			Name:        "echo",
			Description: "Echo a message",
			Usage:       ":echo <message>",
			AdminOnly:   false,
		},
		{
			Name:        "echo-admin",
			Description: "Echo a message (admin only)",
			Usage:       ":echo-admin <message>",
			AdminOnly:   true,
		},
	}
}

// main function for the plugin
func main() {
	plugin := NewEchoPlugin()

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
func (p *EchoPlugin) handleRequest(req sdk.PluginRequest) sdk.PluginResponse {
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

		// Extract config
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

		// Handle echo command
		if req.Command == "echo" && len(args) > 0 {
			// Handle quoted arguments
			var content string
			if len(args) == 1 && strings.HasPrefix(args[0], `"`) && strings.HasSuffix(args[0], `"`) {
				// Single quoted argument
				content = strings.Trim(args[0], `"`)
			} else {
				// Join all arguments into a single message
				content = strings.Join(args, " ")
			}

			log.Printf("Echoing content: '%s'", content)

			echoMsg := sdk.Message{
				Sender:    "EchoBot",
				Content:   content,
				CreatedAt: time.Now(),
			}

			responseData, _ := json.Marshal(echoMsg)
			return sdk.PluginResponse{
				Type:    "message",
				Success: true,
				Data:    responseData,
			}
		}

		// Handle echo-admin command
		if req.Command == "echo-admin" && len(args) > 0 {
			// Handle quoted arguments
			var content string
			if len(args) == 1 && strings.HasPrefix(args[0], `"`) && strings.HasSuffix(args[0], `"`) {
				// Single quoted argument
				content = strings.Trim(args[0], `"`)
			} else {
				// Join all arguments into a single message
				content = strings.Join(args, " ")
			}

			log.Printf("Echoing admin content: '%s'", content)

			echoMsg := sdk.Message{
				Sender:    "EchoBot",
				Content:   content,
				CreatedAt: time.Now(),
			}

			responseData, _ := json.Marshal(echoMsg)
			return sdk.PluginResponse{
				Type:    "message",
				Success: true,
				Data:    responseData,
			}
		}

		return sdk.PluginResponse{
			Type:    "command",
			Success: false,
			Error:   "unknown command",
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
