package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api"
)

type URLDeletter interface {
	RemoveURL(alias string) error
}

func New(log *slog.Logger, deletter URLDeletter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.delete.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		err := deletter.RemoveURL(alias)
		if err != nil {
			log.Error(err.Error(), slog.String("alias", alias))

			render.JSON(w, r, api.Error("Try again"))
			return
		}
		log.Info("Url removed", slog.String("alias", alias))

		render.JSON(w, r, api.OK())
	}
}
