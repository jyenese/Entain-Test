package db

import (
	"database/sql"
	"fmt"
	"sports/proto/sports"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error
	//List will return a list of sports.
	List(filter *sports.ListEventsRequestFilter, orderBy *sports.OrderBy) ([]*sports.Event, error)
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

func (s *sportsRepo) List(filter *sports.ListEventsRequestFilter, orderBy *sports.OrderBy) ([]*sports.Event, error) {
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

		// We only allow ordering by the following directions.
		// We default to desc if no direction is provided.
		// We also uppercase the direction to ensure we're consistent.
		query += " ORDER BY " + orderBy.Field + " " + strings.ToUpper(orderBy.Direction)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return s.scanEvents(rows)
}

func (s *sportsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
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

func (s *sportsRepo) scanEvents(rows *sql.Rows) ([]*sports.Event, error) {
	var (
		err    error
		events []*sports.Event
	)

	for rows.Next() {
		var (
			id                     int64
			name, category, mascot string
			advertisedStartTime    time.Time
		)

		err = rows.Scan(&id, &name, &category, &mascot, &advertisedStartTime)
		if err != nil {
			return nil, err
		}

		event := &sports.Event{
			Id:        id,
			Name:      name,
			Category:  category,
			Mascot:    mascot,
			StartTime: ptypes.TimestampNow(),
		}

		events = append(events, event)
	}

	return events, nil
}
