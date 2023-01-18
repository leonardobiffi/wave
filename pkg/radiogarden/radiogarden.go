package radiogarden

import (
	"context"
	"fmt"

	radiogarden "github.com/jonasrmichel/radio-garden-go"
)

const RadioGardenAPIUrl = "https://radio.garden/api"

type RadioGarden struct {
	client *radiogarden.ClientWithResponses
}

func New() *RadioGarden {
	client, err := radiogarden.NewClientWithResponses(RadioGardenAPIUrl)
	if err != nil {
		panic(err)
	}

	return &RadioGarden{
		client: client,
	}
}

func (r *RadioGarden) GetStationStream(stationId string) string {
	return fmt.Sprintf("%s/ara/content/listen/%s/channel.mp3", RadioGardenAPIUrl, stationId)
}

func (r *RadioGarden) GetStation(stationId string) (*radiogarden.Channel, error) {
	res, err := r.client.GetAraContentChannelChannelIdWithResponse(
		context.Background(),
		stationId,
	)
	if err != nil {
		return nil, err
	}

	return res.JSON200.Data, nil
}

func (r *RadioGarden) SearchStations(search string) ([]Search, error) {
	res, err := r.client.GetSearchWithResponse(
		context.Background(),
		&radiogarden.GetSearchParams{
			Q: search,
		},
	)
	if err != nil {
		return nil, err
	}

	var stations []Search
	for _, hit := range *res.JSON200.Hits.Hits {
		if *hit.Source.Type == "channel" {
			if *hit.Score < 150 {
				continue
			}

			stations = append(stations, Search{
				Score:    *hit.Score,
				Code:     *hit.Source.Code,
				Subtitle: *hit.Source.Subtitle,
				Title:    *hit.Source.Title,
				Type:     *hit.Source.Type,
				Url:      *hit.Source.Url,
			})
		}
	}

	return stations, nil
}

func (r *RadioGarden) SearchPlaces(search string) ([]Search, error) {
	res, err := r.client.GetSearchWithResponse(
		context.Background(),
		&radiogarden.GetSearchParams{
			Q: search,
		},
	)
	if err != nil {
		return nil, err
	}

	var places []Search
	for _, hit := range *res.JSON200.Hits.Hits {
		if *hit.Source.Type == "place" {
			if *hit.Score < 150 {
				continue
			}

			places = append(places, Search{
				Score:    *hit.Score,
				Code:     *hit.Source.Code,
				Subtitle: *hit.Source.Subtitle,
				Title:    *hit.Source.Title,
				Type:     *hit.Source.Type,
				Url:      *hit.Source.Url,
			})
		}
	}

	return places, nil
}
