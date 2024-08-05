package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/svetlana-mel/url-shortener/internal/config"
	generator "github.com/svetlana-mel/url-shortener/internal/lib/generators/alias"
	slog_lib "github.com/svetlana-mel/url-shortener/internal/lib/logger/slog"
	"github.com/svetlana-mel/url-shortener/internal/repository"
)

type URLSaver interface {
	SaveURL(urlString, alias string) error
	GetAlias(urlString string) (string, error)
}

type SaveRequest struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type SaveResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

func New(saver URLSaver, log *slog.Logger) http.HandlerFunc {
	// handler constructor
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SaveRequest

		// распарсим тело запроса, обработаем ошибки
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// ошибка: пустое тело запроса
				// server log
				log.Error("empty request body")

				// client response
				render.JSON(w, r, &SaveResponse{
					Status: "Error",
					Error:  "empty request body",
				})
				return
			}

			log.Error("failed to decode request body", slog_lib.AddErrorAtribute(err))
			render.JSON(w, r, &SaveResponse{
				Status: "Error",
				Error:  "failed to decode request",
			})
		}

		log.Info("request body decoded", slog.Any("request", req))

		// check if url already has alias
		alias, err := saver.GetAlias(req.URL)
		if alias != "" {
			// url already has alias
			log.Info("alias for url already exists")
			render.JSON(w, r, &SaveResponse{
				Status: "OK",
				Alias:  alias,
			})
			return
		}
		if err != nil && !errors.Is(err, repository.ErrAliasNotFound) {
			log.Error("failed to check alias existence", slog_lib.AddErrorAtribute(err))
			render.JSON(w, r, &SaveResponse{
				Status: "Error",
				Error:  "failed to check alias existence",
			})
			return
		}

		cfg := config.NewConfig()

		alias = req.Alias
		if alias == "" {
			alias = generator.GenerateAlias(cfg.AliasLen)
		}

		err = saver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, repository.ErrAliasExists) {
				log.Error("generator error: get dublicate alias", slog_lib.AddErrorAtribute(err))
				render.JSON(w, r, &SaveResponse{
					Status: "Error",
					Error:  "generator error: get dublicate alias",
				})
				return
			}
			log.Error("failed to save url-alias pair", slog_lib.AddErrorAtribute(err))
			render.JSON(w, r, &SaveResponse{
				Status: "Error",
				Error:  "failed to save url-alias pair",
			})
			return
		}

		log.Info("alias-url pair successfully saved")
		render.JSON(w, r, &SaveResponse{
			Status: "OK",
			Alias:  alias,
		})
	}
}
