package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sergey-kurenkov/test_number_requests/internal/counter"
)

func TestNumberRequests(t *testing.T) {
	if err := counter.RemoveDataFile(); err != nil {
		t.Fatal(err)
	}

	func() {
		app := NewGetNumberRequestsApp(60 * time.Second)
		defer app.Stop()

		srv := httptest.NewServer(app.Handler())
		defer srv.Close()

		res, err := http.Get(fmt.Sprintf("%s/get-number-requests", srv.URL)) // nolint:noctx
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

		res, err := http.Get(fmt.Sprintf("%s/get-number-requests", srv.URL)) // nolint:noctx
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
