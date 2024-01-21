package library

import (
	"fmt"
	"log/slog"

	"ds-lab2-bmstu/pkg/apiutils"
	"ds-lab2-bmstu/pkg/readiness"

	"ds-lab2-bmstu/library/api/http"
	"ds-lab2-bmstu/library/config"
	"ds-lab2-bmstu/library/core"
	"ds-lab2-bmstu/library/services/librarydb"
)

type App struct {
	cfg *config.Config

	http *http.Server
}

func New(lg *slog.Logger, cfg *config.Config) (*App, error) {
	a := App{cfg: cfg}

	probe := readiness.New()

	libraries, err := librarydb.New(lg.With("module", "library"), cfg.Libraries, probe)
	if err != nil {
		return nil, fmt.Errorf("failed to init libraries db: %w", err)
	}

	core, err := core.New(lg.With("module", "core"), probe, libraries)
	if err != nil {
		return nil, fmt.Errorf("failed to init core: %w", err)
	}

	a.http, err = http.New(lg.With("module", "http_api"), probe, core)
	if err != nil {
		return nil, fmt.Errorf("failed to init http server: %w", err)
	}

	return &a, nil
}

func (s *App) Run(lg *slog.Logger) {
	apiutils.Serve(lg,
		apiutils.NewCallable(s.cfg.HTTPAddr, s.http),
	)
}
