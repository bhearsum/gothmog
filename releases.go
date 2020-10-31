package main

import (
	"encoding/json"
)

// TODO: consider fully defining this struct again...
type Release map[string]interface{}

func parseRelease(data []byte) (Release, error) {
	var release Release
	err := json.Unmarshal(data, &release)
	return release, err
}
