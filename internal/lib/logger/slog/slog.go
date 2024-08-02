package slog_lib

import "log/slog"

func AddErrorAtribute(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
