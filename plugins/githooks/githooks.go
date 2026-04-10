package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Cod-e-Codes/marchat/plugin/sdk"
)

// GitHooksPlugin manages git repositories and provides git status updates
type GitHooksPlugin struct {
	*sdk.BasePlugin
	config      sdk.Config
	watchedRepo string
}

// NewGitHooksPlugin creates a new git hooks plugin
func NewGitHooksPlugin() *GitHooksPlugin {
	return &GitHooksPlugin{
		BasePlugin: sdk.NewBasePlugin("githooks"),
	}
}

// Init initializes the git hooks plugin
func (p *GitHooksPlugin) Init(config sdk.Config) error {
	p.config = config

	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git command not found: %w", err)
	}

	log.Printf("Git hooks plugin initialized")
	return nil
}

// OnMessage handles incoming messages
func (p *GitHooksPlugin) OnMessage(msg sdk.Message) ([]sdk.Message, error) {
	return nil, nil
}

// Commands returns the commands this plugin provides
func (p *GitHooksPlugin) Commands() []sdk.PluginCommand {
	return []sdk.PluginCommand{
		{
			Name:        "git-status",
			Description: "Show git status of current directory",
			Usage:       ":git-status [path]",
			AdminOnly:   false,
		},
		{
			Name:        "git-log",
			Description: "Show recent git commits",
			Usage:       ":git-log [n] [path]",
			AdminOnly:   false,
		},
		{
			Name:        "git-branch",
			Description: "Show current git branch and available branches",
			Usage:       ":git-branch [path]",
			AdminOnly:   false,
		},
		{
			Name:        "git-diff",
			Description: "Show git diff of uncommitted changes",
			Usage:       ":git-diff [path]",
			AdminOnly:   false,
		},
		{
			Name:        "git-watch",
			Description: "Watch a repository for changes (admin only)",
			Usage:       ":git-watch <path>",
			AdminOnly:   true,
		},
	}
}

// getGitStatus retrieves git status information
func (p *GitHooksPlugin) getGitStatus(repoPath string) (string, error) {
	if repoPath == "" {
		repoPath = "."
	}

	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return "", fmt.Errorf("not a git repository: %s", repoPath)
	}

	cmd := exec.Command("git", "-C", repoPath, "status", "--short", "--branch")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git status failed: %w\n%s", err, output)
	}

	result := "📁 Git Status\n"
	result += "─────────────\n"

	if len(output) == 0 {
		result += "✓ Working tree clean"
	} else {
		result += string(output)
	}

	return result, nil
}

// getGitLog retrieves recent git commits
func (p *GitHooksPlugin) getGitLog(repoPath string, count int) (string, error) {
	if repoPath == "" {
		repoPath = "."
	}
	if count <= 0 {
		count = 5
	}

	cmd := exec.Command("git", "-C", repoPath, "log",
		fmt.Sprintf("--max-count=%d", count),
		"--pretty=format:%h - %s (%an, %ar)")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w\n%s", err, output)
	}

	result := fmt.Sprintf("📝 Recent Commits (last %d)\n", count)
	result += "─────────────────────\n"
	result += string(output)

	return result, nil
}

// getGitBranch shows current branch and available branches
func (p *GitHooksPlugin) getGitBranch(repoPath string) (string, error) {
	if repoPath == "" {
		repoPath = "."
	}

	cmd := exec.Command("git", "-C", repoPath, "branch", "--show-current")
	currentBranch, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	cmd = exec.Command("git", "-C", repoPath, "branch", "-a")
	branches, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get branches: %w", err)
	}

	result := "🌿 Git Branches\n"
	result += "──────────────\n"
	result += fmt.Sprintf("Current: %s\n", strings.TrimSpace(string(currentBranch)))
	result += "\nAll branches:\n"
	result += string(branches)

	return result, nil
}

// getGitDiff shows uncommitted changes
func (p *GitHooksPlugin) getGitDiff(repoPath string) (string, error) {
	if repoPath == "" {
		repoPath = "."
	}

	cmd := exec.Command("git", "-C", repoPath, "diff", "--stat")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w", err)
	}

	if len(output) == 0 {
		return "No uncommitted changes", nil
	}

	result := "📊 Uncommitted Changes\n"
	result += "────────────────────\n"
	result += string(output)

	return result, nil
}

func main() {
	plugin := NewGitHooksPlugin()
	if err := sdk.RunStdio(plugin, plugin.handleCommand); err != nil {
		log.Fatalf("plugin exited: %v", err)
	}
}

func (p *GitHooksPlugin) handleCommand(command string, args []string) sdk.PluginResponse {
	var result string
	var err error
	repoPath := ""

	switch command {
	case "git-status":
		if len(args) > 0 {
			repoPath = args[0]
		}
		result, err = p.getGitStatus(repoPath)

	case "git-log":
		count := 5
		if len(args) > 0 {
			_, _ = fmt.Sscanf(args[0], "%d", &count)
		}
		if len(args) > 1 {
			repoPath = args[1]
		}
		result, err = p.getGitLog(repoPath, count)

	case "git-branch":
		if len(args) > 0 {
			repoPath = args[0]
		}
		result, err = p.getGitBranch(repoPath)

	case "git-diff":
		if len(args) > 0 {
			repoPath = args[0]
		}
		result, err = p.getGitDiff(repoPath)

	case "git-watch":
		if len(args) == 0 {
			return sdk.PluginResponse{
				Type:    "command",
				Success: false,
				Error:   "repository path required",
			}
		}
		p.watchedRepo = args[0]
		result = fmt.Sprintf("Now watching repository: %s", p.watchedRepo)

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
			Error:   fmt.Sprintf("command failed: %v", err),
		}
	}

	msg := sdk.Message{
		Sender:    "GitBot",
		Content:   result,
		CreatedAt: time.Now(),
	}

	responseData, err := json.Marshal(msg)
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
