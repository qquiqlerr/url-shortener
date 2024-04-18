package url

import (
	"log/slog"
	"net/http"
)

type request struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}

type response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (id int64, err error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
