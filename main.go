package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/leonardobiffi/wave/player"
	"github.com/leonardobiffi/wave/version"
)

var docStyle = lipgloss.NewStyle().Margin(1, 3)

type item struct {
	title     string
	streamURL string
	desc      string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list  list.Model
	dj    player.Dj
	muted bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.dj.Stop()
			return m, tea.Quit
		case "x":
			m.dj.Stop()
		case "+":
			m.dj.VolumeUp()
		case "-":
			m.dj.VolumeDown()
		case "m":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if m.muted {
					m.list.NewStatusMessage(m.dj.FormatPlayStatus(i.title))
					m.dj.Mute()
					m.muted = false

				} else {
					m.list.NewStatusMessage(m.dj.FormatMuteStatus(i.title))
					m.dj.Mute()
					m.muted = true
				}
			}
		}

		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.list.NewStatusMessage(m.dj.FormatPlayStatus(i.title))
				m.dj.Play(m.list.Index())
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	var stations = player.LoadStationsApi()
	if len(stations) == 0 {
		panic("no stations found")
	}

	var items []list.Item
	for _, station := range stations {
		items = append(items, item{station.Name, station.StreamURL, station.Subtitle})
	}

	var pipeChan = make(chan io.ReadCloser)
	var mpv = player.MPV{PlayerName: "mpv", IsPlaying: false, PipeChan: pipeChan}

	var list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.SetShowStatusBar(false)

	m := model{
		list:  list,
		dj:    player.Dj{Player: &mpv, Stations: stations, CurrentStation: -1},
		muted: false,
	}

	m.list.Title = fmt.Sprintf("Wave - Radio Player v%s", version.String())
	m.list.NewStatusMessage("Press enter to play")

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
