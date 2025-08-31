package api

import (
	"net/http"

	"github.com/go-chi/render"
)

type SuccessResponse struct {
	HttpResponseCode int    `json:"-"`
	StatusText       string `json:"status"`
	Data             any    `json:"data,omitempty"`
}

func (response SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, response.HttpResponseCode)
	return nil
}

func SuccessfulResponse(data any) render.Renderer {
	return &SuccessResponse{
		HttpResponseCode: http.StatusOK,
		StatusText:       "Success",
		Data:             data,
	}
}

func SuccessfulResponseWithCode(data any, code int) render.Renderer {
	return &SuccessResponse{
		HttpResponseCode: code,
		StatusText:       "Success",
		Data:             data,
	}
}

type ErrResponse struct {
	HttpResponseCode int    `json:"-"`
	StatusText       string `json:"status"`
	ErrorMessage     string `json:"error,omitempty"`
}

func (err ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, err.HttpResponseCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		HttpResponseCode: http.StatusBadRequest,
		StatusText:       "Bad request",
		ErrorMessage:     err.Error(),
	}
}

func ErrInvalidRequestWithCode(err error, code int) render.Renderer {
	return &ErrResponse{
		HttpResponseCode: code,
		StatusText:       "Bad request",
		ErrorMessage:     err.Error(),
	}
}

func ErrInternalServer() render.Renderer {
	return &ErrResponse{
		HttpResponseCode: http.StatusInternalServerError,
		StatusText:       "Internal server error",
		ErrorMessage:     "An error occured. Please try again later",
	}
}
