package pg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {

	var err error
	var str string
	//str, err = Parse("fk_status='ToDo' and update_at='2023-12-12' and pk_story>0",
	//	"fk_status='Progress'", "fk_status='Done'",
	//	"Progress")
	//_, err = Parse("fk_status='ToDo'", "'Progress'=fk_status")
	//str, err = Parse("'ToDo'=fk_status and update_at='2023-12-12'", "fk_status='Progress'")
	//_, err = Parse("fk_status='ToDo' or update_at='2023-12-12'", "fk_status='Progress' and pk_story>0")
	//_, err = Parse("'ToDo'=fk_status or update_at='2023-12-12'", "fk_status='Progress'")

	fmt.Println("WHERE STR: |", str, "|")
	assert.NoError(t, err, "There should be no error")

}
