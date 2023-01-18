package radiogarden

import "strings"

func ExtractID(url string) string {
	u := strings.Split(url, "/")
	return u[len(u)-1]
}
