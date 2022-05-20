package player

import (
	"bufio"
	"os"
	"os/user"
	"strings"
)

func LoadStations() (stations []RadioStation) {
	usr, _ := user.Current()
	dir := usr.HomeDir
	default_file := dir + "/.go-radio/stations"
	if Exists(default_file) {
		f, err := os.Open(default_file)
		check(err)
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.Trim(scanner.Text(), "\n\r")
			pair := strings.Split(line, ",")
			if len(pair) == 2 {
				stations = append(stations,
					RadioStation{strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1]), strings.TrimSpace(pair[2])})
			}
		}
		check(scanner.Err())
	} else {
		stations = append(stations, RadioStation{"KISS FM 92.5", "http://cloud2.cdnseguro.com:23538/listen.pls", "Rock"})
		stations = append(stations, RadioStation{"WBEZ 91.5", "http://stream.wbez.org/wbez128.mp3", "General"})
	}

	return
}
