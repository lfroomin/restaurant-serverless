package httpResponse

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

import (
	"net/http"
)

func Test_New(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		code    int
		data    any
		expBody string
	}{
		{
			name: "happy path",
			code: http.StatusForbidden,
			data: struct {
				Field1 int
				Field2 string
			}{
				Field1: 1,
				Field2: "two",
			},
			expBody: `{"Field1":1,"Field2":"two"}`,
		},
		{
			name: "nil data",
			code: http.StatusForbidden,
		},
		{
			name: "error marshalling",
			code: http.StatusInternalServerError,
			data: make(chan int),
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := New(tc.code, tc.data)

			assert.Equal(t, tc.code, output.StatusCode)
			assert.Equal(t, tc.expBody, output.Body)
			assert.Equal(t, "true", output.Headers["Access-Control-Allow-Credentials"])
			assert.Equal(t, "*", output.Headers["Access-Control-Allow-Origin"])
		})
	}
}

func Test_NewNoEncode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		code    int
		msg     string
		expBody string
	}{
		{
			name:    "happy path",
			code:    http.StatusForbidden,
			msg:     `{"Field1":1,"Field2":"two"}`,
			expBody: `{"Field1":1,"Field2":"two"}`,
		},
		{
			name:    "empty message",
			code:    http.StatusForbidden,
			expBody: "",
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := NewNoEncode(tc.code, tc.msg)

			assert.Equal(t, tc.code, output.StatusCode)
			assert.Equal(t, tc.expBody, output.Body)
		})
	}
}

func Test_NewMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		code    int
		msg     string
		expBody string
	}{
		{
			name:    "happy path",
			code:    http.StatusForbidden,
			msg:     "happy",
			expBody: `{"Message":"happy"}`,
		},
		{
			name:    "empty message",
			code:    http.StatusForbidden,
			expBody: `{"Message":""}`,
		},
	}

	for _, tc := range testCases {
		// scoped variable
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := NewMessage(tc.code, tc.msg)

			assert.Equal(t, tc.code, output.StatusCode)
			assert.Equal(t, tc.expBody, output.Body)
		})
	}
}

func Test_NewBadRequest(t *testing.T) {
	t.Parallel()

	output := NewBadRequest("happy")
	assert.Equal(t, http.StatusBadRequest, output.StatusCode)
	assert.Equal(t, `{"Message":"happy"}`, output.Body)
}

func Test_NewServerError(t *testing.T) {
	t.Parallel()

	output := NewServerError("happy")
	assert.Equal(t, http.StatusInternalServerError, output.StatusCode)
	assert.Equal(t, `{"Message":"happy"}`, output.Body)
}
