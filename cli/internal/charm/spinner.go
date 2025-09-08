package charm

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vegidio/mediaorient"
)

type orientationMsg struct {
	media mediaorient.Media
}

type orientationDone struct{}

func orientationDirCmd(ch <-chan mediaorient.Result[mediaorient.Media]) tea.Cmd {
	return func() tea.Msg {
		if resp, ok := <-ch; ok {
			return orientationMsg{resp.Data}
		}

		return orientationDone{}
	}
}

type spinnerDModel struct {
	spinner spinner.Model
	result  <-chan mediaorient.Result[mediaorient.Media]
	message string
	media   []mediaorient.Media
}

func initSpinnerDModel(result <-chan mediaorient.Result[mediaorient.Media], message string) *spinnerDModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = pink

	return &spinnerDModel{
		spinner: s,
		result:  result,
		message: message,
		media:   make([]mediaorient.Media, 0),
	}
}

func (m *spinnerDModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		orientationDirCmd(m.result),
	)
}

func (m *spinnerDModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msgValue := msg.(type) {
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case orientationMsg:
		m.media = append(m.media, msgValue.media)
		return m, orientationDirCmd(m.result)

	case orientationDone:
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
	return fmt.Sprintf("%s %s%s Media\n", m.message, m.spinner.View(), bold.Render(strconv.Itoa(len(m.media))))
}

func StartSpinner(result <-chan mediaorient.Result[mediaorient.Media], message string) ([]mediaorient.Media, error) {
	model, err := tea.NewProgram(initSpinnerDModel(result, message)).Run()
	if err != nil {
		return nil, err
	}

	m := model.(*spinnerDModel)
	return m.media, nil
}
