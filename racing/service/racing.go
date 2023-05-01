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

// GetRace returns a single race by its ID.
func (s *racingService) GetRace(ctx context.Context, req *racing.GetRaceRequest) (*racing.GetRaceResponse, error) {
	race, err := s.racesRepo.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &racing.GetRaceResponse{Race: race}, nil
}
func (s *racingService) mustEmbedUnimplementedRacingServer() {}