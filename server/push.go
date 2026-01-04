package server

import (
	"context"
	"sync"
	"time"

	"github.com/lnobach/gonrg/d0"
	"github.com/lnobach/gonrg/obis"
	"github.com/lnobach/gonrg/util"
	log "github.com/sirupsen/logrus"
)

type Listener interface {
	OnChange(change *obis.OBISListResult)
}

type pusher struct {
	sched      *scheduler
	rcv        chan d0.ParseableRawData
	listeners  map[string]chan *obis.OBISMappedResult
	listenerMu sync.Mutex
}

func newPusher(sched *scheduler) *pusher {
	return &pusher{
		sched:     sched,
		rcv:       make(chan d0.ParseableRawData),
		listeners: make(map[string]chan *obis.OBISMappedResult),
	}
}

func (p *pusher) getReceiver() chan d0.ParseableRawData {
	return p.rcv
}

func (p *pusher) parseAndPushForever(ctx context.Context, parsecfg *d0.ParseConfig) {
	for {
		select {
		case raw := <-p.rcv:
			p.safeParseAndPush(parsecfg, raw)
		case <-ctx.Done():
			return
		}
	}
}

func (p *pusher) safeParseAndPush(parsecfg *d0.ParseConfig, raw d0.ParseableRawData) {
	defer func() {
		err := util.PanicToError(recover(), nil)
		if err != nil {
			log.WithError(err).Error("caught panic")
		}
	}()
	now := time.Now()
	result, err := d0.ParseOBISList(parsecfg, raw, now)
	if err != nil {
		log.WithError(err).Error("error parsing obis data")
		return
	}
	p.setAndNotify(result)
}

func (p *pusher) setAndNotify(result *obis.OBISListResult) {
	resultMap := obis.ListToMap(result)
	prev := p.sched.SetAndGetPrevious(resultMap)
	change := obis.GetChanged(prev.GetList(), result)
	if change != nil {
		p.notifyChange(change)
	}
}

func (p *pusher) notifyChange(change *obis.OBISListResult) {
	changeMap := obis.ListToMap(change)
	p.listenerMu.Lock()
	defer p.listenerMu.Unlock()
	for id, l := range p.listeners {
		select {
		case l <- changeMap:
		default:
			log.Warnf("listener id=%s: cannot keep up with changes, dropping...", id)
		}
	}
}

func (p *pusher) addListener(id string, receiver chan *obis.OBISMappedResult) {
	log.WithFields(log.Fields{"id": id}).
		Debug("adding listener for meter updates")
	p.listenerMu.Lock()
	defer p.listenerMu.Unlock()
	p.listeners[id] = receiver
}

func (p *pusher) deleteListener(id string) {
	log.WithFields(log.Fields{"id": id}).
		Debug("removing listener for meter updates")
	p.listenerMu.Lock()
	defer p.listenerMu.Unlock()
	delete(p.listeners, id)
}
