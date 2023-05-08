package service

import (
	"context"
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
)

type RacesRepo interface {
	List(ctx context.Context, in *racing.ListRacesRequest) ([]*racing.Race, error)
	// ...
  }
  


// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
	racing.UnimplementedRacingServer
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) *racingService {
	return &racingService{racesRepo: racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
    races, err := s.racesRepo.List(ctx, in)

    if err != nil {
        return nil, err
    }

    return &racing.ListRacesResponse{Races: races}, nil
}


func (s *racingService) mustEmbedUnimplementedRacingServer() {}