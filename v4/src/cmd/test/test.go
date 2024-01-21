package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

// Get wraps http.Get in CircuitBreaker.
func Get(url string, i int) ([]byte, error) {
	body, err := cb.Execute(func() (interface{}, error) {
		_, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if i > 7 {
			return nil, nil
		}
		return nil, errors.New("demo error")
		//defer resp.Body.Close()
		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	return nil, err
		//}
		//
		//return body, nil
	})
	if err != nil {
		return nil, err
	}

	return body.([]byte), nil
}

func main() {
	for x := 0; x < 10; x++ {
		body, err := Get("http://www.google.com/robots.txt", x)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(body))
	}
}
