package handlers

import (
	"runtime"
	"sort"
	"time"

	"github.com/containerum/events-api/pkg/util/ticker"
	"github.com/containerum/kube-client/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}

func withWS(ctx *gin.Context, limit int, getfuncs ...eventsFunc) {
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Debug(err)
		return
	}
	defer c.Close()
	//Abort context in the end
	defer ctx.Abort()

	var errChan = make(chan error)
	var resultChan = make(chan model.Event)
	var finalChan = make(chan []model.Event)

	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	//Checking for closed connection
	go func() {
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
		}
	}()

	//Limiter. Waiting for all DB request to finish on first run.
	//	var doner = wg.NewWG(len(getfuncs))
	for _, eventSource := range getfuncs {
		go func(es eventsFunc) {
			EventAgregator{
				Ctx:         ctx,
				Limit:       limit,
				EventSource: es,
				EventDrain:  resultChan,
				ErrChan:     errChan,
			}.Run()

		}(eventSource)
	}

	go EventBatcher{
		Ctx:               ctx,
		ErrChan:           errChan,
		Quant:             1 * time.Second,
		BatchDrain:        finalChan,
		PreallocBatchSize: 16,
		EventSource:       resultChan,
	}.Run()

	go func() {
		pingTicker := ticker.NewTicker(1 * time.Second)
		pingTicker.Start()
		defer pingTicker.Stop()
		for {
			select {
			case result, ok := <-finalChan:
				if !ok {
					return
				}
				if len(result) > 0 {
					sort.Slice(result, func(i, j int) bool {
						timei, _ := time.Parse(time.RFC3339, result[i].Time)
						timej, _ := time.Parse(time.RFC3339, result[j].Time)
						return timei.Before(timej)
					})
					logrus.Infof("Writing %v events", len(result))
					if err := c.WriteJSON(result); err != nil {
						errChan <- err
					}
				}
			case <-pingTicker.Ticks():
				if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()

	for {
		select {
		case selecerr := <-errChan:
			if selecerr != nil {
				logrus.Debug(selecerr)
				return
			}
		}
	}
}

type EventBatcher struct {
	Ctx         *gin.Context
	ErrChan     chan error
	Quant       time.Duration
	EventSource <-chan model.Event
	BatchDrain  chan<- []model.Event
	//	FirsWaveNotify    <-chan struct{}
	PreallocBatchSize int
}

func (batcher EventBatcher) Run() {
	var ctx = batcher.Ctx
	var aborted = AbortWaiter(ctx.IsAborted)
	var finalChan = batcher.BatchDrain
	var resultChan = batcher.EventSource

	results := make([]model.Event, 0)
	var timer = ticker.NewTicker(batcher.Quant)
	timer.Start()
	defer timer.Stop()
	for {
		select {
		case <-aborted:
			return
		case <-ctx.Done():
			return
		case <-timer.Ticks():
			finalChan <- results
			results = make([]model.Event, 0)
		case event, ok := <-resultChan:
			if !ok {
				return
			}
			results = append(results, event)
		}
	}
}

type EventAgregator struct {
	Ctx         *gin.Context
	StartAt     time.Time
	Limit       int
	EventSource eventsFunc
	EventDrain  chan<- model.Event
	ErrChan     chan<- error
}

func (agregate EventAgregator) Run() {
	var ctx = agregate.Ctx
	var getfunc = agregate.EventSource
	var funcLimit = agregate.Limit
	var startTime = agregate.StartAt
	var resultChan = agregate.EventDrain
	var errChan = agregate.ErrChan
	var aborted = AbortWaiter(agregate.Ctx.IsAborted)
	var firstEventSend = false

	defer close(agregate.EventDrain)

	for {
		select {
		case <-aborted: //If context aborted finish goroutine
			return
		case <-agregate.Ctx.Done():
			return
		default:
			resp, err := getfunc(ctx.Params, funcLimit, startTime)
			if err != nil {
				errChan <- err
				return
			}
		batchLoop:
			for _, event := range resp.Events {
				select {
				case <-aborted: //If context aborted finish goroutine
					return
				case <-agregate.Ctx.Done():
					return
				case resultChan <- event:
					continue batchLoop
				}
			}
			if !firstEventSend {
				firstEventSend = true
			}
			funcLimit = 0
			//Get only new events
			startTime = time.Now()
			time.Sleep(60 * time.Second)
		}
	}
}

func AbortWaiter(aborted func() bool) <-chan struct{} {
	var wait = make(chan struct{})
	go func() {
		defer close(wait)
		for {
			if aborted() {
				return
			}
			runtime.Gosched()
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return wait
}
