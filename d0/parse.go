package d0

import (
	"errors"
	"fmt"
	"maps"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lnobach/gonrg/obis"
	log "github.com/sirupsen/logrus"
)

type parserImpl struct {
	config *ParseConfig
}

func NewParser(config ParseConfig) (Parser, error) {
	c := &config
	err := parseSetDefaults(c)
	if err != nil {
		return nil, fmt.Errorf("failure setting config: %w", err)
	}
	return parserImpl{config: c}, nil
}

func parseSetDefaults(_ *ParseConfig) error {
	return nil
}

func (p parserImpl) GetOBISMap(data ParseableRawData, measurementTime time.Time) (*obis.OBISMappedResult, error) {

	obismap_exact := make(obis.OBISMap)
	obismap := make(obis.OBISMap)
	obislist := make(obis.OBISList, 0, 20)

	deviceid, err := data.ParseObis(p.config, func(e *obis.OBISEntry) error {
		obismap_exact[e.ExactKey] = e
		obismap[e.SimplifiedKey] = e
		if e.Name != "" {
			obismap[e.Name] = e
		}
		obislist = append(obislist, e)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing obis data: %w", err)
	}

	maps.Copy(obismap, obismap_exact)

	return &obis.OBISMappedResult{DeviceID: deviceid, MeasurementTime: measurementTime,
		List: obislist, Map: obismap}, nil

}

func (p parserImpl) GetOBISList(data ParseableRawData, measurementTime time.Time) (*obis.OBISListResult, error) {

	obislist := make(obis.OBISList, 0, 20)

	deviceid, err := data.ParseObis(p.config, func(e *obis.OBISEntry) error {
		obislist = append(obislist, e)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing obis data: %w", err)
	}

	return &obis.OBISListResult{DeviceID: deviceid, MeasurementTime: measurementTime,
		List: obislist}, nil

}

type RawData struct {
	raw string
}

func RawDataFromString(raw string) *RawData {
	return &RawData{raw: raw}
}

var (
	r_obis = regexp.MustCompile(`((?:[0-9]+-[0-9]+:)?((?:[0-9]+\.)?[0-9]+\.[0-9]+)(?:\*[0-9]+)?)\(([^\)\n]*)\)`)
	r_val  = regexp.MustCompile(`^([+-]?[0-9]+)(\.([0-9]+))?(\*([A-Za-z0-9_-]+))?$`)
)

func (d *RawData) ParseObis(cfg *ParseConfig, foundSet func(*obis.OBISEntry) error) (string, error) {

	beg := strings.SplitN(d.raw, "\n", 2)
	if !strings.HasPrefix(beg[0], "/") {
		return "", errors.New("no preamble in data")
	}
	deviceid := beg[0][1:]
	if len(beg) == 1 {
		return "", errors.New("no data after preamble")
	}

	elems := r_obis.FindAllStringSubmatch(beg[1], -1)
	for _, elem := range elems {

		name := obis.GetFromCatalogue(elem[1])
		if name == "" {
			name = obis.GetFromCatalogue(elem[2])
		}

		e := &obis.OBISEntry{
			ExactKey:      elem[1],
			SimplifiedKey: elem[2],
			Name:          name,
		}

		valRaw := elem[3]
		val_elems := r_val.FindStringSubmatch(valRaw)

		if val_elems != nil {
			var err error
			pre_dot := val_elems[1]
			post_dot := val_elems[3]
			int_val := pre_dot + post_dot
			e.ValueNum, err = strconv.ParseInt(int_val, 10, 64)
			if err != nil {
				if cfg.StrictMode {
					return "", fmt.Errorf("could not parse integer value '%s': %w", int_val, err)
				} else {
					log.WithError(err).Errorf("could not parse integer value '%s'", int_val)
					continue
				}
			}
			e.ValueScale = -len(post_dot)
			obis.Floatify(e)
			e.Unit = val_elems[5]
		} else {
			e.ValueText = valRaw
		}

		err := foundSet(e)
		if err != nil {
			if cfg.StrictMode {
				return "", fmt.Errorf("could not use obis set: %w", err)
			} else {
				log.WithError(err).Errorf("could not use obis set")
				continue
			}
		}

	}

	return deviceid, nil

}
