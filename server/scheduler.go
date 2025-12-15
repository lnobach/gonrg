package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/obis"
	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	config  *ServedMeterConfig
	lastVal *obis.OBISMappedResult
	lock    sync.Mutex
	device  d0.Device
	parser  d0.Parser
}

func NewScheduler(config *ServedMeterConfig) (*Scheduler, error) {

	device, err := d0.NewDevice(config.Device)
	if err != nil {
		return nil, fmt.Errorf("could not configure device: %w", err)
	}

	parser, err := d0.NewParser(config.Parser)
	if err != nil {
		return nil, fmt.Errorf("could not configure parser: %w", err)
	}

	return &Scheduler{
		config: config,
		device: device,
		parser: parser,
	}, nil

}

func (s *Scheduler) GetValue() (*obis.OBISMappedResult, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now()
	if s.lastVal == nil || s.lastVal.MeasurementTime.Add(s.config.RateLimitMaxAge).Before(now) {
		if s.lastVal != nil {
			log.Debugf("value timed out (age %s), fetching new value for meter", s.lastVal.MeasurementTime)
		} else {
			log.Debugf("value has not yet been set")
		}
		err := s.fetchValueSafe()
		if err != nil {
			return nil, fmt.Errorf("could not fetch new measurement: %w", err)
		}
	} else {
		log.Debugf("value still cached, fetching from cache")
	}
	return s.lastVal, nil
}

func (s *Scheduler) fetchValueSafe() error {
	var err error
	defer func() {
		if pan := recover(); pan != nil {
			switch x := pan.(type) {
			case string:
				err = fmt.Errorf("fetch panicked %s", x)
			case error:
				err = fmt.Errorf("fetch panicked: %w", x)
			default:
				err = fmt.Errorf("fetch panicked: %v", x)
			}
		}
	}()
	err = s.fetchValueUnsafe()
	return err
}

func (s *Scheduler) fetchValueUnsafe() error {
	rawVal, err := s.device.Get()
	now := time.Now()
	if err != nil {
		return fmt.Errorf("could not get meter raw data from device: %w", err)
	}
	log.WithField("raw", rawVal).Debug("raw data received")
	result, err := s.parser.GetOBISMap(rawVal, now)
	if err != nil {
		return fmt.Errorf("could not parse raw data obtained from meter: %w", err)
	}
	s.lastVal = result
	log.Debug("successfully fetch new meter value")
	return nil
}
