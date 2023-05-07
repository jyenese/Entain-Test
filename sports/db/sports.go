package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
	_ "github.com/mattn/go-sqlite3"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error
	//List will return a list of sports.
	List(filter *sports.ListEventsRequestFilter, orderBy *racing.OrderBy) ([]*racing.Event, error)
}

type sportsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewSportsRepo creates a new sports repository.
func NewSportsRepo(db *sql.DB) SportsRepo {
	return &sportsRepo{db: db}
}

// Init prepares the sports repository dummy data.
func (s *sportsRepo) Init() error {
	var err error

	s.init.Do(func() {
		// Seeding in use for test/example purposes.
		err = s.seed()
	})

	return err
}

func (s *sportsRepo) List(filter *sports.ListEventsRequestFilter, orderBy *racing.OrderBy) ([]*racing.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()["list"]

	query, args = s.applyFilter(query, filter)

	if orderBy != nil {
		// We only allow ordering by the following fields.
		if orderBy.Field != "advertised_start_time" && orderBy.Field != "id" {
			return nil, fmt.Errorf("invalid order by field: %s", orderBy.Field)
		}
		query = fmt.Sprintf("%s ORDER BY %s %s", query, orderBy.Field, orderBy.Direction)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanEvents(rows)
}

func (s *sportsRepo) applyFilter(query string, filter *racing.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}
	if filter.Id != nil {
		clauses = append(clauses, "id = ?")
		args = append(args, filter.Id)
	}
	if filter.Name != nil {
		clauses = append(clauses, "name = ?")
		args = append(args, filter.Name)
	}
	if filter.Category != nil {
		clauses = append(clauses, "category = ?")
		args = append(args, filter.Category)
	}
	if filter.Mascot != nil {
		clauses = append(clauses, "mascot = ?")
		args = append(args, filter.Mascot)
	}

	return query, args
}

func (s *sportsRepo) scanEvents(rows *sql.Rows) ([]*racing.Event, error) {
	var (
		err    error
		events []*racing.Event
	)

	for rows.Next() {
		var (
			event sports.Event
			start time.Time
		)

		if err = rows.Scan(&event.Id, &event.Name, &event.Category, &event.Mascot, &start); err != nil {
			return nil, err
		}

		event.AdvertisedStartTime, err = ptypes.TimestampProto(start)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	return events, nil
}
