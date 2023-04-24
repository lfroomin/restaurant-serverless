package print

import (
	"encoding/json"
	"log"
)

func Json(label string, data any) {
	str, _ := json.Marshal(data)
	log.Printf("%s: %s\n", label, str)
}
