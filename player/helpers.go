package player

import (
	"bufio"
	"os"
	"os/user"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func LoadStations() []RadioStation {
	usr, _ := user.Current()
	dir := usr.HomeDir
	defaultFile := dir + "/.go-radio/stations"
	if Exists(defaultFile) {
		f, err := os.Open(defaultFile)
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
	}

	return stations
}
