package main

import (
	protobuf "github.com/flixurapp/flixur/proto/go"
	"go.uploadedlobster.com/musicbrainzws2"
)

func CheckInitialization() *protobuf.FeatureError {
	if MusicBrainz == nil {
		return &protobuf.FeatureError{
			Code:    0,
			Message: "Client Not Initialized",
		}
	}

	return nil
}

// clamps the limit to the maximum allowed by MB api
func ClampLimit(limit int) int {
	return min(musicbrainzws2.MaxLimit, limit)
}
