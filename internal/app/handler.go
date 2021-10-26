package app

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sergey-kurenkov/test_number_requests/internal/rate_limiter"

	"github.com/sergey-kurenkov/test_number_requests/internal/counter"
)

type Application struct {
	counter     *counter.Counter
	rateLimiter *rate_limiter.RateLimiter
	stopped     int64
}

func NewGetNumberRequestsApp(duration time.Duration) *Application {
	const defaultCapacity = 5

	app := &Application{
		counter:     counter.NewCounter(duration),
		rateLimiter: rate_limiter.NewRateLimiter(defaultCapacity),
	}

	if err := app.counter.Start(); err != nil {
		log.Fatal(err)
	}

	return app
}

func (app *Application) Stop() {
	atomic.AddInt64(&app.stopped, 1)

	if err := app.counter.Stop(); err != nil {
		log.Fatal(err)
	}
}

func (app *Application) Handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/get-number-requests", app.handleGetNumberRequests)

	return r
}

func (app *Application) handleGetNumberRequests(w http.ResponseWriter, r *http.Request) {
	if atomic.LoadInt64(&app.stopped) == 1 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	const defaultSleep = 2 * time.Second

	app.rateLimiter.AddRequest()
	defer app.rateLimiter.OnFinishRequest()

	number := app.counter.OnRequest()

	time.Sleep(defaultSleep)

	w.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprintf(w, "%d\n", number); err != nil {
		log.Println(err)
	}
}
