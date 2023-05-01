package db

import (
    "context"
    "database/sql"
    "git.neds.sh/matty/entain/racing/proto/racing"
    _ "github.com/mattn/go-sqlite3"
    "github.com/golang/protobuf/ptypes"
    "strings"
    "sync"
    "time"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
    Init() error
    // Update the signature of the List method to match the proto definition.
    List(context.Context, *racing.ListRacesRequest) ([]*racing.Race, error)

    // Get: This method returns a single race identified by its unique identifier id.
    Get(ctx context.Context, id int64) (*racing.Race, error)
}

type racesRepo struct {
    db   *sql.DB
    init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
    return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
    var err error

    r.init.Do(func() {
        // For test/example purposes, we seed the DB with some dummy races.
        err = r.seed()
    })

    return err
}

func (r *racesRepo) List(ctx context.Context, req *racing.ListRacesRequest) ([]*racing.Race, error) {
    var (
        err   error
        query string
        args  []interface{}
    )

    query = getRaceQueries()[racesList]

    // Apply filters to the SQL query
    query, args = r.applyFilter(query, req.MeetingIds, req.Filter.VisibleOnly)

    // Add the order by clause to the SQL query if sorting is specified
    if req.SortBy != "" {
        query += " ORDER BY advertised_start_time " + req.SortBy
    }

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }

    return r.scanRaces(rows)
}

// Get returns a single race identified by its unique identifier id.
func (r *racesRepo) Get(ctx context.Context, id int64) (*racing.Race, error) {
    var (
        race           racing.Race
        advertisedTime time.Time
        status         string
    )

    // Get the SQL query to retrieve a single race by its ID from the query map.
    query := getRaceQueries()[raceByID]

    // Query the database with the race ID and scan the returned row into the race struct.
    err := r.db.QueryRowContext(ctx, query, id).Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedTime, &status)
    if err != nil {
        // Return nil if no rows are returned for the given race ID.
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    // Convert the advertised start time to a Protobuf timestamp.
    ts, err := ptypes.TimestampProto(advertisedTime)
    if err != nil {
        return nil, err
    }

    // Set the advertised start time and status of the race struct and return the race.
    race.AdvertisedStartTime = ts
    race.Status = status

    return &race, nil
}



// applyFilter returns the SQL query with the specified filters applied.
func (r *racesRepo) applyFilter(query string, meetingIDs []int64, visibleOnly bool) (string, []interface{}) {
    var (
        clauses []string
        args    []interface{}
    )

    if len(meetingIDs) > 0 {
        clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(meetingIDs)-1)+"?)")

        for _, meetingID := range meetingIDs {
            args = append(args, meetingID)
        }
    }

    if visibleOnly {
        clauses = append(clauses, "visible = 1")
    } else {
        clauses = append(clauses, "visible = 0")
    }

    if len(clauses) != 0 {
        query += " WHERE " + strings.Join(clauses, " AND ")
    }

    return query, args
}


func (m *racesRepo) scanRaces(
    rows *sql.Rows,
) ([]*racing.Race, error) {
    var races []*racing.Race

    for rows.Next() {
        var race racing.Race
        var advertisedStart time.Time
        var status string

        if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart, &status); err != nil {
            if err == sql.ErrNoRows {
                return nil, nil
            }

            return nil, err
        }

        ts, err := ptypes.TimestampProto(advertisedStart)
        if err != nil {
            return nil, err
        }

        race.AdvertisedStartTime = ts
        race.Status = status // For Status message 

        races = append(races, &race)
    }

    return races, nil
}
