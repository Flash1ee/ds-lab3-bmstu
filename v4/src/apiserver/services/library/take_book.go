package library

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ds-lab3-bmstu/apiserver/core/ports/library"
	v1 "ds-lab3-bmstu/library/api/http/v1"
)

func (c *Client) ReturnBook(
	_ context.Context, libraryID string, bookID string,
) (library.Book, error) {
	body, err := json.Marshal(v1.TakeBookRequest{
		BookID:    bookID,
		LibraryID: libraryID,
	})
	if err != nil {
		return library.Book{}, fmt.Errorf("failed to format json body: %w", err)
	}

	resp, err := c.conn.R().
		SetHeader("Content-Type", "application/json").
		SetPathParam("lib_id", libraryID).
		SetPathParam("book_id", bookID).
		SetBody(body).
		SetResult(&v1.ReturnBookResponse{}).
		Post("/api/v1/libraries/{lib_id}/books/{book_id}/return")
	if err != nil {
		return library.Book{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.Book{}, fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	data, _ := resp.Result().(*v1.ReturnBookResponse)

	return library.Book(data.Book), nil
}
