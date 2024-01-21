package library

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"

	"ds-lab3-bmstu/apiserver/core/ports/library"
	"ds-lab3-bmstu/pkg/circuit_breaker"
	"ds-lab3-bmstu/pkg/readiness"
	"ds-lab3-bmstu/pkg/readiness/httpprober"
)

const probeKey = "http-library-client"

var ErrInvalidStatusCode = errors.New("invalid status code")

type Client struct {
	lg     *slog.Logger
	conn   *resty.Client
	cbLib  circuit_breaker.CircuitBreaker
	cbBook circuit_breaker.CircuitBreaker
}

func New(lg *slog.Logger, cfg library.Config, probe *readiness.Probe) (*Client, error) {
	client := resty.New().
		SetTransport(&http.Transport{
			MaxIdleConns:       10,               //nolint: gomnd
			IdleConnTimeout:    30 * time.Second, //nolint: gomnd
			DisableCompression: true,
		}).
		SetBaseURL(fmt.Sprintf("http://%s", net.JoinHostPort(cfg.Host, cfg.Port))).EnableTrace()

	//cbLib := gobreaker.NewCircuitBreaker(gobreaker.Settings{
	//	Name: "library_cb",
	//	// кол-во попыток успешных запросов для перехода из полуоткрытого состояния в закрытое
	//	MaxRequests: uint32(cfg.MaxErrorsTrying),
	//	// время, через которое осуществляется переход из открытого в полуоткрытое состояние
	//	Timeout: time.Second * 5,
	//	// время, через которое очищаются счетчики в закрытом состоянии
	//	Interval: time.Minute,
	//	// проверка, если true, состояние меняется с закрытого на открытое
	//	ReadyToTrip: func(counts gobreaker.Counts) bool {
	//		return counts.ConsecutiveSuccesses >= 5
	//	},
	//	OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
	//		lg.Info(fmt.Sprintf("[library] circuitBreaker change from %s to %s", from.String(), to.String()))
	//	},
	//})
	c := Client{
		lg:   lg,
		conn: client,
		cbLib: circuit_breaker.New(circuit_breaker.Settings{
			Name:                          "library_cb",
			MaxErrorsFromHalfToCloseState: uint32(cfg.MaxErrorsTrying),
			TimeoutFromOpenToHalfState:    time.Second * 5,
			ClearCountsInCloseState:       time.Minute,
			FailureRequestsToOpenState:    1,
		}, lg),
		cbBook: circuit_breaker.New(circuit_breaker.Settings{
			Name:                          "library_book_cb",
			MaxErrorsFromHalfToCloseState: uint32(cfg.MaxErrorsTrying),
			TimeoutFromOpenToHalfState:    time.Second * 5,
			ClearCountsInCloseState:       time.Minute,
			FailureRequestsToOpenState:    1,
		}, lg),
	}

	go httpprober.New(lg, client).Ping(probeKey, probe)

	return &c, nil
}
