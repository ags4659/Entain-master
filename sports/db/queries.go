package db

const (
	eventsList = "list"
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventsList: `
			SELECT 
				id, 
				event_id, 
				name, 
				location, 
				team, 
				winner, 
				runner_up, 
				advertised_start_time 
			FROM events
		`,
	}
}
