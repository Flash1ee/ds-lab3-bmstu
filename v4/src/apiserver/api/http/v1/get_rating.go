package v1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"ds-lab3-bmstu/apiserver/services/rating"
)

type RatingRequest struct {
	AuthedRequest `valid:"optional"`
}

type RatingResponse struct {
	Rating
}

func (a *api) GetRating(c echo.Context, req RatingRequest) error {
	data, err := a.core.GetUserRating(c.Request().Context(), req.Username)
	if err != nil {
		if errors.Is(err, rating.ErrUnavaliable) {
			return c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Bonus Service unavailable"})
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := RatingResponse{
		Rating: Rating(data),
	}

	return c.JSON(http.StatusOK, &resp)
}
