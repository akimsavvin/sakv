package filesegment

import (
	"fmt"
	"github.com/akimsavvin/sakv/internal/database/storage/wal"
	"time"
)

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) Create(dir string) (wal.Segment, error) {
	now := time.Now()
	name := fmt.Sprintf("%s/segment_%d.txt", dir, now.Unix())
	return New(name)
}
