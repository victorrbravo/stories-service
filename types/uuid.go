package types

import (
	"database/sql/driver"
	"encoding/json"
	twinj "github.com/twinj/uuid"

	"stories/constants"
)

type UUID twinj.UUID

var nilUuid UUID = UUID{}

func NewUUID(original string) (UUID, error) {
	var value *twinj.UUID
	var err error
	value, err = twinj.Parse(original)
	if err != nil {
		return UUID{}, err
	}
	return UUID(*value), nil
}

func NewUUIDMust(original string) UUID {
	var value UUID
	var err error
	value, err = NewUUID(original)
	if err != nil {
		panic(err)
	}
	return value
}

func (uuid UUID) Value() (driver.Value, error) {
	var converted twinj.UUID = twinj.UUID(uuid)
	if twinj.Equal(converted, twinj.Nil) {
		return nil, nil
	}
	return converted.String(), nil
}

func (uuid *UUID) Scan(src interface{}) error {
	var parsed *twinj.UUID
	var err error
	switch input := src.(type) {
	case nil:
		*uuid = UUID{}
		return nil
	case string:
		parsed, err = twinj.Parse(input)
		break
	case []byte:
		parsed, err = twinj.Parse(string(input))
		break
	}
	if err != nil {
		return err
	}
	*uuid = UUID(*parsed)
	return nil
}

func (uuid UUID) String() string {
	var converted twinj.UUID = twinj.UUID(uuid)
	return converted.String()
}

func (uuid *UUID) UnmarshalJSON(src []byte) error {
	var decoded *twinj.UUID
	var err error
	var str string
	if string(src) == constants.JSONNil {
		*uuid = UUID(twinj.UUID{})
		return nil
	}
	err = json.Unmarshal(src, &str)
	if err != nil {
		return err
	}
	decoded, err = twinj.Parse(str)
	if err != nil {
		return err
	}
	*uuid = UUID(*decoded)
	return nil
}

func (uuid UUID) isNil() bool {
	return twinj.Equal(twinj.UUID(uuid), twinj.UUID{})
}

func NewNilUUID() UUID {
	return UUID{}
}

func (uuid UUID) MarshalJSON() (data []byte, err error) {
	if uuid.IsNil() {
		var p *UUID = nil
		return json.Marshal(p)
	}
	data, err = json.Marshal(uuid.String())
	return data, err
}

func UUIDEqual(uuid1 UUID, uuid2 UUID) bool {
	return twinj.Equal(twinj.UUID(uuid1), twinj.UUID(uuid2))
}

func (uuid UUID) IsNil() bool {
	return uuid == nilUuid
}
