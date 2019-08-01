package main

import (
	"github.com/astaxie/beego/logs"
	"sync"
	"sync/atomic"
	"time"
)

type statiBenchmarkProducerSnapshot struct {
	sendRequestSuccessCount     int64
	sendRequestFailedCount      int64
	receiveResponseSuccessCount int64
	receiveResponseFailedCount  int64
	sendMessageSuccessTimeTotal int64
	sendMessageMaxRT            int64
	createdAt                   time.Time
	next                        *statiBenchmarkProducerSnapshot
}

type produceSnapshots struct {
	sync.RWMutex
	head, tail, cur *statiBenchmarkProducerSnapshot
	len             int
}

func (s *produceSnapshots) TakeSnapshot() {
	b := s.cur
	sn := new(statiBenchmarkProducerSnapshot)
	sn.sendRequestSuccessCount = atomic.LoadInt64(&b.sendRequestSuccessCount)
	sn.sendRequestFailedCount = atomic.LoadInt64(&b.sendRequestFailedCount)
	sn.receiveResponseSuccessCount = atomic.LoadInt64(&b.receiveResponseSuccessCount)
	sn.receiveResponseFailedCount = atomic.LoadInt64(&b.receiveResponseFailedCount)
	sn.sendMessageSuccessTimeTotal = atomic.LoadInt64(&b.sendMessageSuccessTimeTotal)
	sn.sendMessageMaxRT = atomic.LoadInt64(&b.sendMessageMaxRT)
	sn.createdAt = time.Now()

	s.Lock()
	if s.tail != nil {
		s.tail.next = sn
	}
	s.tail = sn
	if s.head == nil {
		s.head = s.tail
	}

	s.len++
	if s.len > 10 {
		s.head = s.head.next
		s.len--
	}
	s.Unlock()
}

func (s *produceSnapshots) TrintStati() {
	s.RLock()
	if s.len < 10 {
		s.RUnlock()
		return
	}

	f, l := s.head, s.tail
	respSucCount := float64(l.receiveResponseSuccessCount - f.receiveResponseSuccessCount)
	sendTps := respSucCount / l.createdAt.Sub(f.createdAt).Seconds()
	avgRT := float64(l.sendMessageSuccessTimeTotal-f.sendMessageSuccessTimeTotal) / respSucCount
	maxRT := atomic.LoadInt64(&s.cur.sendMessageMaxRT)
	s.RUnlock()

	logs.Debug(
		"SendTPS: %d Max RT: %d Average RT: %7.3f Send Failed: %d Response Failed: %d Total:%d\n",
		int64(sendTps), maxRT, avgRT, l.sendRequestFailedCount, l.receiveResponseFailedCount, l.receiveResponseSuccessCount,
	)
}
