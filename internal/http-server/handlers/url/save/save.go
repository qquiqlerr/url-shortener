package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"
)

// TODO: move to config
const alias_lenght = 5

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	api.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (id int64, err error)
}

func New(log *slog.Logger, urlSaver URLSaver, addr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "internal.http-server.handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, api.Error("failed to decode request body"))
			return
		}
		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", sl.Err(err))
			render.JSON(w, r, api.Error("invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.RandomString(alias_lenght)
		}
		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlExists) {
				log.Info("url already exists", slog.String("URL", req.URL))

				render.JSON(w, r, api.Error("url already exists"))

				return
			}
			log.Error("failed to add url", slog.String("URL", req.URL))
			render.JSON(w, r, api.Error("failed to add url"))
			return
		}

		log.Info("url added", slog.Int64("Id", id))

		render.JSON(w, r, Response{
			Response: api.OK(),
			Alias:    addr + "/" + alias,
		})
	}

}
