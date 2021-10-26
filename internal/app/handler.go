package app

import (
	"fmt"
	"github.com/sergey-kurenkov/test_number_requests/internal/rate_limiters"
	"log"
	"net/http"
	"time"


	"github.com/sergey-kurenkov/test_number_requests/internal/counter"
)

type Application struct {
	counter     *counter.Counter
	rateLimiters *rate_limiters.RateLimiters
}

func NewGetNumberRequestsApp(duration time.Duration) *Application {
	app := &Application{
		counter:     counter.NewCounter(duration),
		rateLimiters: rate_limiters.NewRateLimiters(),
	}

	if err := app.counter.Start(); err != nil {
		log.Fatal(err)
	}

	return app
}

func (app *Application) Stop() {
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
	address := getAddress(r)
	rateLimiter := app.rateLimiters.GetRateLimiter(address)

	rateLimiter.AddRequest()
	defer rateLimiter.OnFinishRequest()

	number := app.counter.OnRequest()

	time.Sleep(2 * time.Second)

	w.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprintf(w, "%d\n", number); err != nil {
		log.Println(err)
	}
}

func getAddress(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
