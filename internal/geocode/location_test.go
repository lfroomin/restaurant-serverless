package geocode

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/location"
	"github.com/aws/aws-sdk-go-v2/service/location/types"
	"github.com/lfroomin/restaurant-serverless/internal/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_Geocode(t *testing.T) {

	addressNumber := "123"
	street := "street"
	line1 := addressNumber + " " + street
	line2 := "line2"
	city := "city"
	state := "state"
	zip := "zip"
	country := "country"
	address := model.Address{
		Line1:   &line1,
		Line2:   &line2,
		City:    &city,
		State:   &state,
		ZipCode: &zip,
		Country: &country,
	}

	geocode := "123.000000,456.000000"
	subRegion := "sub" + state
	locationExp := model.Location{
		Geocode:       &geocode,
		AddressNumber: &addressNumber,
		Street:        &street,
		Municipality:  &city,
		PostalCode:    &zip,
		Region:        &state,
		SubRegion:     &subRegion,
		Country:       &country,
	}

	testCases := []struct {
		name      string
		address   model.Address
		loc       model.Location
		stubError string
		errMsg    string
	}{
		{
			name:    "happy path",
			address: address,
			loc:     locationExp,
		},
		{
			name:      "error",
			stubError: "an error occurred",
			errMsg:    "an error occurred",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			lc := LocationService{
				Client:     placeSearcherStub{error: tc.stubError},
				PlaceIndex: "",
			}
			loc, timezoneName, err := lc.Geocode(tc.address)

			if tc.errMsg != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.errMsg, err.Error())
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.loc, loc)
				assert.NotNil(t, timezoneName)
			}
		})
	}
}

type placeSearcherStub struct {
	error string
}

func (s placeSearcherStub) SearchPlaceIndexForText(_ context.Context, input *location.SearchPlaceIndexForTextInput, _ ...func(*location.Options)) (*location.SearchPlaceIndexForTextOutput, error) {
	if s.error != "" {
		return nil, errors.New(s.error)
	}

	// This depends on input.Text that looks like "123 street line2 city state zip country"
	inputText := strings.Split(*input.Text, " ")

	geometry := types.PlaceGeometry{Point: []float64{456, 123}}
	timezoneStr := "timezone"
	timezone := types.TimeZone{Name: &timezoneStr, Offset: new(int32)}
	subRegion := "sub" + inputText[4]
	place := types.Place{
		Geometry:      &geometry,
		AddressNumber: &inputText[0],
		Street:        &inputText[1],
		Municipality:  &inputText[3],
		PostalCode:    &inputText[5],
		Region:        &inputText[4],
		SubRegion:     &subRegion,
		Country:       &inputText[6],
		TimeZone:      &timezone,
	}

	return &location.SearchPlaceIndexForTextOutput{Results: []types.SearchForTextResult{{Place: &place}}}, nil
}
