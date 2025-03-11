package confirmation

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Result represents the outcome of the confirmation
type Result int

const (
	// Undecided means the user hasn't made a choice yet
	Undecided Result = iota
	// Confirmed means the user accepted/confirmed
	Confirmed
	// Cancelled means the user declined/cancelled
	Cancelled
)

// ConfirmationModel represents a confirmation dialog model
type Model struct {
	question     string
	cursor       int
	choices      []string
	choiceStyles []lipgloss.Style
	result       Result
	width        int
	height       int
}

var (
	// Title style with white text on purple background
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2).
			MarginBottom(1)

	// Container style for the entire view
	containerStyle = lipgloss.NewStyle().
			MarginLeft(1)

	// Question style
	questionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			MarginBottom(1)

	// Active choice style
	activeChoiceStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Bold(true).
				PaddingLeft(1).
				PaddingRight(1)

	// Inactive choice style
	inactiveChoiceStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#AAAAAA")).
				PaddingLeft(1).
				PaddingRight(1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

// New creates a new confirmation dialog model
func New(question string, defaultYes bool) Model {
	cursor := 0
	if defaultYes {
		cursor = 0 // "Yes" is selected by default
	} else {
		cursor = 1 // "No" is selected by default
	}

	// Set up initial choice styles
	choices := []string{"Yes", "No"}
	choiceStyles := make([]lipgloss.Style, len(choices))
	for i := range choiceStyles {
		if i == cursor {
			choiceStyles[i] = activeChoiceStyle
		} else {
			choiceStyles[i] = inactiveChoiceStyle
		}
	}

	return Model{
		question:     question,
		cursor:       cursor,
		choices:      choices,
		choiceStyles: choiceStyles,
		result:       Undecided,
		width:        80, // Default value, will be updated
		height:       24, // Default value, will be updated
	}
}

// Confirm runs the confirmation dialog and returns the result
func Confirm(question string, defaultYes bool) (bool, error) {
	m := New(question, defaultYes)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	m = finalModel.(Model)
	return m.result == Confirmed, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: 80, Height: 24}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.result = Cancelled
			return m, tea.Quit

		case "enter":
			if m.cursor == 0 {
				m.result = Confirmed
			} else {
				m.result = Cancelled
			}
			return m, tea.Quit

		case "left", "h":
			m.cursor = 0 // Select "Yes"
			m.updateChoiceStyles()
			return m, nil

		case "right", "l":
			m.cursor = 1 // Select "No"
			m.updateChoiceStyles()
			return m, nil

		case "y", "Y":
			m.cursor = 0 // Select "Yes"
			m.result = Confirmed
			return m, tea.Quit

		case "n", "N":
			m.cursor = 1 // Select "No"
			m.result = Cancelled
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the model
func (m Model) View() string {
	var view strings.Builder

	view.WriteString(titleStyle.Render(" Confirmation "))
	view.WriteString("\n")

	view.WriteString(questionStyle.Render(m.question))
	view.WriteString("\n")

	// Render choices
	for i, choice := range m.choices {
		view.WriteString(m.choiceStyles[i].Render(choice))
		view.WriteString(" ")
	}
	view.WriteString("\n\n")

	view.WriteString(helpStyle.Render("Use arrow keys to select, Enter to confirm, Esc to cancel"))

	return containerStyle.Render(view.String())
}

// updateChoiceStyles updates the styles of choices based on the cursor position
func (m *Model) updateChoiceStyles() {
	for i := range m.choiceStyles {
		if i == m.cursor {
			m.choiceStyles[i] = activeChoiceStyle
		} else {
			m.choiceStyles[i] = inactiveChoiceStyle
		}
	}
}
