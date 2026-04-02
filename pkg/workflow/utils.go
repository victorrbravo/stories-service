package workflow

import (
	"log"
	db "stories/database"
)

func GetPgQuery(config db.DatabaseConfig) (db.PostgresQuery, error) {
	var query db.PostgresQuery
	query.Initialize(&config)
	err := query.Connect()
	if err != nil {
		log.Println("Db Connect error")
		return query, err
	}
	return query, nil
}

// SliceDiff returns the set difference of two slices (A - B).
// It returns a new slice containing all elements that exist in slice 'a'
// but are not present in slice 'b'.
//
// The function preserves the relative order of elements as they appear in 'a'.
// It uses a map-based lookup to ensure O(n + m) time complexity.
func SliceDiff[T comparable](a, b []T) []T {
	mb := make(map[T]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []T
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
