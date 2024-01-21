package http

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"ds-lab2-bmstu/pkg/httpvalidator"
	"ds-lab2-bmstu/pkg/readiness"
	"ds-lab2-bmstu/rating/api/http/common"
	v1 "ds-lab2-bmstu/rating/api/http/v1"
)

type Core interface {
	v1.Core
}

type Server struct {
	mx *echo.Echo
}

func New(lg *slog.Logger, probe *readiness.Probe, core Core) (*Server, error) {
	mx := echo.New()
	mx.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.RequestID(),
	)

	mx.Debug = false
	mx.HideBanner = true
	mx.HidePort = true
	mx.HTTPErrorHandler = func(err error, c echo.Context) {
		mx.DefaultHTTPErrorHandler(err, c)
	}
	mx.Validator = &httpvalidator.CustomValidator{}

	s := Server{mx: mx}

	err := common.InitListener(s.mx, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init common apis: %w", err)
	}

	err = v1.InitListener(s.mx, lg.With("api", "v1"), core)
	if err != nil {
		return nil, fmt.Errorf("failed to init v1 apis: %w", err)
	}

	return &s, nil
}

func (s *Server) ListenAndServe(addr string) error {
	return s.mx.Start(addr) //nolint: wrapcheck
}
