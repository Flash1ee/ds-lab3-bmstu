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

func (c *Client) GetLibraries(
	ctx context.Context, city string, page uint64, size uint64,
) (library.Infos, error) {
	data, err := c.cbLib.Execute(func() (any, error) {
		return c.getLibraries(ctx, city, page, size)
	})
	if err == nil {
		res, ok := data.(library.Infos)
		if !ok {
			return library.Infos{}, nil
		}

		return res, nil
	}

	if errors.Is(err, gobreaker.ErrOpenState) {
		return library.Infos{}, nil
	}
	return library.Infos{}, fmt.Errorf("get libraries: %w", err)
}

func (c *Client) getLibraries(
	_ context.Context, city string, page uint64, size uint64,
) (library.Infos, error) {
	q := map[string]string{
		"city": city,
		"page": strconv.FormatUint(page, 10),
	}

	if size == 0 {
		size = math.MaxUint64
	}

	q["size"] = strconv.FormatUint(size, 10)

	resp, err := c.conn.R().
		SetQueryParams(q).
		SetResult(&v1.LibrariesResponse{}).
		Get("/api/v1/libraries")
	if err != nil {
		return library.Infos{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.Infos{}, fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	data, _ := resp.Result().(*v1.LibrariesResponse)

	libraries := library.Infos{Total: data.Total}
	for _, book := range data.Items {
		libraries.Items = append(libraries.Items, library.Info(book))
	}

	return libraries, nil
}

func (c *Client) GetLibrariesByIDs(ctx context.Context, ids []string) (library.Infos, error) {
	data, err := c.cbLib.Execute(func() (any, error) {
		return c.getLibrariesByIDs(ctx, ids)
	})
	if err == nil {
		res, ok := data.(library.Infos)
		if !ok {
			return library.Infos{}, nil
		}

		return res, nil
	}

	if errors.Is(err, gobreaker.ErrOpenState) {
		return library.Infos{}, nil
	}

	return library.Infos{}, fmt.Errorf("get libraries: %w", err)

}

func (c *Client) getLibrariesByIDs(ctx context.Context, ids []string) (library.Infos, error) {
	id, err := json.Marshal(ids)
	if err != nil {
		return library.Infos{}, fmt.Errorf("failed to marshal data: %w", err)
	}

	resp, err := c.conn.R().
		SetContext(ctx).
		SetQueryParam("ids", string(id)).
		SetResult(&v1.LibrariesResponse{}).
		Get("/api/v1/libraries")
	if err != nil {
		return library.Infos{}, fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return library.Infos{}, fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	data, _ := resp.Result().(*v1.LibrariesResponse)

	libraries := library.Infos{Total: data.Total}
	for _, book := range data.Items {
		libraries.Items = append(libraries.Items, library.Info(book))
	}

	return libraries, nil
}
