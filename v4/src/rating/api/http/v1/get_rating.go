package v1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"ds-lab3-bmstu/rating/core/ports/ratings"
)

type RatingRequest struct {
	AuthedRequest `valid:"optional"`
}

type RatingResponse struct {
	Stars uint64 `json:"stars" valid:"range(0|100),required"`
}

func (a *api) GetRating(c echo.Context, req RatingRequest) error {
	data, err := a.core.GetUserRating(c.Request().Context(), req.Username)
	isEmptyUser := errors.Is(err, ratings.ErrNotFound)
	if err != nil && !isEmptyUser {
		return c.NoContent(http.StatusInternalServerError)
	}
	if isEmptyUser {
		data.Stars = 101
	}

	resp := RatingResponse{
		Stars: uint64(data.Stars),
	}

	return c.JSON(http.StatusOK, &resp)
}
