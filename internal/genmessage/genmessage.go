package genmessage

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hamzabow/co/internal/prompts"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/joho/godotenv/autoload"
)

var (
	ErrMissingAPIKey     = errors.New("missing OPENAI_API_KEY environment variable")
	ErrFailedToGetDiffs  = errors.New("failed to get git diffs")
	ErrNoChangesInRepo   = errors.New("no staged changes detected in the repository; use 'git add' to stage changes")
	ErrOpenAIFetchFailed = errors.New("failed to fetch response from OpenAI API")
)

// Define a custom model that embeds spinner.Model and implements tea.Model
type customSpinnerModel struct {
	spinner.Model
	message string
}

// Implement the Init method for the custom model
func (m customSpinnerModel) Init() tea.Cmd {
	return m.Model.Tick
}

// Implement the Update method for the custom model
func (m customSpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := m.Model.Update(msg)
	m.Model = newModel
	return m, cmd
}

// Implement the View method for the custom model
func (m customSpinnerModel) View() string {
	// Apply the style to both the spinner and the message text
	return m.Model.View() + " " + m.message
}

func GenerateCommitMessage(key string) (string, error) {

	client := openai.NewClient(option.WithAPIKey(key))

	diff, err := getGitDiff()

	if err != nil {
		return "", ErrFailedToGetDiffs
	}

	if diff == "" {
		// Prompt the user if they want to stage all changes
		fmt.Print("No staged changes detected. Would you like to stage all changes? (Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response == "" || response == "y" || response == "yes" {
			// Stage all changes
			err = stageAllChanges()
			if err != nil {
				return "", fmt.Errorf("failed to stage changes: %w", err)
			}
			fmt.Println("All changes staged successfully.")

			// Get the diff again after staging
			diff, err = getGitDiff()
			if err != nil {
				return "", ErrFailedToGetDiffs
			}

			if diff == "" {
				return "", errors.New("still no changes to commit after staging all changes")
			}
		} else {
			return "", ErrNoChangesInRepo
		}
	}

	// prompt := fmt.Sprintf(prompts.GitmojiPrompt, diff)
	prompt := fmt.Sprintf(prompts.LongConventionalCommitsPrompt, diff)

	// Create a new custom spinner model
	s := customSpinnerModel{
		Model: spinner.New(),
		message: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true).
			PaddingLeft(2).
			PaddingRight(2).
			Render(" Generating Commit Message "),
	}
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		PaddingLeft(2).
		PaddingTop(1)

	// Start the spinner in a separate goroutine
	p := tea.NewProgram(s)
	go func() {
		if err := p.Start(); err != nil {
			fmt.Println("Error starting spinner:", err)
		}
	}()

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})

	// Ensure the spinner stops completely
	p.Quit()
	// Give a small pause to allow the spinner goroutine to clean up
	time.Sleep(100 * time.Millisecond)

	if err != nil {
		return "", ErrOpenAIFetchFailed
	}
	response := chatCompletion.Choices[0].Message.Content
	return response, nil
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// stageAllChanges stages all changes in the repository using git add .
func stageAllChanges() error {
	cmd := exec.Command("git", "add", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add command failed: %s (%w)", output, err)
	}
	return nil
}
