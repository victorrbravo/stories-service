package pg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {

	var err error
	_, err = Parse("fk_status='ToDo'", "fk_status='Progress'")
	//_, err = Parse("fk_status='ToDo'", "'Progress'=fk_status")
	//_, err = Parse("'ToDo'=fk_status or update_at='2023-12-12'", "fk_status='Progress'")
	//_, err = Parse("fk_status='ToDo' or update_at='2023-12-12'", "fk_status='Progress' and pk_story>0")
	//_, err = Parse("'ToDo'=fk_status or update_at='2023-12-12'", "fk_status='Progress'")

	assert.NoError(t, err, "There should be no error")

}
