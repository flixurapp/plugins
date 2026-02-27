package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/flixurapp/flixur/pluginkit"
	protobuf "github.com/flixurapp/flixur/proto/go"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"go.uploadedlobster.com/musicbrainzws2"
)

var INFO = protobuf.PacketInfo{
	Id:          "musicbrainz",
	Version:     1,
	Features:    []protobuf.Features{protobuf.Features_ARTIST_SEARCH},
	Name:        "MusicBrainz",
	Description: "Integration with MusicBrainz.org",
	Author:      "xela.codes",
}

var MusicBrainz *musicbrainzws2.Client

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "3:04:05PM",
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("[%s] %s", INFO.Id, i)
		},
	})

	log.Info().Msg("Initializing plugin...")

	if err := pluginkit.WriteMessage(&protobuf.PluginPacket{
		Id:   ulid.Make().String(),
		Type: protobuf.PacketType_INFO,
	}, &INFO, os.Stdout); err != nil {
		log.Err(err).Msg("Failed to write info packet.")
		panic(0)
	}

	listener := pluginkit.StartReadingPackets(os.Stdin, func(err error) {
		log.Err(err).Msg("Failed to read packet from stdin.")
	})
	pluginkit.AddPacketListener(listener, protobuf.PacketType_INIT, func(data *protobuf.PacketInit, pkt *protobuf.PluginPacket) {
		log.Info().Interface("d", data).Msg("Initializing MusicBrainz client...")
		MusicBrainz = musicbrainzws2.NewClient(musicbrainzws2.AppInfo{
			Name:    "Flixur MusicBrainz Plugin",
			Version: fmt.Sprintf("v%d,%d", INFO.Version, data.GetVersion()),
		})
	})
	pluginkit.AddPacketListener(listener, protobuf.PacketType_DESTROY, func(data *protobuf.PacketDestroy, pkt *protobuf.PluginPacket) {
		log.Info().Msg("Destroying...")
		if MusicBrainz != nil {
			MusicBrainz.Close()
		}
	})

	pluginkit.ImplementFeature(listener, protobuf.Features_ARTIST_SEARCH, func(req *protobuf.FeatureArtistSearchRequest, _ *protobuf.PluginPacket) (*protobuf.FeatureArtistSearchResponse, *protobuf.FeatureError) {
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

	// never exit
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
