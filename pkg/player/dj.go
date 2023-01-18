package player

type Dj struct {
	Player         RadioPlayer
	Stations       []RadioStation
	CurrentStation int
}

type RadioPlayer interface {
	Play(streamURL string)
	Mute()
	Pause()
	IncVolume()
	DecVolume()
	Close()
}

type RadioStation struct {
	Name      string
	StreamURL string
	Subtitle  string
}

func (dj *Dj) Play(station int) {
	if 0 <= station && station < len(dj.Stations) && dj.CurrentStation != station {
		if dj.CurrentStation >= 0 {
			dj.Player.Close()
		}

		dj.CurrentStation = station
		dj.Player.Play(dj.Stations[dj.CurrentStation].StreamURL)
	}
}

func (dj *Dj) Stop() {
	if dj.CurrentStation >= 0 {
		dj.Player.Close()
		dj.CurrentStation = -1
	}
}

func (dj *Dj) Mute() {
	if dj.CurrentStation >= 0 {
		dj.Player.Mute()
	}
}

func (dj *Dj) VolumeUp() {
	if dj.CurrentStation >= 0 {
		dj.Player.IncVolume()
	}
}

func (dj *Dj) VolumeDown() {
	if dj.CurrentStation >= 0 {
		dj.Player.DecVolume()
	}
}

func (dj *Dj) FormatPlayStatus(station string) string {
	return station + ": Playing... 🔊"
}

func (dj *Dj) FormatMuteStatus(station string) string {
	return station + ": Muted 🙊"
}
