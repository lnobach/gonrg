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

type Pusher struct {
	sched      *Scheduler
	rcv        chan d0.ParseableRawData
	listeners  map[string]chan *obis.OBISMappedResult
	listenerMu sync.Mutex
}

func NewPusher(sched *Scheduler) *Pusher {
	return &Pusher{
		sched:     sched,
		rcv:       make(chan d0.ParseableRawData),
		listeners: make(map[string]chan *obis.OBISMappedResult),
	}
}

func (p *Pusher) GetReceiver() chan d0.ParseableRawData {
	return p.rcv
}

func (p *Pusher) ParseAndPushForever(ctx context.Context, parsecfg *d0.ParseConfig) {
	for {
		select {
		case raw := <-p.rcv:
			p.safeParseAndPush(parsecfg, raw)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Pusher) safeParseAndPush(parsecfg *d0.ParseConfig, raw d0.ParseableRawData) {
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
	}
	p.setAndNotify(result)
}

func (p *Pusher) setAndNotify(result *obis.OBISListResult) {
	resultMap := obis.ListToMap(result)
	prev := p.sched.SetAndGetPrevious(resultMap)
	change := obis.GetChanged(prev.GetList(), result)
	if change != nil {
		p.NotifyChange(change)
	}
}

func (p *Pusher) NotifyChange(change *obis.OBISListResult) {
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

func (p *Pusher) AddListener(id string, receiver chan *obis.OBISMappedResult) {
	log.WithFields(log.Fields{"id": id}).
		Debug("adding listener for meter updates")
	p.listenerMu.Lock()
	defer p.listenerMu.Unlock()
	p.listeners[id] = receiver
}

func (p *Pusher) DeleteListener(id string) {
	log.WithFields(log.Fields{"id": id}).
		Debug("removing listener for meter updates")
	p.listenerMu.Lock()
	defer p.listenerMu.Unlock()
	delete(p.listeners, id)
}
