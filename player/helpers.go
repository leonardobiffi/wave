package player

import (
	"os"
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
