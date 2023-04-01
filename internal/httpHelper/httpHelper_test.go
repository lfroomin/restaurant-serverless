package httpHelper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ResponseBodyMsg(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		expOutput string
	}{
		{"happy path",
			"happy",
			`{"Message":"happy"}`,
		},
		{"empty string",
			"",
			`{"Message":""}`,
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := ResponseBodyMsg(tc.input)

			assert.Equal(t, tc.expOutput, output)
		})
	}
}