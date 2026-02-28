package main

import (
	"context"

	pb "github.com/flixurapp/flixur/proto/go"
	"github.com/rs/zerolog/log"
	"go.uploadedlobster.com/musicbrainzws2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (p *Plugin) ArtistSearch(ctx context.Context, req *pb.ArtistSearchRequest) (*pb.ArtistSearchResponse, error) {
	if err := p.isInitialized(); err != nil {
		return nil, err
	}

	log.Info().Str("query", req.Query).Int32("limit", req.Limit).Msg("ArtistSearch called")

	res, err := p.client.SearchArtists(ctx, musicbrainzws2.SearchFilter{
		Query:  req.Query,
		Dismax: true,
	}, musicbrainzws2.Paginator{Limit: ClampLimit(int(req.Limit))})

	if err != nil {
		log.Err(err).Msg("MusicBrainz search failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	results := make([]*pb.Artist, len(res.Artists))
	for i, artist := range res.Artists {
		var location *string
		if artist.Area != nil {
			location = &artist.Area.Name
		}

		results[i] = &pb.Artist{
			Id:       string(artist.ID),
			Provider: PluginInfo.Id,
			Name:     artist.Name,
			Location: location,
		}
	}

	return &pb.ArtistSearchResponse{
		Results: results,
	}, nil
}
