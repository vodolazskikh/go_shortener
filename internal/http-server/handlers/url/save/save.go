package save

import (
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	random "url-shortener/internal/lib/utils"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

const aliasLength = 8

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New "

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			errorMessage := "failed to decode request body"
			log.Error(errorMessage, sl.Err(err))
			render.JSON(w, r, resp.Error(errorMessage))

			return
		}

		if err := validator.New().Struct(req); err != nil {
			// @TODO привести ошибки к человекопонятному типу - что конкретно не провалидировалось
			errorMessage := "failed to validate request body"
			log.Error(errorMessage, sl.Err(err))
			render.JSON(w, r, resp.Error(errorMessage))

			return
		}

		alias := req.Alias

		if alias == "" {
			alias = random.RandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.URL, alias)

		if err != nil {
			errorMessage := "failed to save url"
			log.Info(errorMessage, sl.Err(err))
			render.JSON(w, r, resp.Error(errorMessage))

			return
		}

		log.Info("url saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})

	}
}
