package genmessage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	_ "github.com/joho/godotenv/autoload"
)

var (
	ErrMissingAPIKey     = errors.New("missing OPENAI_API_KEY environment variable")
	ErrFailedToGetDiffs  = errors.New("failed to get git diffs")
	ErrNoChangesInRepo   = errors.New("No changes detected in the repository")
	ErrOpenAIFetchFailed = errors.New("failed to fetch response from OpenAI API")
)

func GenerateCommitMessage() (string, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return "", ErrMissingAPIKey
	}

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
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
	})
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
