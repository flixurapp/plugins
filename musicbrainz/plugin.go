package main

import (
	"context"
	"fmt"

	"github.com/flixurapp/flixur/pluginkit"
	pb "github.com/flixurapp/flixur/proto/go"
	"github.com/rs/zerolog/log"
	"go.uploadedlobster.com/musicbrainzws2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var PluginInfo = &pb.PluginInfo{
	Id:          "musicbrainz",
	Name:        "MusicBrainz",
	Version:     1,
	Features:    []pb.Feature{pb.Feature_MUSIC_METADATA},
	Icon:        "simple-icons:musicbrainz",
	Description: "Integration with MusicBrainz.org",
	Author:      "xela.codes",
	Url:         "https://musicbrainz.org",
}

type Plugin struct {
	pb.UnimplementedFlixurPluginServer
	client *musicbrainzws2.Client
}

func (p *Plugin) GetPluginInfo(ctx context.Context) (*pb.PluginInfo, error) {
	return PluginInfo, nil
}

// helper to test if the musicbrainz client is initialized before making requests
func (p *Plugin) isInitialized() error {
	if p.client == nil {
		return status.Error(codes.Unavailable, "MusicBrainz client not initialized")
	}
	return nil
}

func main() {
	pluginkit.SetupPluginLogger(PluginInfo)
	log.Info().Msg("MusicBrainz plugin starting...")

	plugin := &Plugin{
		client: musicbrainzws2.NewClient(musicbrainzws2.AppInfo{
			Name:    "Flixur MusicBrainz Plugin",
			Version: fmt.Sprintf("v%d", PluginInfo.Version),
		}),
	}
	defer plugin.client.Close()

	pluginkit.Serve(plugin)
}
