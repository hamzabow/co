package apikeyinput

import (
	"errors"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var ErrEmptyApiKey = errors.New("you have entered an empty API key multiple times. Please try again later")

func PromptApiKeyWithRetries() (string, error) {
	key := PromptApiKey()
	if key == "" {
		max_retries := 3
		retries := max_retries
		for key == "" && retries > 0 {
			fmt.Println("-------------------------------------------------")
			fmt.Println("Key is empty! Please enter a valid OpenAI API KEY")
			fmt.Println("-------------------------------------------------")
			key = PromptApiKey()
			retries -= 1
		}
		if retries <= 0 && key == "" {
			return "", ErrEmptyApiKey
		}
	}
	return key, nil
}

func PromptApiKey() string {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	text := m.(model).textInput.Value()

	return text
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	err       error
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "..."
	ti.Focus()
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	// ti.CharLimit = 156
	// ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlP:

			if m.textInput.EchoMode != textinput.EchoPassword {
				m.textInput.EchoMode = textinput.EchoPassword
				m.textInput.EchoCharacter = '•'
			} else {
				m.textInput.EchoMode = textinput.EchoNormal
			}

			// m.toggleEchoMode()
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Enter your OpenAI API key:\n\n%s\n\n%s",
		m.textInput.View(),
		"(Escape to quit, Ctrl+P to show/hide password)",
	) + "\n"
}

func (m model) toggleEchoMode() {
	if m.textInput.EchoMode != textinput.EchoPassword {
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.EchoCharacter = '•'
		return

	}
	m.textInput.EchoMode = textinput.EchoNormal

}
