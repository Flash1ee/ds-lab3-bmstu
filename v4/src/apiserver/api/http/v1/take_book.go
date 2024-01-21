package v1

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sony/gobreaker"

	"ds-lab3-bmstu/apiserver/core"
	"ds-lab3-bmstu/apiserver/services/rating"
)

type TakeBookRequest struct {
	AuthedRequest `valid:"optional"`
	BookID        string `json:"bookUid" valid:"uuidv4,required"`
	LibraryID     string `json:"libraryUid" valid:"uuidv4,required"`
	End           Time   `json:"tillDate" valid:"required"`
}

type TakeBookResponse struct {
	ID      string    `json:"reservationUid"`
	Status  string    `json:"status"`
	Start   time.Time `json:"-"`
	End     time.Time `json:"-"`
	Book    Book      `json:"book"`
	Library Library   `json:"library"`
	Rating  Rating    `json:"rating"`
}

func (r TakeBookResponse) MarshalJSON() ([]byte, error) {
	type Alias TakeBookResponse

	return json.Marshal(&struct { //nolint: wrapcheck
		Alias
		Start string `json:"startDate"`
		End   string `json:"tillDate"`
	}{
		Alias: (Alias)(r),
		Start: r.Start.Format(time.DateOnly),
		End:   r.End.Format(time.DateOnly),
	})
}

func (a *api) TakeBook(c echo.Context, req TakeBookRequest) error {
	data, err := a.core.TakeBook(
		c.Request().Context(), req.Username, req.LibraryID, req.BookID, req.End.Time,
	)
	if err != nil {
		if errors.Is(err, rating.ErrUnavaliable) || errors.Is(err, gobreaker.ErrOpenState) {
			return c.JSON(http.StatusServiceUnavailable, ErrorResponse{Message: "Bonus Service unavailable"})
		}
		status := http.StatusInternalServerError
		if errors.Is(err, core.ErrInsufficientRating) {
			status = http.StatusPreconditionFailed
		}

		return c.NoContent(status)
	}

	resp := TakeBookResponse{
		ID:      data.ID,
		Status:  data.Status,
		Start:   data.Start,
		End:     data.End,
		Book:    Book(data.ReservedBook.Book),
		Library: Library(data.ReservedBook.Library),
		Rating:  Rating(data.Rating),
	}

	return c.JSON(http.StatusOK, &resp)
}
