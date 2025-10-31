package delete

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
		status    int
	}{
		{
			name:   "Success",
			alias:  "test_alias",
			status: http.StatusOK,
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "missing alias",
			status:    http.StatusOK,
		},
		{
			name:      "DeleteURL Error",
			alias:     "test_alias",
			respError: "failed to delete url",
			mockError: errors.New("unexpected error"),
			status:    http.StatusOK,
		},
		{
			name:      "Not Found Error",
			alias:     "non_existent",
			respError: "failed to delete url",
			mockError: errors.New("url not found"),
			status:    http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewMockURLDeleter(t)

			if tc.alias != "" && tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			handler := New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, "/"+tc.alias, nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.status, rr.Code)

			body := rr.Body.String()

			var resp response.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
