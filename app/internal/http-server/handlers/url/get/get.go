package get

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	slog_lib "github.com/svetlana-mel/url-shortener/internal/lib/logger/slog"
	"github.com/svetlana-mel/url-shortener/internal/repository"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

type GetURLResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	URL    string `json:"alias,omitempty"`
}

func New(getter URLGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.get.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		resUrl, err := getter.GetURL(alias)

		if err != nil {
			if errors.Is(err, repository.ErrURLNotFound) {
				log.Error("url not found")
				render.JSON(w, r, &GetURLResponse{
					Status: "Error",
					Error:  "url not found",
				})
				return
			}
			log.Error("error get url", slog_lib.AddErrorAtribute(err))
			render.JSON(w, r, &GetURLResponse{
				Status: "Error",
				Error:  "internal error",
			})
			return
		}

		if resUrl == "" {
			log.Error("alias with empty url")
			render.JSON(w, r, &GetURLResponse{
				Status: "Error",
				Error:  "internal error",
			})
			return
		}

		render.JSON(w, r, &GetURLResponse{
			Status: "OK",
			URL:    resUrl,
		})
	}
}
