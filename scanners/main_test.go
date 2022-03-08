package scanners_test

import (
	"testing"

	"stories/tests"
)

func TestMain(m *testing.M) {
	tests.PrepareDatabase()
	m.Run()
}
