package db

const (
	racesList = "list"
	raceByID = "by_id"
)

func getRaceQueries() map[string]string {
	return map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time,
				CASE 
					WHEN strftime('%s', advertised_start_time) <= strftime('%s', 'now') THEN 'CLOSED' 
					ELSE 'OPEN'
				END AS status
			FROM races
		`,
		// Will compute the status from db using strftime 
		// To get single race my meeting id 
		raceByID: `
			SELECT
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time,
				CASE 
					WHEN strftime('%s', advertised_start_time) <= strftime('%s', 'now') THEN 'CLOSED'
					ELSE 'OPEN'
				END AS status
			FROM races
			WHERE id = ?
		`,
	}
}


