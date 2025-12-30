package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lnobach/gonrg/obis"
	"github.com/lnobach/gonrg/util"
	"github.com/lnobach/gonrg/version"
	log "github.com/sirupsen/logrus"
)

type serverImpl struct {
	debug          bool
	config         *ServerConfig
	meters         map[string]*Scheduler
	sessioncounter int64
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

func (s *serverImpl) getWSUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			return slices.Contains(s.config.AllowOrigins, origin)
		},
	}
}

func (s *serverImpl) ListenAndServe() error {
	log.Infof("%s server, version %s", version.GonrgName, version.GonrgVersion)

	for mtr, msched := range s.meters {
		err := msched.Init()
		if err != nil {
			return fmt.Errorf("error setting up scheduler for meter %s: %w", mtr, err)
		}
	}

	if s.debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	corscfg := cors.DefaultConfig()
	corscfg.AllowOrigins = s.config.AllowOrigins
	router.Use(cors.New(corscfg))

	err := router.SetTrustedProxies(s.config.TrustedProxies)
	if err != nil {
		return fmt.Errorf("error setting trusted proxies: %w", err)
	}
	router.GET("/meters", s.getMeters)
	router.GET("/meter/:meter", s.getMeter)
	router.GET("/meter/:meter/:obiskey", s.getMeterValue)

	router.GET("/ws/meter/:meter", s.getPushHandler(false))
	router.GET("/ws/meter/:meter/:obiskey", s.getPushHandler(true))

	return router.Run(s.config.ListenAddr)
}

func (s *serverImpl) getPushHandler(obisval bool) func(c *gin.Context) {

	return func(c *gin.Context) {

		connid := s.newSessionID(c)

		conn, err := s.getWSUpgrader().Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.WithField("connid", connid).WithError(err).
				Debug("could not upgrade websocket session")
			return
		}
		defer util.LogDeferWarn(conn.Close)

		meter := c.Param("meter")

		var obiskey string
		if obisval {
			obiskey = c.Param("obiskey")
		}

		sched, exists := s.meters[meter]
		if !exists {
			log.WithFields(log.Fields{"connid": connid, "meter": meter}).
				Debug("could not find meter, terminating websocket")
			return
		}

		if sched.pusher == nil {
			log.WithFields(log.Fields{"connid": connid, "meter": meter}).
				Debug("push not supported for this meter, terminating websocket")
			return
		}

		rcv := make(chan *obis.OBISMappedResult)
		sched.pusher.AddListener(connid, rcv)
		defer sched.pusher.DeleteListener(connid)

	outer:
		for {
			select {
			case upd := <-rcv:

				var result any
				if obisval {
					if upd == nil {
						continue
					}
					entry, exists := upd.Map[obiskey]
					if !exists {
						continue
					}
					result = &obis.OBISSingleResult{
						MeasurementTime: upd.MeasurementTime,
						DeviceID:        upd.DeviceID,
						Entry:           entry,
					}
				} else {
					result = upd.GetList()
				}

				jsonMsg, err := json.Marshal(result)
				if err != nil {
					log.WithFields(log.Fields{"connid": connid, "meter": meter}).
						WithError(err).
						Debug("error marshaling json, terminating websocket")
					continue
				}
				err = conn.WriteMessage(websocket.TextMessage, jsonMsg)
				if err != nil {
					log.WithFields(log.Fields{"connid": connid, "meter": meter}).
						WithError(err).
						Debug("error writing message to socket, terminating socket")
					break outer
				}
			case <-c.Done():
				break outer
			}
		}

	}
}

func (s *serverImpl) getMeters(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, util.KeysFromMap(s.meters))
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
		log.WithError(err).Warn("error while trying to get value")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value, see server logs"})
		return
	}

	c.IndentedJSON(http.StatusOK, val.GetList())

}

func (s *serverImpl) getMeterValue(c *gin.Context) {
	meter := c.Param("meter")

	sched, exists := s.meters[meter]
	if !exists {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "meter not found"})
		return
	}

	val, err := sched.GetValue()
	if err != nil {
		log.WithError(err).Warn("error while trying to get value")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value, see server logs"})
		return
	}
	if val == nil {
		log.Debug("value we tried to get is nil")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value because no data"})
		return
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

func (s *serverImpl) newSessionID(c *gin.Context) string {
	sidnum := s.sessioncounter
	s.sessioncounter++
	return fmt.Sprintf("%s-%d", c.Request.RemoteAddr, sidnum)
}
