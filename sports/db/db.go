package db

import (
	"time"

	"syreclabs.com/go/faker"
    "math/rand"
)

func (s *sportsRepo) seed() error {
	statement, err := s.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INTEGER PRIMARY KEY, event_id INTEGER, name TEXT, location TEXT, team TEXT, winner TEXT, runner_up TEXT, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}
	event := []string{"Basketball", "Soccer", "Swimming", "Tennis"}
	for i := 1; i <= 100; i++ {
		statement, err = s.db.Prepare(`INSERT OR IGNORE INTO events(id, event_id, name, location, team, winner, runner_up, advertised_start_time) VALUES (?,?,?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Number().Between(1, 10),
				event[rand.Intn(len(event))],
				faker.Address().City(),
				faker.Team().Name(),
				faker.Name().Name(),
				faker.Name().Name(),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
