package request

import (
	"errors"
	"io"
	"log/slog"

	"github.com/go-chi/render"
	slog_lib "github.com/svetlana-mel/url-shortener/internal/lib/logger/slog"
)

func Parse[T any](
	reqBody io.ReadCloser,
	log *slog.Logger,
	req *T,
) string {
	err := render.DecodeJSON(reqBody, &req)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// ошибка: пустое тело запроса
			// server log
			log.Error("empty request body")
			return "empty request body"
		}

		log.Error("failed to decode request body", slog_lib.AddErrorAtribute(err))

		return "failed to decode request"
	}

	return ""
}
