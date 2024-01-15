package debug

import (
	"encoding/json"
	"log"
)

func Die(v ...any) {
	t, _ := json.Marshal(v)
	log.Fatal(string(t))
}
