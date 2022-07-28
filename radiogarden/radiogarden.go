package radiogarden

import (
	"context"
	"fmt"

	radiogarden "github.com/jonasrmichel/radio-garden-go"
)

type RadioGarden struct {
	client *radiogarden.ClientWithResponses
}

func New() *RadioGarden {
	client, err := radiogarden.NewClientWithResponses("https://radio.garden/api")
	if err != nil {
		panic(err)
	}

	return &RadioGarden{
		client: client,
	}
}

func (r *RadioGarden) GetStationStream(stationId string) string {
	return fmt.Sprintf("https://radio.garden/api/ara/content/listen/%s/channel.mp3", stationId)
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

func (r *RadioGarden) Search(search string) (*[]radiogarden.SearchResult, error) {
	res, err := r.client.GetSearchWithResponse(
		context.Background(),
		&radiogarden.GetSearchParams{
			Q: search,
		},
	)
	if err != nil {
		return nil, err
	}

	return res.JSON200.Hits.Hits, nil
}
