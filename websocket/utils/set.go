package utils

import "sync"

type Set struct {
	data map[string]struct{}
	sync.RWMutex
}

func NewSet(values ...string) *Set {
	s := new(Set)
	s.data = make(map[string]struct{})
	s.Lock()
	defer s.Unlock()
	for _, value := range values {
		s.data[value] = struct{}{}
	}
	return s
}

func (s *Set) Add(value string) {
	s.Lock()
	defer s.Unlock()
	s.data[value] = struct{}{}
}

func (s *Set) Remove(value string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, value)
}

func (s *Set) Has(value string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.data[value]
	return ok
}
