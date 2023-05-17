package utils

import "sync"

type Set struct {
	data map[int]struct{}
	sync.RWMutex
}

func NewSet(values ...int) *Set {
	s := new(Set)
	s.data = make(map[int]struct{})
	s.Lock()
	defer s.Unlock()
	for _, value := range values {
		s.data[value] = struct{}{}
	}
	return s
}

func (s *Set) Add(value int) {
	s.Lock()
	defer s.Unlock()
	s.data[value] = struct{}{}
}

func (s *Set) Remove(value int) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, value)
}

func (s *Set) Has(value int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.data[value]
	return ok
}
