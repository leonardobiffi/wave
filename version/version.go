package version

var Version = "0.0.0-dev"
var template = `Wave {{printf "%s" .Version}}`

func String() string {
	return Version
}

func Template() string {
	return template
}
