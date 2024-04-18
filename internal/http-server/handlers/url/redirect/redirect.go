package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/storage"
)

type URLGetter interface {
	GetURL(alias string) (url string, err error)
}

func New(log *slog.Logger, getter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "url.redirect.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, api.Error("Url is empty"))

			return
		}
		url, err := getter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				log.Error("url not found", slog.String("alias", alias))
				render.JSON(w, r, api.Error("Url not found"))
				return
			}
			log.Error("Internal error", slog.String("alias", alias))
			render.JSON(w, r, api.Error("Internal error"))
			return
		}
		log.Info("Get url", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
