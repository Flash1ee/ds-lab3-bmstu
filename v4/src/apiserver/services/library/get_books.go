package library

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/sony/gobreaker"

	"ds-lab3-bmstu/apiserver/core/ports/library"
	v1 "ds-lab3-bmstu/library/api/http/v1"
)

func (c *Client) GetBooks(
	ctx context.Context, libraryID string, showAll bool, page uint64, size uint64,
) (library.Books, error) {
	data, err := c.cbBook.Execute(func() (any, error) {
		return c.getBooks(ctx, libraryID, showAll, page, size)
	})
	if err == nil {
		res, ok := data.(library.Books)
		if !ok {
			return library.Books{}, nil
		}

		return res, nil
	}
	if errors.Is(err, gobreaker.ErrOpenState) {
		return library.Books{}, nil
	}

	return library.Books{}, fmt.Errorf("get books: %w", err)
}

func (c *Client) getBooks(
	ctx context.Context, libraryID string, showAll bool, page uint64, size uint64,
) (library.Books, error) {
	if size == 0 {
		size = math.MaxUint64
	}

	q := map[string]string{
		"size": strconv.FormatUint(size, 10),
		"page": strconv.FormatUint(page, 10),
	}

	if showAll {
		q["show_all"] = "1"
	}

	resp, err := c.conn.R().
		SetContext(ctx).
		SetQueryParams(q).
		SetPathParam("library_id", libraryID).
		SetResult(&v1.BooksResponse{}).
		Get("/api/v1/libraries/{library_id}/books")
	if err != nil {
		return library.Books{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.Books{}, fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	data, _ := resp.Result().(*v1.BooksResponse)

	books := library.Books{Total: data.Total}
	for _, book := range data.Items {
		books.Items = append(books.Items, library.Book(book))
	}

	return books, nil
}

func (c *Client) GetBooksByIDs(
	ctx context.Context, ids []string,
) (library.Books, error) {
	data, err := c.cbBook.Execute(func() (any, error) {
		return c.getBooksByIDs(ctx, ids)
	})
	if err == nil {
		res, ok := data.(library.Books)
		if !ok {
			return library.Books{}, nil
		}

		return res, nil
	}
	if errors.Is(err, gobreaker.ErrOpenState) {
		return library.Books{}, nil
	}

	return library.Books{}, fmt.Errorf("get books by id: %w", err)
}

func (c *Client) getBooksByIDs(
	_ context.Context, ids []string,
) (library.Books, error) {
	id, err := json.Marshal(ids)
	if err != nil {
		return library.Books{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := c.conn.R().
		SetQueryParam("ids", string(id)).
		SetResult(&v1.BooksResponse{}).
		Get("/api/v1/books")
	if err != nil {
		return library.Books{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.Books{}, fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	data, _ := resp.Result().(*v1.BooksResponse)

	books := library.Books{Total: data.Total}
	for _, book := range data.Items {
		books.Items = append(books.Items, library.Book(book))
	}

	return books, nil
}
