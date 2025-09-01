package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
)

type FollowType string

const (
	FollowTypeRequested FollowType = "Requested"
	FollowTypeDenied    FollowType = "Denied"
	FollowTypeAccepted  FollowType = "Accepted"
	FollowTypeBlocked   FollowType = "Blocked"
	FollowTypeExpired   FollowType = "Expired"
)

func (e FollowType) IsValid() error {
	switch e {
	case FollowTypeRequested:
	case FollowTypeDenied:
	case FollowTypeAccepted:
	case FollowTypeBlocked:
	case FollowTypeExpired:
		return nil
	}
	return fmt.Errorf("Field must be type %s: %s", "FollowType", "Requested, Denied, Accepted, Blocked, Expired")
}

func (e FollowType) Next() FollowType {
	switch e {
	case FollowTypeRequested:
		return FollowTypeDenied
	case FollowTypeDenied:
		return FollowTypeAccepted
	case FollowTypeAccepted:
		return FollowTypeBlocked
	case FollowTypeBlocked:
		return FollowTypeExpired
	case FollowTypeExpired:

		return FollowTypeRequested
	}
	return e
}

func (e FollowType) Accepted() FollowType {
	if e == FollowTypeRequested {
		return FollowTypeAccepted
	}
	log.Printf("FollowType: %s not accepted", e)
	return e
}

func (e FollowType) Denied() FollowType {
	if e == FollowTypeRequested {
		return FollowTypeAccepted
	}
	log.Printf("FollowType: %s not denied", e)
	return e
}

func (e FollowType) Blocked() FollowType {
	return FollowTypeBlocked
}

func (e FollowType) Expired() FollowType {
	if e == FollowTypeRequested {
		return FollowTypeAccepted
	}
	log.Printf("FollowType: %s not expired", e)
	return e
}

func (e FollowType) isNil() bool {
	return len(string(e)) == 0
}

func (e FollowType) MarshalJSON() ([]byte, error) {
	if len(string(e)) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(string(e))
}

func (e *FollowType) UnmarshalJSON(src []byte) error {
	if string(src) != "null" {
		var s string
		if err := json.Unmarshal(src, &s); err != nil {
			return err
		}
		*e = FollowType(s)
	}
	return nil
}

func (e *FollowType) Scan(src interface{}) error {
	switch value := src.(type) {
	case nil:
		*e = ""
	case string:
		*e = FollowType(value)
	case []byte:
		*e = FollowType(string(value))
	default:
		return fmt.Errorf("Cannot read a FollowType from non-string")
	}
	return nil
}

func (e FollowType) Value() (driver.Value, error) {
	if len(string(e)) == 0 {
		return nil, nil
	}
	return string(e), nil
}
