package scanners

import (
	"fmt"
	"testing"

	"stories/database"
	"stories/types"
)

const (
	str         = "string"
	flt         = 0.5
	typesString = "another string"
)

var (
	uuid                    = types.NewUUIDMust("cca25c3f-445d-4837-bdf7-8f3a86ade4f2")
	reflectSqlString string = fmt.Sprintf(
		"SELECT '%v' as \"Uuid\", '%s' as \"String\", %.1f AS \"Float\", '%s' AS \"TypesString\", null AS \"NullString\"",
		uuid,
		str,
		flt,
		typesString,
	)
)

type testType struct {
	String      string
	Float       float64
	TypesString types.String
	NullString  types.String
	Uuid        types.UUID
}

func TestReflect_ScanRow(t *testing.T) {
	var query database.IQuery
	var scanner Reflect
	value := testType{}
	var err error
	scanner = NewReflectScanner(&value)
	query = database.NewQuery()
	_, err = query.QueryMany(reflectSqlString, &scanner)
	if err != nil {
		t.Fatal(err)
	}
	if value.String != str {
		t.Fatalf("the `String' field should be populated with the value \"%s\"", str)
	}
	if value.Float != flt {
		t.Fatalf("the `Float' field should be populated with the value %.1f", flt)
	}
	if !types.UUIDEqual(value.Uuid, uuid) {
		t.Fatalf("the `Uuid' field should be populated with the value %v", uuid)
	}
	if !types.IsNil(value.NullString) {
		t.Fatalf("the `NullString' field should be NULL")
	}
}

func TestNewReflectScanner(t *testing.T) {
	var scanner Reflect
	var anything interface{}
	scanner = NewReflectScanner(&anything)
	if scanner.interfaceValue != &anything {
		t.Fatalf("scanner not initialzed correctly")
	}
}

// TODO: test `scanAndDiscover()'
