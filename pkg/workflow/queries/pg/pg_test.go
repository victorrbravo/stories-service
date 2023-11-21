package pg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {

	var err error
	_, err = Parse("fk_status='ToDo' or update_at='2023-12-12'", "fk_status='Progress'")

	assert.NoError(t, err, "There should be no error")

}
