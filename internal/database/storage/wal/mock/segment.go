package walmock

import (
	"github.com/stretchr/testify/mock"
)

type Segment struct {
	mock.Mock
}

func NewSegment(name ...string) *Segment {
	s := new(Segment)
	if len(name) > 0 {
		s.On("Name").Return(name[0])
	}

	s.On("Close").Return(nil)

	return s
}

func (s *Segment) SetContent(content []byte) {
	l := len(content)
	s.On("Read", mock.AnythingOfType("[]uint8")).
		Run(func(args mock.Arguments) {
			p := args.Get(0).([]byte)
			for i := range l {
				p[i] = content[i]
			}
		}).
		Return(l, nil)

	s.On("Size").Return(uint(l), nil)

}

func (s *Segment) Name() string {
	args := s.Called()
	return args.Get(0).(string)
}

func (s *Segment) Read(p []byte) (n int, err error) {
	args := s.Called(p)
	return args.Int(0), args.Error(1)
}

func (s *Segment) Write(p []byte) (int, error) {
	args := s.Called(p)
	return args.Int(0), args.Error(1)
}

func (s *Segment) Close() error {
	args := s.Called()
	return args.Error(0)
}

func (s *Segment) Size() (uint, error) {
	args := s.Called()
	return args.Get(0).(uint), args.Error(1)
}
