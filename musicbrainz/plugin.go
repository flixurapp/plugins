package main

import (
	"embed"
	"fmt"

	"forge.xela.codes/xela/flixur/pluginkit"
	pb "forge.xela.codes/xela/flixur/pluginkit/proto"
	"github.com/rs/zerolog/log"
	"go.uploadedlobster.com/musicbrainzws2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed manifest.json
var manifest embed.FS
var Manifest = pluginkit.LoadPluginManifest(manifest)

type Plugin struct {
	pb.UnimplementedFlixurPluginServiceServer
	client *musicbrainzws2.Client
}

// helper to test if the musicbrainz client is initialized before making requests
func (p *Plugin) isInitialized() error {
	if p.client == nil {
		return status.Error(codes.Unavailable, "MusicBrainz client not initialized")
	}
	return nil
}

func main() {
	pluginkit.SetupPluginLogger(Manifest)
	log.Info().Msg("MusicBrainz plugin starting...")

	plugin := &Plugin{
		client: musicbrainzws2.NewClient(musicbrainzws2.AppInfo{
			Name:    "Flixur MusicBrainz Plugin",
			Version: fmt.Sprintf("v%s", Manifest.Version),
		}),
	}
	defer plugin.client.Close()

	pluginkit.Serve(plugin)
}
