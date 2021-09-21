package http_server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"github.com/sergey-kurenkov/test_number_requests/internal/counter"
	"testing"
	"time"
)

func TestNumberRequests(t *testing.T) {
	counter.RemoveDataFile()

	func() {
		app := NewGetNumberRequestsApp(60 * time.Second)
		defer app.Stop()

		srv := httptest.NewServer(app.Handler())
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/get-number-requests", srv.URL))
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("status not OK")
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(body) != "1\n" {
			t.Fail()
		}
	}()

	func() {
		app := NewGetNumberRequestsApp(60 * time.Second)
		defer app.Stop()

		srv := httptest.NewServer(app.Handler())
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/get-number-requests", srv.URL))
		if err != nil {
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("status not OK")
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if string(body) != "2\n" {
			t.Fail()
		}
	}()
}
