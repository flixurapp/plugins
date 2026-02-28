package main

import (
	"go.uploadedlobster.com/musicbrainzws2"
)

// clamps the limit to the maximum allowed by MB api
func ClampLimit(limit int) int {
	return min(musicbrainzws2.MaxLimit, limit)
}
