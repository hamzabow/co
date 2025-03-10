package genmessage

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
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
	ErrNoChangesInRepo   = errors.New("no changes detected in the repository")
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
		return "", ErrNoChangesInRepo
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
