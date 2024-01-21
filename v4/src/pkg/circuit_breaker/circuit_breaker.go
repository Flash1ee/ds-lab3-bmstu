package circuit_breaker

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sony/gobreaker"
)

type Settings struct {
	Name string

	MaxErrorsFromHalfToCloseState uint32        // Кол-во успешных запросов для перехода из полуоткрытого состояния в закрытое
	TimeoutFromOpenToHalfState    time.Duration // Время, через которое осуществляется переход из открытого в полуоткрытое состояние
	ClearCountsInCloseState       time.Duration // Время, через которое очищаются счетчики в закрытом состоянии
	FailureRequestsToOpenState    uint32        // Кол-во ошибок для перехода в открытое состояние
}

type CircuitBreaker interface {
	Execute(req func() (interface{}, error)) (interface{}, error)
}

func New(st Settings, lg *slog.Logger) CircuitBreaker {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        st.Name,
		MaxRequests: st.MaxErrorsFromHalfToCloseState,
		Timeout:     st.TimeoutFromOpenToHalfState,
		Interval:    st.ClearCountsInCloseState,
		// проверка, если true, состояние меняется с закрытого на открытое
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			lg.Info(fmt.Sprintf("[%s] failures requests = %d", st.Name, counts.ConsecutiveFailures))
			return counts.ConsecutiveFailures >= st.FailureRequestsToOpenState
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			lg.Info(fmt.Sprintf("[%s] circuitBreaker change from %s to %s", name, from.String(), to.String()))
		},
	})
	return cb
}
