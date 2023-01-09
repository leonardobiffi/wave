package player

import (
	"os"
	"os/user"

	"github.com/leonardobiffi/wave/radiogarden"
	"gopkg.in/yaml.v3"
)

var defaultStations = []RadioStation{
	{
		Name:      "Super Rádio Tupi 96.5 FM",
		StreamURL: "https://radio.garden/api/ara/content/listen/RMuzTKOn/channel.mp3",
		Subtitle:  "Rio de Janeiro RJ, Brazil",
	},
}

func LoadStationsApi() []RadioStation {
	var stations []RadioStation

	client := radiogarden.New()

	personalStations := loadPersonalStations()
	for _, station := range personalStations.Stations.Search {

		result, err := client.SearchStations(station)
		if err != nil {
			panic(err)
		}

		// get first result, highest score
		if len(result) > 0 {
			station := result[0]

			url := client.GetStationStream(radiogarden.ExtractID(station.Url))
			stations = append(stations, RadioStation{
				Name:      station.Title,
				StreamURL: url,
				Subtitle:  station.Subtitle,
			})
		}
	}

	stations = append(stations, defaultStations...)
	return stations
}

type PersonalRadioStation struct {
	Stations struct {
		Search []string `yaml:"search"`
	} `yaml:"stations"`
}

func loadPersonalStations() (stations PersonalRadioStation) {

	usr, _ := user.Current()
	dir := usr.HomeDir
	defaultFile := dir + "/.wave/stations.yaml"
	if Exists(defaultFile) {
		f, err := os.ReadFile(defaultFile)
		check(err)

		err = yaml.Unmarshal(f, &stations)
		check(err)
	}

	return stations
}
