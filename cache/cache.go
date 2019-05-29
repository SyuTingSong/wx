package cache

import (
	"math"
	"time"
)

type result struct {
	value interface{}
	ttl   uint
	err   error
}

type FetchFromSource func(key string, timeout uint) (value interface{}, ttl uint, err error)

type cachePromise struct {
	result   result
	expireAt time.Time
	done     chan struct{}
}

func (p *cachePromise) fetchFromSource(r request) {
	value, ttl, err := r.fetch(r.key, r.timeout)
	p.result = result{value, ttl, err}
	p.expireAt = time.Now().Add(time.Duration(ttl) * time.Second)
	close(p.done)
}

func (p *cachePromise) delivery(t chan<- result) {
	<-p.done
	timeDiff := p.expireAt.Sub(time.Now())
	var ttl uint = 0
	if timeDiff > 0 {
		ttl = uint(math.RoundToEven(float64(timeDiff) / float64(time.Second)))
	}
	t <- result{p.result.value, ttl, p.result.err}
}

const (
	cmdGet = iota
	cmdDel
)

type request struct {
	command  uint
	key      string
	timeout  uint
	fetch    FetchFromSource
	response chan<- result
}

type Store struct {
	requests chan request
	promises map[string]*cachePromise
}

func NewStore() *Store {
	store := &Store{
		requests: make(chan request),
		promises: make(map[string]*cachePromise),
	}
	go store.start()
	return store
}

func (s *Store) Close() {
	close(s.requests)
}

func (s *Store) start() {
	for r := range s.requests {
		switch r.command {
		case cmdGet:
			promise := s.promises[r.key]
			if promise == nil || time.Now().After(promise.expireAt) {
				promise = &cachePromise{
					expireAt: time.Now().Add(time.Duration(r.timeout+1) * time.Second),
					done:     make(chan struct{}),
				}
				s.promises[r.key] = promise
				go promise.fetchFromSource(r)
			}
			go promise.delivery(r.response)
		case cmdDel:
			if _, ok := s.promises[r.key]; ok {
				delete(s.promises, r.key)
			}
			close(r.response)
		}
	}
}

func (s *Store) Get(key string, fetch FetchFromSource, timeout uint) (val interface{}, ttl uint, err error) {
	res := make(chan result)
	s.requests <- request{
		cmdGet,
		key,
		timeout,
		fetch,
		res,
	}
	result := <-res
	return result.value, result.ttl, result.err
}

func (s *Store) ForceUpdate(key string, fetch FetchFromSource, timeout uint) (val interface{}, ttl uint, err error) {
	s.Del(key)
	return s.Get(key, fetch, timeout)
}

func (s *Store) Del(key string) {
	res := make(chan result)
	s.requests <- request{
		cmdDel,
		key,
		0,
		nil,
		res,
	}
	<-res
}
