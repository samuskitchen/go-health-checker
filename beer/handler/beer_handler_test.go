// Package handler_test contains unit tests for Beer HTTP handlers.
// Uses Echo to create request contexts, testify for assertions and mocks
// to simulate BeerService behaviors.
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_mocksService "github.com/samuskitchen/go-health-checker/beer/mocks/interfaces"
	"github.com/samuskitchen/go-health-checker/beer/model"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const (
	fakeBeerIdUint uint = 1 // Example beer ID
)

// HTTPContext encapsulates parameters to simulate HTTP requests in Echo.
type HTTPContext struct {
	Req         *http.Request              // Simulated request
	Res         *httptest.ResponseRecorder // Response recorder
	EchoContext echo.Context               // Built Echo context
}

// SetupHTTPContext prepares an Echo context with:
// - Method, route and body (if body != nil).
// - Path parameters (pathParams).
// - Query params (queryParams).
// - Content-Type header (mediaType).
func SetupHTTPContext(
	method, routePattern string,
	body interface{},
	pathParams map[string]string,
	queryParams map[string]string,
	mediaType string,
) HTTPContext {
	// 1) Serialize body if it exists
	var reader io.Reader
	if body != nil {
		payload, _ := json.Marshal(body)
		reader = bytes.NewReader(payload)
	}

	// 2) Create Echo server, request and response recorder
	e := echo.New()
	req := httptest.NewRequest(method, routePattern, reader)
	rec := httptest.NewRecorder()
	// Set Content-Type if specified
	if mediaType != "" {
		req.Header.Set(echo.HeaderContentType, mediaType)
	}

	// 3) Add query params to the URL
	q := req.URL.Query()
	for k, v := range queryParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// 4) Create Echo context
	ctx := e.NewContext(req, rec)

	// 5) Assign path parameters at once
	if len(pathParams) > 0 {
		keys, vals := make([]string, 0, len(pathParams)), make([]string, 0, len(pathParams))
		for k, v := range pathParams {
			keys = append(keys, k)
			vals = append(vals, v)
		}
		ctx.SetParamNames(keys...)
		ctx.SetParamValues(vals...)
	}

	return HTTPContext{Req: req, Res: rec, EchoContext: ctx}
}

// Test_beerHandler_GetAllBeersHandler tests the endpoint to retrieve all beers.
func Test_beerHandler_GetAllBeersHandler(t *testing.T) {
	mockService := _mocksService.NewMockBeerService(t)
	handler := NewBeerHandler(mockService)
	ctx := context.Background()

	t.Run("Successful Response", func(t *testing.T) {
		beers := []model.BeersResponse{
			{
				ID:        fakeBeerIdUint,
				Name:      "Gulden Draak",
				Brewery:   "Blót",
				Country:   "BE",
				Price:     6.50,
				Currency:  "EUR",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/beers",
			nil,
			nil,
			nil,
			"",
		)
		mockService.On("GetAllBeers", ctx).Return(beers, nil).Once()

		res := httpContext.Res
		err := handler.GetAllBeersHandler(httpContext.EchoContext)

		expectedResponse := `[
				{
					"id": 1,
					"name": "Gulden Draak",
					"brewery": "Blót",
					"country": "BE",
					"price": 6.50,
					"currency": "EUR",
					"created_at": "2023-01-01T00:00:00Z",
					"updated_at": "2023-01-01T00:00:00Z"
				}
		]`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		httpContext := SetupHTTPContext(
			http.MethodGet,
			"/beers",
			nil,
			nil,
			nil,
			"",
		)
		mockService.On("GetAllBeers", ctx).Return(nil, assert.AnError).Once()

		res := httpContext.Res
		err := handler.GetAllBeersHandler(httpContext.EchoContext)

		expectedResponse := `{"message": "assert.AnError general error for testing"}`

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.JSONEq(t, expectedResponse, res.Body.String())
		mockService.AssertExpectations(t)
	})
}
