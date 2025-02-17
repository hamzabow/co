package genmessage

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	_ "github.com/joho/godotenv/autoload"
)

func GenerateCommitMessage() {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		fmt.Println("Error: OPENAI_API_KEY env variable is not set")
		return
	}

	client := openai.NewClient(option.WithAPIKey(key))

	diff, err := getGitDiff()

	if err != nil {
		fmt.Println("Failed to get diffs:", err)
		return
	}

	if diff == "" {
		fmt.Println("No changes detected in the repository.")
		return
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
		fmt.Println("Failed to prompt OpenAI API:", err)
		return
	}
	println(chatCompletion.Choices[0].Message.Content)

}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.CombinedOutput()
	return string(output), err
}
