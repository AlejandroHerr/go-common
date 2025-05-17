package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type ErrRepsonse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string      `json:"status"`          // user-level status message
	ErrorText  string      `json:"error,omitempty"` // application-level error message, for debugging
	Details    interface{} `json:"details,omitempty"`
}

func NewErrorResponse(
	err error,
	statusCode int,
	statusText string,
	errorText string,
	details interface{},
) *ErrRepsonse {
	return &ErrRepsonse{
		Err:            err,
		HTTPStatusCode: statusCode,
		StatusText:     statusText,
		ErrorText:      errorText,
		Details:        details,
	}
}

func (e *ErrRepsonse) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("error", e.Err.Error()),
		slog.Int("http_status_code", e.HTTPStatusCode),
		slog.String("status_text", e.StatusText),
		slog.String("error_text", e.ErrorText),
		slog.Any("details", e.Details),
	)
}

func (e *ErrRepsonse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)

	return nil
}

func RenderErrorResponse(err error) *ErrRepsonse {
	return &ErrRepsonse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
		Details:        nil,
	}
}
