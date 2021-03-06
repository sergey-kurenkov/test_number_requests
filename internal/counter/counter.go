package counter

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"
)

type Counter struct {
	interval       time.Duration
	mt             sync.Mutex
	reqTimes       *list.List
	stopCounter    chan struct{}
	stoppedCounter chan struct{}
}

func NewCounter(interval time.Duration) *Counter {
	return &Counter{
		interval:       interval,
		reqTimes:       list.New(),
		stopCounter:    make(chan struct{}),
		stoppedCounter: make(chan struct{}),
	}
}

func (c *Counter) Start() error {
	if err := c.readOnStart(); err != nil {
		return err
	}

	go c.removeOldEntries()

	return nil
}

func (c *Counter) Stop() error {
	close(c.stopCounter)
	<-c.stoppedCounter

	return c.writeOnExit()
}

func (c *Counter) removeOldEntries() {
	const defaultCleaningInterval time.Duration = 100 * time.Millisecond

	tc := time.NewTicker(defaultCleaningInterval)

	for {
		select {
		case <-c.stopCounter:
			close(c.stoppedCounter)
			return
		case <-tc.C:
			func() {
				c.mt.Lock()
				defer c.mt.Unlock()

				c.removeOld()
			}()
		}
	}
}

func (c *Counter) Size() int64 {
	c.mt.Lock()
	defer c.mt.Unlock()

	return int64(c.reqTimes.Len())
}

func (c *Counter) OnRequest() int64 {
	c.mt.Lock()
	defer c.mt.Unlock()

	ts := time.Now()
	c.reqTimes.PushBack(ts)
	c.removeOld()

	return int64(c.reqTimes.Len())
}

func (c *Counter) removeOld() {
	for c.reqTimes.Len() > 0 {
		front := c.reqTimes.Front()

		frontTS := front.Value.(time.Time)
		if time.Since(frontTS) <= c.interval {
			return
		}

		c.reqTimes.Remove(front)

		continue
	}
}

const (
	counterFileName string      = "./counter.txt"
	defaultFilePerm os.FileMode = 0o644
)

func (c *Counter) readOnStart() error {
	file, err := os.OpenFile(counterFileName, os.O_RDONLY, defaultFilePerm)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	if err != nil {
		return err
	}

	defer file.Close()

	l := list.New()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		if txt == "" {
			continue
		}

		t, err := time.Parse(time.RFC3339Nano, txt)
		if err != nil {
			return err
		}

		if time.Since(t) > c.interval {
			continue
		}

		l.PushBack(t)
	}

	c.mt.Lock()
	defer c.mt.Unlock()
	c.reqTimes = l

	return nil
}

func (c *Counter) writeOnExit() error {
	c.mt.Lock()
	defer c.mt.Unlock()

	file, err := os.OpenFile(counterFileName, os.O_CREATE|os.O_WRONLY, defaultFilePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for e := c.reqTimes.Front(); e != nil; e = e.Next() {
		ts := e.Value.(time.Time)

		_, err := writer.WriteString(fmt.Sprintf("%s\n", ts.Format(time.RFC3339Nano)))
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

func RemoveDataFile() error {
	err := os.Remove(counterFileName)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	return err
}
