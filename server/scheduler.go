package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/obis"
	"github.com/lnobach/gonrg/sml"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type Scheduler struct {
	config   *ServedMeterConfig
	lastVal  *obis.OBISMappedResult
	lock     sync.RWMutex // for value access and non-cron device access
	cronlock sync.Mutex   // for device access only in cron mode
	device   d0.Device
	parser   d0.Parser
	cr       *cron.Cron
}

func NewScheduler(config *ServedMeterConfig) (*Scheduler, error) {

	var err error
	var device d0.Device
	if config.SML {
		device, err = sml.NewDevice(config.Device)
	} else {
		device, err = d0.NewDevice(config.Device)
	}
	if err != nil {
		return nil, fmt.Errorf("could not configure device: %w", err)
	}

	parser, err := d0.NewParser(config.Parser)
	if err != nil {
		return nil, fmt.Errorf("could not configure parser: %w", err)
	}

	var cr *cron.Cron

	if config.Cron != "" {
		cr = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional |
			cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)))
	}

	return &Scheduler{
		config: config,
		device: device,
		parser: parser,
		cr:     cr,
	}, nil

}

func (s *Scheduler) Init() error {
	if s.cr == nil {
		// init only for cron'ed meters
		return nil
	}
	log.Debugf("initializing cron for meter %s and fetch first value...", s.config.Name)

	result, err := s.fetchValueSafe()
	if err != nil {
		log.WithError(err).Errorf("error fetching initial value for meter %s, will retry with first cron", s.config.Name)
	}
	s.lastVal = result

	_, err = s.cr.AddFunc(s.config.Cron, s.fetchValueCron)
	if err != nil {
		return fmt.Errorf("could not initialize cron job: %w", err)
	}
	s.cr.Start()
	return nil
}

func (s *Scheduler) GetValue() (*obis.OBISMappedResult, error) {
	if s.cr != nil {
		// RLock slightly improves performance in the cron case
		s.lock.RLock()
		defer s.lock.RUnlock()
		return s.lastVal, nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.lastVal == nil || s.lastVal.MeasurementTime.Add(s.config.RateLimitMaxAge).Before(time.Now()) {
		if s.lastVal != nil {
			log.Debugf("value timed out (age %s), fetching new value for meter", s.lastVal.MeasurementTime)
		} else {
			log.Debugf("value has not yet been set")
		}
		result, err := s.fetchValueSafe()
		if err != nil {
			return nil, fmt.Errorf("could not fetch new measurement: %w", err)
		}
		s.lastVal = result
	} else {
		log.Debugf("value still cached, fetching from cache")
	}
	return s.lastVal, nil
}

func (s *Scheduler) fetchValueCron() {
	if !s.cronlock.TryLock() {
		log.Warnf("will not run fetch for meter %s because the previous job is still running", s.config.Name)
		return
	}
	defer s.cronlock.Unlock()
	start := time.Now()
	result, err := s.fetchValueSafe()
	if err != nil {
		result = nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lastVal = result
	if err != nil {
		log.WithError(err).Errorf("error fetching value for meter %s, duration %s", s.config.Name, time.Since(start))
	} else {
		log.Debugf("fetched value for meter %s, duration %s", s.config.Name, time.Since(start))
	}
}

func (s *Scheduler) fetchValueSafe() (result *obis.OBISMappedResult, err error) {
	defer func() {
		if pan := recover(); pan != nil {
			switch x := pan.(type) {
			case string:
				err = fmt.Errorf("fetch panicked: %s", x)
			case error:
				err = fmt.Errorf("fetch panicked: %w", x)
			default:
				err = fmt.Errorf("fetch panicked: %v", x)
			}
		}
	}()
	result, err = s.fetchValueUnsafe()
	return
}

func (s *Scheduler) fetchValueUnsafe() (*obis.OBISMappedResult, error) {
	rawVal, err := s.device.Get()
	now := time.Now()
	if err != nil {
		return nil, fmt.Errorf("could not get meter raw data from device: %w", err)
	}
	result, err := s.parser.GetOBISMap(rawVal, now)
	if err != nil {
		return nil, fmt.Errorf("could not parse raw data obtained from meter: %w", err)
	}
	log.Debug("successfully fetch new meter value")
	return result, nil
}
