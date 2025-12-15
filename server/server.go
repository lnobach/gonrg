package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lnobach/gonrg/obis"
)

type serverImpl struct {
	debug  bool
	config *ServerConfig
	meters map[string]*Scheduler
}

func NewServer(config *ServerConfig, debug bool) (Server, error) {
	err := serverSetDefaults(config)
	if err != nil {
		return nil, fmt.Errorf("failure setting config: %w", err)
	}
	srv := &serverImpl{
		config: config,
		meters: make(map[string]*Scheduler),
		debug:  debug,
	}
	for _, m := range config.Meters {
		sched, err := NewScheduler(m)
		if err != nil {
			return nil, fmt.Errorf("could not create scheduler for meter %s: %w", m.Name, err)
		}
		srv.meters[m.Name] = sched
	}
	return srv, nil
}

func serverSetDefaults(_ *ServerConfig) error {
	return nil
}

func (s *serverImpl) ListenAndServe() error {
	if s.debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	err := router.SetTrustedProxies(s.config.TrustedProxies)
	if err != nil {
		return fmt.Errorf("error setting trusted proxies: %w", err)
	}
	router.GET("/meter/:meter", s.getMeter)
	router.GET("/meter/:meter/:obiskey", s.getMeterValue)
	return router.Run(s.config.ListenAddr)
}

func (s *serverImpl) getMeter(c *gin.Context) {

	meter := c.Param("meter")

	sched, exists := s.meters[meter]
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "meter not found or no data yet"})
		return
	}

	val, err := sched.GetValue()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value, see server logs"})
	}

	c.IndentedJSON(http.StatusOK, val.GetList())

}

func (s *serverImpl) getMeterValue(c *gin.Context) {
	meter := c.Param("meter")

	sched, exists := s.meters[meter]
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "meter not found or no data yet"})
		return
	}

	val, err := sched.GetValue()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value, see server logs"})
	}

	obiskey := c.Param("obiskey")
	entry, exists := val.Map[obiskey]
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "obis value not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, &obis.OBISSingleResult{
		MeasurementTime: val.MeasurementTime,
		DeviceID:        val.DeviceID,
		Entry:           entry,
	})

}
