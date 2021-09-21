package counter

import (
	"testing"
	"time"
)

func TestEmptyCounter(t *testing.T) {
	RemoveDataFile()

	c := NewCounter(60 * time.Second)

	err := c.Start()
	if err != nil {
		t.Fatal(err)
	}

	err = c.Stop()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotEmptyStartStop(t *testing.T) {
	RemoveDataFile()

	func() {
		c := NewCounter(60 * time.Second)

		err := c.Start()
		if err != nil {
			t.Fatal(err)
		}

		n := c.OnRequest()
		if n != 1 {
			t.Error(n)
		}

		n = c.OnRequest()
		if n != 2 {
			t.Error(n)
		}

		n = c.OnRequest()
		if n != 3 {
			t.Error(n)
		}

		err = c.Stop()
		if err != nil {
			t.Fatal(err)
		}
	}()

	func() {
		c := NewCounter(60 * time.Second)

		err := c.Start()
		if err != nil {
			t.Fatal(err)
		}

		n := c.Size()
		if n != 3 {
			t.Error(n)
		}

		n = c.OnRequest()
		if n != 4 {
			t.Error(n)
		}

		err = c.Stop()
		if err != nil {
			t.Fatal(err)
		}
	}()
}

func TestEviction(t *testing.T) {
	RemoveDataFile()

	func() {
		c := NewCounter(1 * time.Second)

		err := c.Start()
		if err != nil {
			t.Fatal(err)
		}

		n := c.OnRequest()
		if n != 1 {
			t.Error(n)
		}

		n = c.Size()
		if n != 1 {
			t.Error(n)
		}

		time.Sleep(1500 * time.Millisecond)

		n = c.Size()
		if n != 0 {
			t.Error(n)
		}

		err = c.Stop()
		if err != nil {
			t.Fatal(err)
		}
	}()
}

func TestEvictionOnLoad(t *testing.T) {
	RemoveDataFile()

	func() {
		c := NewCounter(1 * time.Second)

		err := c.Start()
		if err != nil {
			t.Fatal(err)
		}

		n := c.OnRequest()
		if n != 1 {
			t.Error(n)
		}

		n = c.Size()
		if n != 1 {
			t.Error(n)
		}

		err = c.Stop()
		if err != nil {
			t.Fatal(err)
		}
	}()

	func() {
		time.Sleep(time.Second)

		c := NewCounter(1 * time.Second)

		err := c.Start()
		if err != nil {
			t.Fatal(err)
		}

		n := c.Size()
		if n != 0 {
			t.Error(n)
		}

		err = c.Stop()
		if err != nil {
			t.Fatal(err)
		}
	}()
}
