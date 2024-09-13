package listener

import (
	"bufio"
	"context"
	"fmt"
	"os"
)

type QueryHandler interface {
	HandleQuery(ctx context.Context, query string) string
}

type Listener struct {
	qh QueryHandler
}

func New(qh QueryHandler) *Listener {
	return &Listener{
		qh: qh,
	}
}

func (l *Listener) StartListening(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Println("Enter a query:")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			query := scanner.Text()
			resp := l.qh.HandleQuery(ctx, query)
			fmt.Println(resp)
		}
	}
}
