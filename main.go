package main

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/leonardobiffi/wave/pkg/config"
	"github.com/leonardobiffi/wave/pkg/helpers"
	"github.com/leonardobiffi/wave/pkg/player"
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
	list         list.Model
	dj           player.Dj
	muted        bool
	keys         *listKeyMap
	delegateKeys *config.DelegateKeyMap
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case helpers.Contains(msg.String(), []string{"q", "esc", "ctrl+c"}):
			m.dj.Stop()
			return m, tea.Quit
		case key.Matches(msg, m.delegateKeys.Stop()):
			m.dj.Stop()
			m.list.NewStatusMessage("Press enter to play")
		case key.Matches(msg, m.keys.volupeUp):
			m.dj.VolumeUp()
		case key.Matches(msg, m.keys.volupeDown):
			m.dj.VolumeDown()
		case key.Matches(msg, m.delegateKeys.Mute()):
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

		if key.Matches(msg, m.delegateKeys.Play()) {
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

type listKeyMap struct {
	volupeUp   key.Binding
	volupeDown key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		volupeUp: key.NewBinding(
			key.WithKeys("+"),
			key.WithHelp("+", "increase volume"),
		),
		volupeDown: key.NewBinding(
			key.WithKeys("-"),
			key.WithHelp("-", "decresase volume"),
		),
	}
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

	var listKeys = newListKeyMap()
	var delegateKeys = config.NewDelegateKeyMap()
	var pipeChan = make(chan io.ReadCloser)
	var mpv = player.MPV{PlayerName: "mpv", IsPlaying: false, PipeChan: pipeChan}

	delegate := config.NewItemDelegate(delegateKeys)
	var list = list.New(items, delegate, 0, 0)
	list.SetShowStatusBar(false)

	m := model{
		list:         list,
		dj:           player.Dj{Player: &mpv, Stations: stations, CurrentStation: -1},
		muted:        false,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}

	m.list.Title = fmt.Sprintf("Wave - Radio Player v%s", version.String())
	m.list.NewStatusMessage("Press enter to play")
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.volupeUp,
			listKeys.volupeDown,
		}
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
