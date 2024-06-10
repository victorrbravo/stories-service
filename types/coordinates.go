package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/wkb"
	"log"
	"regexp"
	"stories/constants"
	"strconv"
)

type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	notNil    bool
}

var NilCoordinates = Coordinates{-1, -1, false}

func (v *Coordinates) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	lon, lat, err := parseCoordinatesInYaml(value)
	if err != nil {
		return err
	}
	*v = Coordinates{lon, lat, true}
	return nil
}

func (coordinates *Coordinates) MarshalYAML() (interface{}, error) {
	if !coordinates.notNil {
		return constants.YAMLNil, nil
	} else {
		value := fmt.Sprintf("(%f,%f)", coordinates.Latitude, coordinates.Longitude)
		return value, nil
	}
}

func (c Coordinates) isNil() bool {
	return !c.notNil
}

func parseCoordinatesInYaml(input string) (float64, float64, error) {
	// Regular expression pattern to match float values inside parentheses
	re := regexp.MustCompile(`\(([-+]?[0-9]*\.?[0-9]+),\s*([-+]?[0-9]*\.?[0-9]+)\)`)

	// Find string submatches
	matches := re.FindStringSubmatch(input)

	// Check if matches were found
	if len(matches) < 3 {
		return 0, 0, fmt.Errorf("invalid input format")
	}

	// Convert the string matches to float64
	x, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing x value: %w", err)
	}

	y, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("error parsing y value: %w", err)
	}

	return x, y, nil
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
	if value == nil {
		coordinates.notNil = false
		return nil
	}
	geomBytes, ok := value.([]byte)
	g, err := ewkb.Unmarshal(geomBytes)
	if err != nil {
		log.Fatalf("failed to unmarshal WKB: %v", err)
	}
	// Type assert to a geom.Point
	point, ok := g.(*geom.Point)
	if !ok {
		return errors.New("expected a Point geometry")
	}

	coordinates.Longitude = point.Coords().X()
	coordinates.Latitude = point.Coords().Y()
	coordinates.notNil = true
	return nil
}

func (coordinates Coordinates) Value() (driver.Value, error) {
	//return fmt.Sprintf("POINT (%f %f)", coordinates.Latitude, coordinates.Longitude), nil
	if coordinates.notNil == false {
		return nil, nil
	}
	point := geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{coordinates.Longitude, coordinates.Latitude})

	wkbData, err := wkb.Marshal(point, binary.LittleEndian)
	if err != nil {
		log.Println("Failed to convert point to WKB: %v", err)
		return "", err
	}

	// Encode WKB to hexadecimal string
	hexString := fmt.Sprintf("%s", hex.EncodeToString(wkbData))
	return hexString, nil

}

func NewCoordinates(latitude float64, longitude float64) Coordinates {
	return Coordinates{Latitude: latitude, Longitude: longitude, notNil: true}
}

func NewNilCoordinates() Coordinates {
	return Coordinates{}
}
