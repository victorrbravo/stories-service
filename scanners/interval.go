package scanners

import (
	"github.com/victorrbravo/stories-service/database"
	"strings"
	"time"
)

type Interval struct {
	Value time.Duration
	raw   string
}

func (i *Interval) Scan(value interface{}) error {
	data := value.([]byte)
	str := string(data)
	labels := []rune{'h', 'm', 's'}
	if strings.Count(str, ":") > 2 {
		labels = append([]rune{'d'}, labels...)
	}
	forDuration := func() func(c rune) rune {
		p := 0
		return func(c rune) rune {
			if c == ':' && p < len(labels) {
				value := labels[p]
				p++
				return value
			}
			return c
		}
	}
	str = strings.Map(forDuration(), str)
	str = str + "s"
	duration, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	i.Value = duration
	return nil
}

func (scanner *Interval) ScanRow(row database.RowInterface) (interface{}, error) {
	var value Interval
	var err error
	err = row.Scan(&value)
	if err != nil {
		return nil, err
	}
	scanner.Value = value.Value
	return scanner, nil
}
