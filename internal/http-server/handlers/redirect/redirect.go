package redirect

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "storage.postgres.GetURL"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("redirect: missing alias")

			render.JSON(w, r, response.Error("missing alias"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			log.Error("failed to get url", "error", err)

			render.JSON(w, r, response.Error("internal error"))

			return
		}

		log.Info("redirect: got url", "url", resURL)

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
