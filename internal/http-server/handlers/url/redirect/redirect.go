package redirect

import (
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Response struct {
	resp.Response
	Url string
}

type URLGetter interface {
	GetUrl(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			errorMessage := "Empty alias"
			render.JSON(w, r, resp.Error(errorMessage))

			return
		}

		saverUrl, err := urlGetter.GetUrl(alias)

		if err != nil {
			errorMessage := "failed to get url"
			log.Error(errorMessage, sl.Err(err))
			render.Status(r, 404)
			render.JSON(w, r, resp.Error(errorMessage))

			return
		}

		log.Info("url getted", slog.String("url", saverUrl))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Url:      saverUrl,
		})
	}
}
