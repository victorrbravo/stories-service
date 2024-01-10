package settings

import (
	"os"
	"stories/actions/constants"
	"strings"
)

// Variable SHOW_EMPTY_TASKS

var ShowEmptyTasks = func() bool {
	value := os.Getenv("SHOW_EMPTY_TASKS")
	if strings.Compare(value, constants.True) == 0 {
		return true
	}
	return false
}
