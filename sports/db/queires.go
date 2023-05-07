package db

const (
	eventList = "list"
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventList: `
			SELECT 
				id, 
				name, 
				category, 
				mascot, 
				advertised_start_time, 
				advertised_end_time 
			FROM events
		`,
	}
}
