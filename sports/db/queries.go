package db

const (
	eventList = "list"
)

// The getEventQueries function returns a map of queries that can be used to retrieve data from the database.
func getEventQueries() map[string]string {
	return map[string]string{
		eventList: `
			SELECT 
				id, 
				name, 
				category, 
				mascot, 
				advertised_start_time
			FROM events
		`,
	}
}
