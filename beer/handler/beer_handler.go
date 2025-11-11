// Package handler defines HTTP controllers (handlers) for the Beer entity.
// They are responsible for receiving requests, validating them and delegating business
// logic to the corresponding service. Also handles responses and errors.
package handler

import (
	"net/http"

	"github.com/samuskitchen/go-health-checker/beer/interfaces"

	// Models are used only in Swagger annotations, hence the blank import
	_ "github.com/samuskitchen/go-health-checker/beer/model"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// errorResponse structure for generic error responses
type errorResponse struct {
	Message string `json:"message"`
}

// beerHandler implements BeerHandler and encapsulates the Beer service.
type beerHandler struct {
	beerService interfaces.BeerService
}

// BeerHandler groups handler methods for Beer endpoints.
type BeerHandler interface {
	GetAllBeersHandler(c echo.Context) error // List all beers
}

// NewBeerHandler builds a BeerHandler with the service implementation.
func NewBeerHandler(service interfaces.BeerService) BeerHandler {
	return &beerHandler{beerService: service}
}

// GetAllBeersHandler retrieves all beers from the database and returns JSON.
// @Description Get all beers
// @Tags Beer
// @ID GetAllBeersHandler
// @Success 200 {array} model.BeersResponse
// @Failure 500 {object} errorResponse
// @Router /beers [GET]
func (bh *beerHandler) GetAllBeersHandler(c echo.Context) error {
	ctx := c.Request().Context()

	beers, err := bh.beerService.GetAllBeers(ctx)
	if err != nil {
		log.Error().Msgf("error GetAllBeers: %v", err)
		return c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, beers)
}
