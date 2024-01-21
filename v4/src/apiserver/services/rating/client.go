package rating

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"ds-lab3-bmstu/apiserver/core/ports/library"
	"ds-lab3-bmstu/pkg/circuit_breaker"
	"ds-lab3-bmstu/pkg/readiness/httpprober"
	v1 "ds-lab3-bmstu/rating/api/http/v1"

	"ds-lab3-bmstu/apiserver/core/ports/rating"
	"ds-lab3-bmstu/pkg/readiness"
)

const probeKey = "http-rating-client"

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
	ErrUnavaliable       = library.ErrUnavaliable
)

type Client struct {
	lg   *slog.Logger
	conn *resty.Client
	cb   circuit_breaker.CircuitBreaker
}

func New(lg *slog.Logger, cfg rating.Config, probe *readiness.Probe) (*Client, error) {
	client := resty.New().
		SetTransport(&http.Transport{
			MaxIdleConns:       10,               //nolint: gomnd
			IdleConnTimeout:    30 * time.Second, //nolint: gomnd
			DisableCompression: true,
		}).
		SetBaseURL(fmt.Sprintf("http://%s", net.JoinHostPort(cfg.Host, cfg.Port)))

	c := Client{
		lg:   lg,
		conn: client,
		cb: circuit_breaker.New(circuit_breaker.Settings{
			Name:                          "rating_cb",
			MaxErrorsFromHalfToCloseState: uint32(cfg.MaxErrorsTrying),
			TimeoutFromOpenToHalfState:    time.Second * 5,
			ClearCountsInCloseState:       time.Minute,
			FailureRequestsToOpenState:    1,
		}, lg),
	}

	go httpprober.New(lg, client).Ping(probeKey, probe)

	return &c, nil
}

func (c *Client) UpdateUserRating(
	_ context.Context, username string, diff int,
) error {
	resp, err := c.conn.R().
		SetHeader("X-User-Name", username).
		SetQueryParam("diff", strconv.Itoa(diff)).
		SetResult(&v1.RatingResponse{}).
		Patch("/api/v1/rating")
	if err != nil {
		return fmt.Errorf("failed to execute http request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%d: %w", resp.StatusCode(), ErrInvalidStatusCode)
	}

	return nil
}
