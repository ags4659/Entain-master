package db

import (
	"database/sql"
	//"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"git.neds.sh/matty/entain/sports/proto/sports"
)

// SportsRepo provides repository access to sports.
type SportsRepo interface {
	// Init will initialise our sports repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error)
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
		// For test/example purposes, we seed the DB with some dummy events.
		err = s.seed()
	})

	return err
}

func (s *sportsRepo) List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = s.applyFilter(query, filter)

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

	if len(filter.EventIds) > 0 {
		clauses = append(clauses, "event_id IN ("+strings.Repeat("?,", len(filter.EventIds)-1)+"?)")

		for _, eventID := range filter.EventIds {
			args = append(args, eventID)
		}
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (s *sportsRepo) scanEvents(
	rows *sql.Rows,
) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time

		if err := rows.Scan(&event.Id, &event.EventId, &event.Name, &event.Location, &event.Team, &event.Winner, &event.RunnerUp, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts := timestamppb.New(advertisedStart)

		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}

	return events, nil
}
