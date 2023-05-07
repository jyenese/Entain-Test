package service

import (
	"errors"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
)

type Sports interface {
	// ListEvents will return a collection of races.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)
}

// sportsService implements the sports interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

// newSportService instantiates and returns a new sportService.
func newSportService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

func (s *sportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	// Default to ordering by advertised_start_time desc.
	if in.Order == nil {
		in.Order = &sports.OrderBy{
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

	races, err := s.sportsRepo.List(in.Filter, in.Order)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}
