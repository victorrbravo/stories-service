package scanners_test

import (
	"testing"

	"stories/database"
	"stories/scanners"
)

func TestDate_ScanRow(t *testing.T) {
	var query database.IQuery
	var scanner scanners.Date
	var err error
	query = database.NewQuery()
	_, err = query.QueryOne("SELECT CURRENT_DATE", &scanner)
	if err != nil {
		t.Fatal(err)
	}
}
