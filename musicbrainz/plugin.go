package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/flixurapp/flixur/pluginkit"
	protobuf "github.com/flixurapp/flixur/proto/go"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uploadedlobster.com/musicbrainzws2"
)

var INFO = protobuf.PacketInfo{
	Id:          "musicbrainz",
	Version:     1,
	MinVersion:  1,
	Features:    []protobuf.Features{protobuf.Features_ARTIST_SEARCH},
	Name:        "MusicBrainz",
	Icon:        "simple-icons:musicbrainz",
	Description: "Integration with MusicBrainz.org",
	Author:      "xela.codes",
}

var Listener pluginkit.PacketListenerAdder
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

	Listener = pluginkit.StartReadingPackets(os.Stdin, func(err error) {
		log.Err(err).Msg("Failed to read packet from stdin.")
	})
	pluginkit.AddPacketListener(Listener, protobuf.PacketType_INIT, func(data *protobuf.PacketInit, pkt *protobuf.PluginPacket) {
		log.Info().Interface("d", data).Msg("Initializing MusicBrainz client...")
		MusicBrainz = musicbrainzws2.NewClient(musicbrainzws2.AppInfo{
			Name:    "Flixur MusicBrainz Plugin",
			Version: fmt.Sprintf("v%d,%d", INFO.Version, data.GetVersion()),
		})
	})
	pluginkit.AddPacketListener(Listener, protobuf.PacketType_DESTROY, func(data *protobuf.PacketDestroy, pkt *protobuf.PluginPacket) {
		log.Info().Msg("Destroying...")
		if MusicBrainz != nil {
			MusicBrainz.Close()
		}
	})

	ImplementArtists()

	// never exit
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
