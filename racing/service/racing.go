package service

import (
	"errors"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	// Default to ordering by advertised_start_time desc.
	if in.Order == nil {
		in.Order = &racing.OrderBy{
			Field:     "advertised_start_time",
			Direction: "desc",
		}
	}

	// Validate that the order direction is correct
	if in.Order != nil {
		if in.Order.Direction != "desc" && in.Order.Direction != "asc" {
			return nil, errors.New("invalid order by")
		}
	}

	races, err := s.racesRepo.List(in.Filter, in.Order)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}
