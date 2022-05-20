package player

import (
	"io"
	"os/exec"
	"strings"
)

type MPlayer struct {
	PlayerName string
	IsPlaying  bool
	StreamURL  string
	Command    *exec.Cmd
	In         io.WriteCloser
	Out        io.ReadCloser
	PipeChan   chan io.ReadCloser
}

func (player *MPlayer) Play(stream_url string) {
	if !player.IsPlaying {
		var err error
		is_playlist := strings.HasSuffix(stream_url, ".m3u") || strings.HasSuffix(stream_url, ".pls")
		if is_playlist {
			player.Command = exec.Command(player.PlayerName, "-quiet", "-playlist", stream_url)
		} else {
			player.Command = exec.Command(player.PlayerName, "-quiet", stream_url)
		}
		player.In, err = player.Command.StdinPipe()
		check(err)
		player.Out, err = player.Command.StdoutPipe()
		check(err)

		err = player.Command.Start()
		check(err)

		player.IsPlaying = true
		player.StreamURL = stream_url
		go func() {
			player.PipeChan <- player.Out
		}()
	}
}

func (player *MPlayer) Close() {
	if player.IsPlaying {
		player.IsPlaying = false

		player.In.Write([]byte("q"))
		player.In.Close()
		player.Out.Close()
		player.Command = nil

		player.StreamURL = ""
	}
}

func (player *MPlayer) Mute() {
	if player.IsPlaying {
		player.In.Write([]byte("m"))
	}
}

func (player *MPlayer) Pause() {
	if player.IsPlaying {
		player.In.Write([]byte("p"))
	}
}

func (player *MPlayer) IncVolume() {
	if player.IsPlaying {
		player.In.Write([]byte("*"))
	}
}

func (player *MPlayer) DecVolume() {
	if player.IsPlaying {
		player.In.Write([]byte("/"))
	}
}
