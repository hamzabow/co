package messagetextarea

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// func MessageTextArea(msg string) string {
func MessageTextArea(msg string) string {
	p := tea.NewProgram(initialModel(msg))

	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	text := m.(model).textarea.Value()

	return text
}

type errMsg error

type model struct {
	textarea textarea.Model
	err      error
}

func initialModel(initialValue string) model {
	ti := textarea.New()
	ti.SetValue(initialValue)
	ti.Placeholder = "Commit message ..."
	ti.Focus()
	ti.SetWidth(80)
	ti.SetHeight(10)
	ti.ShowLineNumbers = false

	ti.Prompt = " "

	textColor := lipgloss.Color("212")
	ti.FocusedStyle.Base = ti.FocusedStyle.Base.Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("69"))
	ti.FocusedStyle.Text = ti.FocusedStyle.Text.Foreground(textColor)
	ti.FocusedStyle.CursorLine = ti.FocusedStyle.CursorLine.Foreground(textColor)
	ti.BlurredStyle.Text = ti.BlurredStyle.Text.Foreground(lipgloss.Color("240"))

	return model{
		textarea: ti,
		err:      nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			// Don't do anything special for regular Enter
			// Let it pass to the textarea for normal processing
		case tea.KeyCtrlJ: // Ctrl+Enter is often mapped to Ctrl+J in terminals
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		"Here is the commit message:\n\n%s\n\n%s",
		m.textarea.View(),
		"(ctrl+c to quit)",
	) + "\n\n"
}
