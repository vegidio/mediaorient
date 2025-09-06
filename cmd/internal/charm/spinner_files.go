package charm

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vegidio/mediaorient"
)

func orientationFilesCmd(files []string) tea.Cmd {
	return func() tea.Msg {
		result, err := mediaorient.CalculateFilesOrientation(files)
		return spinnerDoneMsg{result, err}
	}
}

type spinnerFModel struct {
	spinner spinner.Model
	text    string
	files   []string
	result  []mediaorient.Media
	err     error
}

func initSpinnerFModel(files []string) *spinnerFModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = pink

	return &spinnerFModel{
		spinner: s,
		text:    fmt.Sprintf("⏳ Calculating the orientation in %s files...", green.Render(strconv.Itoa(len(files)))),
		files:   files,
	}
}

func (m *spinnerFModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		orientationFilesCmd(m.files),
	)
}

func (m *spinnerFModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *spinnerFModel) View() string {
	return fmt.Sprintf("\n%s %s\n", m.text, m.spinner.View())
}

func SpinnerFiles(files []string) ([]mediaorient.Media, error) {
	model, _ := tea.NewProgram(initSpinnerFModel(files)).Run()
	m := model.(*spinnerFModel)
	return m.result, m.err
}
