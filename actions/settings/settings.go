package settings

import (
	"log"
	"os"
	"path/filepath"
	"github.com/victorrbravo/stories-service/actions/constants"
	"strings"
)

// Variable SHOW_EMPTY_TASKS
const (
	TemporaryName string = "stories_files"
	DefautName    string = "wf.svg"
)

var ShowEmptyTasks = func() bool {
	value := os.Getenv("SHOW_EMPTY_TASKS")
	if strings.Compare(value, constants.True) == 0 {
		return true
	}
	return false
}

var WorflowURL = func() string {
	value := os.Getenv("WORKFLOW_URL")
	if len(value) == 0 {
		dir, err := os.MkdirTemp("", TemporaryName)
		if err != nil {
			log.Print("error creating temporary directory for URLWfFile")
			return ""
		}
		value = filepath.Join(dir, DefautName)
	}
	return value

}
