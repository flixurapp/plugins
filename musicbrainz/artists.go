package main

import (
	"context"
	"time"

	"github.com/flixurapp/flixur/pluginkit"
	protobuf "github.com/flixurapp/flixur/proto/go"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"go.uploadedlobster.com/musicbrainzws2"
)

func ImplementArtists() {
	pluginkit.ImplementFeature(Listener, protobuf.Features_ARTIST_SEARCH, func(req *protobuf.FeatureArtistSearchRequest, _ *protobuf.PluginPacket) (*protobuf.FeatureArtistSearchResponse, *protobuf.FeatureError) {
		if err := CheckInitialization(); err != nil {
			return nil, err
		}
		log.Debug().Interface("req", req).Msg("Received search request.")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if res, err := MusicBrainz.SearchArtists(ctx, musicbrainzws2.SearchFilter{
			Query:  req.GetQuery(),
			Dismax: true,
		}, musicbrainzws2.Paginator{Limit: ClampLimit(int(req.GetLimit()))}); err != nil {
			return nil, &protobuf.FeatureError{
				Code:    int32(protobuf.FeatureErrorCode_UNKNOWN),
				Message: err.Error(),
			}
		} else {
			return &protobuf.FeatureArtistSearchResponse{
				Results: lo.Map(res.Artists, func(artist musicbrainzws2.Artist, _ int) *protobuf.Artist {
					var area *string
					if artist.Area != nil {
						area = &artist.Area.Name
					}

					return &protobuf.Artist{
						Id:       string(artist.ID),
						Provider: INFO.Id,
						Name:     artist.Name,
						Location: area,
					}
				}),
			}, nil
		}
	})
}
