package memorydb

import (
	"fmt"
	"sync"
)

type Set struct {
	m map[interface{}]bool
	sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		m: map[interface{}]bool{},
	}
}

func (s *Set) Add(item interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) AddMany(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, item := range items {
		s.m[item] = true
	}
}

func (s *Set) Pop() (interface{}, bool) {
	s.Lock()
	defer s.Unlock()
	if s.IsEmpty() {
		return nil, false
	}
	return s.ToList()[0], true
}

func (s *Set) Remove(item interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) ToList() []interface{} {
	s.RLock()
	defer s.RUnlock()
	list := []interface{}{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func (s *Set) Len() int {
	return len(s.ToList())
}

func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}

func (s *Set) ToString() string {
	return fmt.Sprintf("%v", s.ToList())
}

func (s *Set) Print() {
	fmt.Printf("len: %d items: %v", len(s.ToList()), s.ToList())
}
