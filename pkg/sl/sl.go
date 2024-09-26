package sl

import "log/slog"

func Err(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func Comp(comp string) slog.Attr {
	return slog.String("component", comp)
}
