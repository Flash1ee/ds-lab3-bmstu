package rating

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"ds-lab3-bmstu/apiserver/core/ports/rating"
	v1 "ds-lab3-bmstu/rating/api/http/v1"
)

func (c *Client) GetUserRating(ctx context.Context, username string) (rating.Rating, error) {
	data, err := c.cb.Execute(func() (interface{}, error) {
		return c.getUserRating(ctx, username)
	})
	if err == nil {
		res, ok := data.(rating.Rating)
		if !ok {
			return rating.Rating{}, nil
		}

		return res, nil
	}
	return rating.Rating{}, fmt.Errorf("get rating error: %w", err)
}

func (c *Client) getUserRating(ctx context.Context, username string) (rating.Rating, error) {
	resp, err := c.conn.R().
		SetContext(ctx).
		SetHeader("X-User-Name", username).
		SetResult(&v1.RatingResponse{}).
		Get("/api/v1/rating")
	if err != nil {
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) {
			err = ErrUnavaliable
		}

		return rating.Rating{}, fmt.Errorf("execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return rating.Rating{}, fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	data, ok := resp.Result().(*v1.RatingResponse)
	if !ok {
		return rating.Rating{}, errors.New("parse rating response error")
	}

	return rating.Rating(*data), nil
}
