package httpHelper

import (
	"encoding/json"
)

var CORSHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
}

func ResponseBodyMsg(msg string) string {
	data := map[string]string{"Message": msg}
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
