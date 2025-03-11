package apikeyinput

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	ErrEmptyApiKey = errors.New("you have entered an empty API key multiple times. Please try again later")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2).
			MarginBottom(0)

	// Add a new container style for the entire view
	containerStyle = lipgloss.NewStyle().
			MarginTop(1).
			MarginLeft(1)

	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1).
			Width(50)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	attemptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			MarginTop(1)
)

func PromptApiKeyWithRetries() (string, error) {
	m := initialModel()
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	m = finalModel.(model)

	if m.userQuit {
		return "", errors.New("user cancelled the operation")
	}

	if m.attempts >= m.maxAttempts && m.textInput.Value() == "" {
		return "", ErrEmptyApiKey
	}

	return m.textInput.Value(), nil
}

func PromptApiKey() (string, bool) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	finalModel := m.(model)
	text := finalModel.textInput.Value()

	return text, finalModel.userQuit
}

type errMsg error

type model struct {
	textInput   textinput.Model
	err         error
	userQuit    bool
	attempts    int
	maxAttempts int
	showError   bool
	width       int
	height      int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "sk-..."
	ti.Focus()
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	return model{
		textInput:   ti,
		err:         nil,
		userQuit:    false,
		attempts:    0,
		maxAttempts: 3,
		showError:   false,
		width:       80,
		height:      24,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			return tea.WindowSizeMsg{Width: 80, Height: 24}
		},
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		inputWidth := m.width - 6
		if inputWidth < 20 {
			inputWidth = 20
		}

		m.textInput.Width = inputWidth - 2

		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.textInput.Value() == "" {
				m.attempts++
				m.showError = true

				if m.attempts >= m.maxAttempts {
					return m, tea.Quit
				}

				return m, nil
			}
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			m.userQuit = true
			return m, tea.Quit

		case tea.KeyCtrlP:
			if m.textInput.EchoMode != textinput.EchoPassword {
				m.textInput.EchoMode = textinput.EchoPassword
				m.textInput.EchoCharacter = '•'
			} else {
				m.textInput.EchoMode = textinput.EchoNormal
			}

		default:
			if m.showError && len(msg.String()) > 0 {
				m.showError = false
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var view strings.Builder

	view.WriteString(titleStyle.Render(" OpenAI API Key "))
	view.WriteString("\n\n")

	dynamicInputBoxStyle := inputBoxStyle.Copy().Width(m.width - 4)

	view.WriteString(dynamicInputBoxStyle.Render(m.textInput.View()))
	view.WriteString("\n")

	if m.showError {
		remainingAttempts := m.maxAttempts - m.attempts
		errorMessage := "API key cannot be empty!"

		if remainingAttempts > 0 {
			errorMessage += fmt.Sprintf(" (%d attempts remaining)", remainingAttempts)
		} else {
			errorMessage += " (last attempt)"
		}

		view.WriteString(errorStyle.Render(errorMessage))
		view.WriteString("\n")
	}

	if m.attempts > 0 && !m.showError {
		view.WriteString(attemptStyle.Render(fmt.Sprintf("Attempt %d of %d", m.attempts+1, m.maxAttempts)))
		view.WriteString("\n")
	}

	view.WriteString(helpStyle.Render("Press Enter to submit, Esc to quit, Ctrl+P to toggle visibility"))

	return containerStyle.Render(view.String())
}

func (m model) toggleEchoMode() {
	if m.textInput.EchoMode != textinput.EchoPassword {
		m.textInput.EchoMode = textinput.EchoPassword
		m.textInput.EchoCharacter = '•'
		return

	}
	m.textInput.EchoMode = textinput.EchoNormal

}
