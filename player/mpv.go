package player

import (
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/DexterLB/mpvipc"
	"github.com/shirou/gopsutil/v3/process"
)

type MPV struct {
	PlayerName string
	IsPlaying  bool
	StreamURL  string
	Command    *exec.Cmd
	In         io.WriteCloser
	Out        io.ReadCloser
	PipeChan   chan io.ReadCloser
	Connection *mpvipc.Connection
}

func (player *MPV) Play(stream_url string) {
	if !player.IsPlaying {
		var err error
		is_playlist := strings.HasSuffix(stream_url, ".m3u") || strings.HasSuffix(stream_url, ".pls")
		if is_playlist {
			player.Command = exec.Command(player.PlayerName, "-quiet", "-playlist", stream_url, "--input-ipc-server=/tmp/mpv.sock")
		} else {
			player.Command = exec.Command(player.PlayerName, "-quiet", stream_url, "--input-ipc-server=/tmp/mpv.sock")
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

		var wait time.Duration = 10

		for {
			if wait == 0 {
				log.Fatal("Could not connect to mpv socket")
			}

			time.Sleep(time.Second * 1)
			conn := mpvipc.NewConnection("/tmp/mpv.sock")
			err = conn.Open()
			if err == nil {
				player.Connection = conn
				break
			}

			wait--
		}
	}
}

func (player *MPV) Close() {
	if player.IsPlaying {
		player.IsPlaying = false

		pid, err := player.Connection.Get("pid")
		if err != nil {
			log.Fatal(err)
		}
		val := int64(pid.(float64))

		p, err := process.NewProcess(int32(val))
		if err != nil {
			log.Fatal(err)
		}
		err = p.Kill()
		if err != nil {
			log.Fatal(err)
		}

		player.In.Close()
		player.Out.Close()
		player.Command = nil

		player.StreamURL = ""

		player.Connection.Close()
	}
}

func (player *MPV) Mute() {
	if player.IsPlaying {
		err := player.Connection.Set("mute", true)
		if err != nil {
			log.Fatal(err)
		}

		player.IsPlaying = false
		return
	}

	err := player.Connection.Set("mute", false)
	if err != nil {
		log.Fatal(err)
	}
	player.IsPlaying = true
}

func (player *MPV) Pause() {
	if player.IsPlaying {
		err := player.Connection.Set("pause", true)
		if err != nil {
			log.Fatal(err)
		}

		player.IsPlaying = false
		return
	}

	err := player.Connection.Set("pause", false)
	if err != nil {
		log.Fatal(err)
	}
	player.IsPlaying = true
}

func (player *MPV) IncVolume() {
	if player.IsPlaying {
		_, err := player.Connection.Call("cycle", "volume", "up")
		if err != nil {
			log.Fatal(err)
		}

		vol, err := player.Connection.Get("volume")
		if err != nil {
			log.Fatal(err)
		}

		log.Println(vol)
	}
}

func (player *MPV) DecVolume() {
	if player.IsPlaying {
		_, err := player.Connection.Call("cycle", "volume", "down")
		if err != nil {
			log.Fatal(err)
		}

		vol, err := player.Connection.Get("volume")
		if err != nil {
			log.Fatal(err)
		}

		log.Println(vol)
	}
}
