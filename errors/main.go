package errors

import (
	"errors"
	"github.com/victorrbravo/stories-service/constants"
)

var DuplicateObjectsError = errors.New(constants.SQLErrorDuplicateKey)
var NoTrackingFoundError = errors.New(constants.NoTrackingFound)
