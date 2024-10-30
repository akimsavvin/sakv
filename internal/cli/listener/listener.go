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
	queries := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			scanner.Scan()
			queries <- scanner.Text()
		}
	}()

	for {
		fmt.Println("Enter a query:")

		select {
		case <-ctx.Done():
			fmt.Println("Stopped listening")
			return ctx.Err()
		case query := <-queries:
			resp := l.qh.HandleQuery(ctx, query)
			fmt.Println(resp)
		}
	}
}
