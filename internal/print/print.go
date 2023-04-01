package print

import (
	"encoding/json"
	"fmt"
)

func Json(label string, data any) {
	str, _ := json.Marshal(data)
	fmt.Printf("%s: %s\n", label, str)
}
