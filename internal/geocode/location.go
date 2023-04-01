package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/location"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"strings"
)

type placeSearcher interface {
	SearchPlaceIndexForText(ctx context.Context, input *location.SearchPlaceIndexForTextInput, optFns ...func(*location.Options)) (*location.SearchPlaceIndexForTextOutput, error)
}

type LocationService struct {
	Client     placeSearcher
	PlaceIndex string
}

func (ls LocationService) Geocode(address model.Address) (model.Location, string, error) {

	text := join(address.Line1, address.Line2, address.City, address.State, address.ZipCode, address.Country)

	fmt.Printf("Geocode address: %s\n", text)

	input := &location.SearchPlaceIndexForTextInput{
		IndexName:  &ls.PlaceIndex,
		Text:       &text,
		MaxResults: 10,
	}

	data, err := ls.Client.SearchPlaceIndexForText(context.Background(), input)
	if err != nil {
		return model.Location{}, "", err
	}

	d, _ := json.Marshal(data)
	fmt.Printf("Location output: %s\n", d)

	loc := model.Location{}
	var timezoneName string
	if data != nil && len(data.Results) > 0 {
		place := data.Results[0].Place
		geocode := fmt.Sprintf("%f,%f", place.Geometry.Point[1], place.Geometry.Point[0])
		loc = model.Location{
			Geocode:       &geocode,
			AddressNumber: place.AddressNumber,
			Street:        place.Street,
			Municipality:  place.Municipality,
			PostalCode:    place.PostalCode,
			Region:        place.Region,
			SubRegion:     place.SubRegion,
			Country:       place.Country,
		}
		timezoneName = *place.TimeZone.Name
	}

	return loc, timezoneName, nil
}

func join(strs ...*string) string {
	var sb strings.Builder
	for _, str := range strs {
		if str != nil {
			sb.WriteString(*str + " ")
		}
	}
	return sb.String()
}
