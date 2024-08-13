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

type URLDeleter interface {
	DeleteURL(alias string) error
	GetURL(alias string) (string, error)
}

type DeleteURLResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	URL    string `json:"url,omitempty"`
}

func New(deleter URLDeleter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		resUrl, err := deleter.GetURL(alias)

		if err != nil {
			if errors.Is(err, repository.ErrURLNotFound) {
				log.Error("url not found")
				render.JSON(w, r, &DeleteURLResponse{
					Status: "Error",
					Error:  "no url with provided alias",
				})
				return
			}
			log.Error("error get url", slog_lib.AddErrorAtribute(err))
			render.JSON(w, r, &DeleteURLResponse{
				Status: "Error",
				Error:  "internal error",
			})
			return
		}

		err = deleter.DeleteURL(alias)

		if err != nil {
			log.Error("error get url", slog_lib.AddErrorAtribute(err))
			render.JSON(w, r, &DeleteURLResponse{
				Status: "Error",
				Error:  "internal error",
			})
			return
		}

		render.JSON(w, r, &DeleteURLResponse{
			Status: "OK",
			URL:    resUrl,
		})
	}
}
