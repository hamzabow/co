package messagetextarea

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CommitResult int

const (
	ResultCancel CommitResult = iota
	ResultCommit
)

// func MessageTextArea(msg string) string {
func MessageTextArea(msg string) (string, CommitResult) {
	p := tea.NewProgram(initialModel(msg))

	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	finalModel := m.(model)
	text := finalModel.textarea.Value()

	return text, finalModel.result
}

type errMsg error

type model struct {
	textarea textarea.Model
	err      error
	result   CommitResult
	width    int
	height   int
}

var (
	// Title style with white text on purple background
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2).
			MarginLeft(1).
			MarginBottom(0)

	// Container style for the entire view
	containerStyle = lipgloss.NewStyle().
			MarginLeft(1)

	// Input box style
	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

func initialModel(initialValue string) model {
	ti := textarea.New()
	ti.SetValue(initialValue)
	ti.Placeholder = "Commit message ..."
	ti.Focus()
	ti.ShowLineNumbers = false

	ti.Prompt = " "

	// Match the color scheme with the API key input
	ti.FocusedStyle.Base = lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("#7D56F4"))
	ti.FocusedStyle.Text = ti.FocusedStyle.Text.Foreground(lipgloss.Color("#FAFAFA"))
	ti.FocusedStyle.CursorLine = ti.FocusedStyle.CursorLine.Foreground(lipgloss.Color("#FAFAFA"))
	ti.BlurredStyle.Text = ti.BlurredStyle.Text.Foreground(lipgloss.Color("#626262"))

	return model{
		textarea: ti,
		err:      nil,
		result:   ResultCancel,
		// We'll set proper width and height when we get a WindowSizeMsg
		width:  80, // Default value, will be updated
		height: 24, // Default value, will be updated
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		func() tea.Msg {
			return tea.WindowSizeMsg{Width: 80, Height: 24}
		},
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Store the window size
		m.width = msg.Width
		m.height = msg.Height

		// Responsive behavior: automatically adjust dimensions based on terminal size
		// Calculate dynamic width for the textarea
		// Leave some margin on both sides and account for borders
		textareaWidth := m.width - 6 // 2 for container margin + 4 for borders/padding
		if textareaWidth < 20 {      // Minimum reasonable width
			textareaWidth = 20
		}

		// Update the textarea width
		m.textarea.SetWidth(textareaWidth)

		// Adjust height if needed, leaving space for title and help text
		m.textarea.SetHeight(10) // Or adjust dynamically if needed

		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			m.result = ResultCancel
			return m, tea.Quit
		case tea.KeyCtrlY: // Use Ctrl+Y as the commit shortcut
			m.result = ResultCommit
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
	var view strings.Builder

	view.WriteString(titleStyle.Render(" Commit Message "))
	view.WriteString("\n\n")

	// Dynamically set the width of the input box style based on terminal width
	dynamicInputBoxStyle := inputBoxStyle.Copy().Width(m.width - 4) // Account for container margin

	// Wrap the textarea in the dynamicInputBoxStyle
	view.WriteString(dynamicInputBoxStyle.Render(m.textarea.View()))
	view.WriteString("\n")

	view.WriteString(helpStyle.Render("  Ctrl+C to quit, Ctrl+Y to commit"))

	// Apply container style to the entire view
	return containerStyle.Render(view.String())
}
