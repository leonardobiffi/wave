package config

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

func NewItemDelegate(keys *DelegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	help := []key.Binding{keys.play, keys.mute, keys.stop}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type DelegateKeyMap struct {
	play key.Binding
	mute key.Binding
	stop key.Binding
}

// Additional short help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d DelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.play,
		d.mute,
		d.stop,
	}
}

// Additional full help entries. This satisfies the help.KeyMap interface and
// is entirely optional.
func (d DelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.play,
			d.mute,
			d.stop,
		},
	}
}

func (d DelegateKeyMap) Play() key.Binding {
	return d.play
}

func (d DelegateKeyMap) Mute() key.Binding {
	return d.mute
}

func (d DelegateKeyMap) Stop() key.Binding {
	return d.stop
}

func NewDelegateKeyMap() *DelegateKeyMap {
	return &DelegateKeyMap{
		play: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "play"),
		),
		mute: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "mute"),
		),
		stop: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "stop"),
		),
	}
}
