package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type RendererFunc func(http.ResponseWriter, *http.Request) (render.Renderer, *ErrRepsonse)

func HandleRendererFunc(fn RendererFunc, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp render.Renderer

		okResp, errResp := fn(w, r)
		if errResp != nil {
			logger.ErrorContext(
				r.Context(),
				"error in handler",
				slog.Any("error", errResp),
			)

			resp = errResp
		} else {
			resp = okResp
		}

		if err := render.Render(w, r, resp); err != nil {
			logger.ErrorContext(r.Context(), "error rendering response", slog.Any("error", err))

			render.Render(w, r, RenderErrorResponse(err)) //nolint: errcheck,gosec // ignore error
		}
	}
}
