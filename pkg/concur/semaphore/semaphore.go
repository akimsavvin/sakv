package semaphore

import (
	"sync"
)

type Semaphore struct {
	c            *sync.Cond
	current, max int
}

func New(max int) Semaphore {
	if max < 1 {
		panic("semaphore: max cannot be less than 1")
	}

	return Semaphore{
		c:   sync.NewCond(new(sync.Mutex)),
		max: max,
	}
}

func (s *Semaphore) Acquire() {
	s.c.L.Lock()
	defer s.c.L.Unlock()

	for s.current == s.max {
		s.c.Wait()
	}
	s.current++
}

func (s *Semaphore) Release() {
	s.c.L.Lock()
	defer s.c.L.Unlock()

	if s.current == 0 {
		panic("semaphore: zero count")
	}

	s.current--
	s.c.Signal()
}

func (s *Semaphore) ReleaseAll() {
	s.c.L.Lock()
	defer s.c.L.Unlock()

	if s.current == 0 {
		panic("semaphore: zero count")
	}

	s.current = 0
	s.c.Broadcast()
}
