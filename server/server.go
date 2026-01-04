package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	meters         map[string]*scheduler
	sessioncounter int64
	srv            *http.Server
}

func NewServer(config *ServerConfig, debug bool) (Server, error) {
	err := serverSetDefaults(config)
	if err != nil {
		return nil, fmt.Errorf("failure setting config: %w", err)
	}
	srv := &serverImpl{
		config: config,
		meters: make(map[string]*scheduler),
		debug:  debug,
	}
	for _, m := range config.Meters {
		sched, err := newScheduler(m)
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
	log.Infof("%s server, version %s", version.GonrgName, version.GonrgVersion)

	for mtr, msched := range s.meters {
		err := msched.start()
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

	if len(s.config.AllowOrigins) > 0 {
		corscfg := cors.DefaultConfig()
		corscfg.AllowOrigins = s.config.AllowOrigins
		corscfg.AllowMethods = []string{"GET"}
		router.Use(cors.New(corscfg))
	}

	err := router.SetTrustedProxies(s.config.TrustedProxies)
	if err != nil {
		return fmt.Errorf("error setting trusted proxies: %w", err)
	}
	router.GET("/meters", s.getMeters)
	router.GET("/meter/:meter", s.getMeter)
	router.GET("/meter/:meter/:obiskey", s.getMeterValue)

	router.GET("/ws/meter/:meter", s.getPushHandler(false))
	router.GET("/ws/meter/:meter/:obiskey", s.getPushHandler(true))

	s.srv = &http.Server{
		Addr:    s.config.ListenAddr,
		Handler: router.Handler(),
	}
	return s.srv.ListenAndServe()
}

func (s *serverImpl) Shutdown(ctx context.Context) error {
	for mtr, msched := range s.meters {
		err := msched.stop(ctx)
		if err != nil {
			return fmt.Errorf("error shutting down up scheduler for meter %s: %w", mtr, err)
		}
	}
	return s.srv.Shutdown(ctx)
}

func (s *serverImpl) getPushHandler(obisval bool) func(c *gin.Context) {

	return func(c *gin.Context) {

		connid := s.newSessionID(c)

		upgr := &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return util.CorsWsIsOriginAllowed(r, s.config.AllowOrigins)
			},
		}

		conn, err := upgr.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.WithField("connid", connid).WithError(err).
				Debug("could not upgrade websocket session")
			return
		}
		defer util.LogDeferWarn(conn.Close)

		conn.SetReadLimit(1)

		go func(c *websocket.Conn) {
			for {
				_, _, err := c.ReadMessage()
				log.WithField("connid", connid).
					Debug("received message from client, but not supposed to, dropping")
				if err != nil {
					log.WithField("connid", connid).WithError(err).
						Debug("read channel closed")
					return
				}
			}
		}(conn)

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
		sched.pusher.addListener(connid, rcv)
		defer sched.pusher.deleteListener(connid)

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

	val, err := sched.getValue()
	if err != nil {
		log.WithError(err).Warn("error while trying to get value")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value, see server logs"})
		return
	}
	if val == nil {
		log.Debug("value we tried to get has not been set (yet)")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value because no data"})
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

	val, err := sched.getValue()
	if err != nil {
		log.WithError(err).Warn("error while trying to get value")
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "could not get value, see server logs"})
		return
	}
	if val == nil {
		log.Debug("value we tried to get has not been set (yet)")
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
