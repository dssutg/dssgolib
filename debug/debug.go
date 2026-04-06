// Package debug provides convenient routines for debugging.
package debug

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

func PrintJSON(x any) {
	data, err := json.Marshal(x, jsontext.WithIndent("  "), json.Deterministic(true))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
