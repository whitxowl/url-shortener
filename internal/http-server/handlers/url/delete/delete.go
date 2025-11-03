package delete

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, deleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			log.Info("missing alias")

			render.JSON(w, r, response.Error("missing alias"))

			return
		}

		err := deleter.DeleteURL(alias)
		if errors.Is(err, &storage.ErrNoUrlWithAlias{}) {
			log.Info("no url with alias", slog.String("alias", alias))

			render.JSON(w, r, response.Error("no url with such alias"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", "error", err)

			render.JSON(w, r, response.Error("failed to delete url"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		render.JSON(w, r, response.OK())
	}
}
