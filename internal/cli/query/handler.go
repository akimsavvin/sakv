package query

import (
	"context"
	"errors"
	"io"
	"net"
)

type NetHandler struct {
	conn net.Conn
}

func NewNetHandler(conn net.Conn) *NetHandler {
	return &NetHandler{
		conn: conn,
	}
}

func (h *NetHandler) HandleQuery(ctx context.Context, query string) string {
	_, err := h.conn.Write([]byte(query))
	if err != nil {
		return h.err(err)
	}

	resp := make([]byte, 4096)
	_, err = h.conn.Read(resp)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return "[ERROR]: connection closed"
		}

		return h.err(err)
	}

	return string(resp)
}

func (h *NetHandler) err(err error) string {
	return "[ERROR]: " + err.Error()
}
