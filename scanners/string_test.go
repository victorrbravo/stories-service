package scanners

import (
	"testing"

	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/types"
)

func TestString_ScanRow(t *testing.T) {
	var query database.IQuery = database.NewQuery() // or appropriate constructor
	var scanner String
	var err error
	_, err = query.QueryOne("SELECT 'string value'", &scanner)
	if err != nil {
		t.Fatal(err)
	}
	if types.AsString(scanner.Value) != "string value" {
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
