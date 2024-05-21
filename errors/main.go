package errors

import (
	"errors"
	"stories/constants"
)

var DuplicateObjectsError = errors.New(constants.SQLErrorDuplicateKey)
var NoTrackingFoundError = errors.New(constants.NoTrackingFound)
