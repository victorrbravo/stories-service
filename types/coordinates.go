package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"stories/constants"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	notNil    bool
}

func (coordinates *Coordinates) UnmarshalJSON(input []byte) error {
	var values map[string]float64
	var err error
	var stringValue string
	stringValue = string(input)
	if stringValue == constants.JSONNil {
		coordinates.notNil = false
		return nil
	}
	values = make(map[string]float64)
	err = json.Unmarshal(input, &values)
	if err != nil {
		return err
	}
	if latitude, ok := values["latitude"]; ok {
		coordinates.Latitude = latitude
	} else {
		return fmt.Errorf("`latitude' is not present in the message")
	}
	if longitude, ok := values["longitude"]; ok {
		coordinates.Longitude = longitude
	} else {
		return fmt.Errorf("`latitude' is not present in the message")
	}
	coordinates.notNil = true
	return nil
}

func (coordinates Coordinates) MarshalJSON() ([]byte, error) {
	if !coordinates.notNil {
		return []byte(constants.JSONNil), nil
	} else {
		return json.Marshal(map[string]float64{
			"latitude":  coordinates.Latitude,
			"longitude": coordinates.Longitude,
		})
	}
}

func parseCoordinates(data []byte, byteOrder binary.ByteOrder) (Coordinates, error) {
	var geometryType uint32
	var err error
	var coordinates Coordinates
	err = binary.Read(bytes.NewReader(data[1:5]), byteOrder, &geometryType)
	if err != nil {
		return Coordinates{}, err
	}
	if geometryType != 1 {
		return Coordinates{}, fmt.Errorf("only 2D points are supported for the type `types.Coordinate'")
	}
	err = binary.Read(bytes.NewReader(data[5:13]), byteOrder, &coordinates.Latitude)
	if err != nil {
		return Coordinates{}, err
	}
	err = binary.Read(bytes.NewReader(data[13:]), byteOrder, &coordinates.Longitude)
	if err != nil {
		return Coordinates{}, err
	}
	return coordinates, nil
}

func (coordinates *Coordinates) Scan(value interface{}) error {
	var err error
	var byteOrder uint8
	var data []byte
	if value == nil {
		coordinates.notNil = false
		return nil
	}
	data, err = hex.DecodeString(string(value.([]byte)))
	if err != nil {
		return fmt.Errorf(constants.BadConversionFromByte, value)
	}
	byteOrder = data[0:1][0]
	// The first byte defines byteOrder, 0 means big-endian and 1 little endian
	switch byteOrder {
	case 0:
		*coordinates, err = parseCoordinates(data, binary.BigEndian)
		break
	case 1:
		*coordinates, err = parseCoordinates(data, binary.LittleEndian)
		break
	default:
		return fmt.Errorf("invalid endiannes, first byte must be either 01 or 00")
	}
	if err != nil {
		return err
	}
	coordinates.notNil = true
	return nil
}

func (coordinates Coordinates) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT (%f %f)", coordinates.Latitude, coordinates.Longitude), nil
}

func NewCoordinates(latitude float64, longitude float64) Coordinates {
	return Coordinates{Latitude: latitude, Longitude: longitude, notNil: true}
}

func NewNilCoordinates() Coordinates {
	return Coordinates{}
}
