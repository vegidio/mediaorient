package charm

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vegidio/mediaorient"
)

func orientationDirCmd(directory string, mediaType string, recursive bool) tea.Cmd {
	return func() tea.Msg {
		result, err := mediaorient.CalculateDirectoryOrientation(directory, mediaType, recursive)
		return spinnerDoneMsg{result, err}
	}
}

type spinnerDModel struct {
	spinner   spinner.Model
	text      string
	directory string
	mediaType string
	recursive bool
	result    []mediaorient.Media
	err       error
}

func initSpinnerModel(directory string, mediaType string, recursive bool) *spinnerDModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = pink

	return &spinnerDModel{
		spinner:   s,
		text:      fmt.Sprintf("⏳ Calculating the orientation in the directory %s...", green.Render(directory)),
		directory: directory,
		mediaType: mediaType,
		recursive: recursive,
	}
}

func (m *spinnerDModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		orientationDirCmd(m.directory, m.mediaType, m.recursive),
	)
}

func (m *spinnerDModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgValue := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case spinnerDoneMsg:
		m.result = msgValue.result
		m.err = msgValue.err
		return m, tea.Quit

	case tea.KeyMsg:
		switch msgValue.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m *spinnerDModel) View() string {
	return fmt.Sprintf("%s %s\n", m.text, m.spinner.View())
}

func SpinnerDir(directory string, mediaType string, recursive bool) ([]mediaorient.Media, error) {
	model, _ := tea.NewProgram(initSpinnerModel(directory, mediaType, recursive)).Run()
	m := model.(*spinnerDModel)
	return m.result, m.err
}
