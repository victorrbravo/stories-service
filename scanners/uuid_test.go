package scanners

import (
	"testing"

	"stories/database"
	"stories/types"
)

func TestUUID_ScanRow(t *testing.T) {
	var query database.IQuery
	var scanner UUID
	var err error
	query = database.NewQuery()
	_, err = query.QueryOne("SELECT uuid_generate_v4()", &scanner)
	if err != nil {
		t.Fatal(err)
	}
	if types.IsNil(scanner.Value) {
		t.Fatal("unexpected `NULL' uuid value")
	}
}

func TestUUID_ScanRow2(t *testing.T) {
	var query database.IQuery
	var scanner UUID
	var err error
	query = database.NewQuery()
	_, err = query.QueryOne("SELECT NULL", &scanner)
	if err != nil {
		t.Fatal(err)
	}
	if !types.IsNil(scanner.Value) {
		t.Fatal("uuid value must be `NULL' but it's not")
	}
}
