package retry

import (
	"log/slog"
	"slices"
	"sync"
	"time"
)

type Client[T any] struct {
	queue []T
	mx    sync.RWMutex
	log   *slog.Logger
	start sync.Once
}

func New[T any](log *slog.Logger) (*Client[T], error) {
	return &Client[T]{log: log}, nil
}

func (c *Client[T]) Append(v T) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.queue = append(c.queue, v)

	return nil
}

func (c *Client[T]) Start(f func(T) error) error {
	try := func() {
		c.mx.Lock()

		if len(c.queue) == 0 {
			c.mx.Unlock()

			return
		}

		i := c.queue[0]
		c.queue = slices.Delete(c.queue, 0, 1)

		c.mx.Unlock()

		if err := f(i); err != nil {
			if err := c.Append(i); err != nil {
				c.log.Error("failed to append to queue", "err", err)
			}
		}
		c.log.Info("success retry")
	}

	c.start.Do(func() {
		go func() {
			ticker := time.NewTicker(time.Second)
			for {
				<-ticker.C

				try()
			}
		}()
	})

	return nil
}
