package filesegment

import (
	"os"
)

type FileSegment struct {
	file *os.File
}

func New(name string) (*FileSegment, error) {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	return &FileSegment{
		file: f,
	}, nil
}

func (s *FileSegment) Name() string {
	return s.file.Name()
}

func (s *FileSegment) Read(p []byte) (n int, err error) {
	return s.file.Read(p)
}

func (s *FileSegment) Write(p []byte) (int, error) {
	n, err := s.file.Write(p)
	if err != nil {
		return 0, err
	}

	return n, s.file.Sync()
}

func (s *FileSegment) Close() error {
	return s.file.Close()
}

func (s *FileSegment) Size() (uint, error) {
	stat, err := s.file.Stat()
	if err != nil {
		return 0, err
	}

	return uint(stat.Size()), nil
}
