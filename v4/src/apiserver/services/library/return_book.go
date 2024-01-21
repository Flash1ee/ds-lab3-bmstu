package library

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ds-lab3-bmstu/apiserver/core/ports/library"
	v1 "ds-lab3-bmstu/library/api/http/v1"
)

func (c *Client) TakeBook(ctx context.Context, libraryID string, bookID string) (library.ReservedBook, error) {
	body, err := json.Marshal(v1.TakeBookRequest{
		BookID:    bookID,
		LibraryID: libraryID,
	})
	if err != nil {
		return library.ReservedBook{}, fmt.Errorf("failed to format json body: %w", err)
	}

	resp, err := c.conn.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&v1.TakeBookResponse{}).
		Post("/api/v1/books")
	if err != nil {
		return library.ReservedBook{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.ReservedBook{}, fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	data, _ := resp.Result().(*v1.TakeBookResponse)

	return library.ReservedBook{
		Book:    library.Book(data.Book),
		Library: library.Info(data.Library),
	}, nil
}
