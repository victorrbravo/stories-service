package scanners_test

import (
	"testing"

	"stories/database"
	"stories/scanners"
	"stories/types"
)

func TestInt_ScanRow(t *testing.T) {
	var query database.IQuery
	var scanner scanners.Int
	var err error
	query = database.NewQuery()
	_, err = query.QueryOne("SELECT 1234", &scanner)
	if err != nil {
		t.Fatal(err)
	}
	if types.AsInt(scanner.Value) != 1234 {
		t.Fatalf("cannot scan a simple string value from the database")
	}
	_, err = query.QueryOne("SELECT NULL", &scanner)
	if err != nil {
		t.Fatal(err)
	}
	if !types.IsNil(scanner.Value) {
		t.Fatalf("failed to scan `NULL' as a `types.String'")
	}
}
