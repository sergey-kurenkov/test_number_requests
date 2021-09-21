package http_server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/sergey-kurenkov/test_number_requests/internal/counter"
	"time"
)

type Application struct {
	counter *counter.Counter
}

func NewGetNumberRequestsApp(duration time.Duration) *Application {
	app := &Application{
		counter: counter.NewCounter(duration),
	}

	if err:= app.counter.Start(); err != nil {
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
	number := app.counter.OnRequest()
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "%d\n", number); err != nil {
		log.Println(err)
	}
}
