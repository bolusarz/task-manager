package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
)

func TestSuccessfulResponses(t *testing.T) {
	tests := []struct {
		name          string
		data          any
		setup         func(data any) render.Renderer
		validate      func(w render.Renderer, data any)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, w render.Renderer)
	}{
		{
			name: "Successful Response: 200",
			data: map[string]string{"foo": "bar"},
			setup: func(data any) render.Renderer {
				return SuccessfulResponse(data)
			},
			validate: func(w render.Renderer, data any) {
				sr, ok := w.(*SuccessResponse)
				require.True(t, ok)

				require.Equal(t, sr.HttpResponseCode, http.StatusOK)
				require.Equal(t, sr.StatusText, "Success")
				require.Equal(t, sr.Data, data)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, w render.Renderer) {
				sr, ok := w.(*SuccessResponse)
				require.True(t, ok)

				require.Equal(t, recorder.Code, sr.HttpResponseCode)
			},
		},
		{
			name: "Successful Response: 201",
			data: map[string]string{"foo": "bar"},
			setup: func(data any) render.Renderer {
				return SuccessfulResponseWithCode(data, http.StatusCreated)
			},
			validate: func(w render.Renderer, data any) {
				sr, ok := w.(*SuccessResponse)
				require.True(t, ok)

				require.Equal(t, sr.HttpResponseCode, http.StatusCreated)
				require.Equal(t, sr.StatusText, "Success")
				require.Equal(t, sr.Data, data)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, w render.Renderer) {
				sr, ok := w.(*SuccessResponse)
				require.True(t, ok)

				require.Equal(t, recorder.Code, sr.HttpResponseCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := tt.setup(tt.data)
			tt.validate(sr, tt.data)

			router := chi.NewRouter()
			router.Get("/", func(w http.ResponseWriter, req *http.Request) {
				render.Render(w, req, sr)
			})

			httpRecorder := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			router.ServeHTTP(httpRecorder, req)

			tt.checkResponse(t, httpRecorder, sr)
		})
	}

}

func TestErrorResponses(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		setup         func(err error) render.Renderer
		validate      func(w render.Renderer, err error)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, w render.Renderer)
	}{
		{
			name: "Bad Request Response: 400",
			err:  fmt.Errorf("Invalid request"),
			setup: func(err error) render.Renderer {
				return ErrInvalidRequest(err)
			},
			validate: func(w render.Renderer, err error) {
				sr, ok := w.(*ErrResponse)
				require.True(t, ok)

				require.Equal(t, sr.HttpResponseCode, http.StatusBadRequest)
				require.Equal(t, sr.StatusText, "Bad request")
				require.Equal(t, sr.ErrorMessage, err.Error())
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, w render.Renderer) {
				sr, ok := w.(*ErrResponse)
				require.True(t, ok)

				require.Equal(t, recorder.Code, sr.HttpResponseCode)
			},
		},
		{
			name: "Internal Server Error Response: 500",
			err:  nil,
			setup: func(err error) render.Renderer {
				return ErrInternalServer()
			},
			validate: func(w render.Renderer, err error) {
				sr, ok := w.(*ErrResponse)
				require.True(t, ok)

				require.Equal(t, sr.HttpResponseCode, http.StatusInternalServerError)
				require.Equal(t, sr.StatusText, "Internal server error")
				require.Equal(t, sr.ErrorMessage, "An error occured. Please try again later")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, w render.Renderer) {
				sr, ok := w.(*ErrResponse)
				require.True(t, ok)

				require.Equal(t, recorder.Code, sr.HttpResponseCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := tt.setup(tt.err)
			tt.validate(sr, tt.err)

			router := chi.NewRouter()
			router.Get("/", func(w http.ResponseWriter, req *http.Request) {
				render.Render(w, req, sr)
			})

			httpRecorder := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			router.ServeHTTP(httpRecorder, req)

			tt.checkResponse(t, httpRecorder, sr)
		})
	}

}
