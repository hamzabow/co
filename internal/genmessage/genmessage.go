package genmessage

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	spinner "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
}

// Implement the Init method for the custom model
func (m customSpinnerModel) Init() tea.Cmd {
	return m.Model.Tick
}

// Implement the Update method for the custom model
func (m customSpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := m.Model.Update(msg)
	return customSpinnerModel{newModel}, cmd
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

	prompt := fmt.Sprintf(
		`
Here are the git diffs:

%s

Generate a concise and clear commit message describing these changes.`, diff)

	// Create a new custom spinner model
	s := customSpinnerModel{spinner.New()}
	s.Spinner = spinner.Dot

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
