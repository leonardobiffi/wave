package player

import (
	"os"
	"os/user"
	"sort"

	"github.com/leonardobiffi/wave/pkg/radiogarden"
	"gopkg.in/yaml.v3"
)

var defaultStations = []string{
	"Super Rádio Tupi, Rio de Janeiro RJ, Brazil",
	"Kiss FM, São Paulo SP, Brazil",
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
	} else {
		// create the file yaml file
		stations.Stations.Search = defaultStations
		d, err := yaml.Marshal(&stations)
		check(err)

		err = os.WriteFile(defaultFile, d, 0644)
		check(err)
	}

	sort.Strings(stations.Stations.Search)
	return stations
}
