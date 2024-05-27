package workers

import (
	"container/heap"
	"fmt"
	"server/pkg/logger"
	"sync"
	"time"
)

type MinHeapEvent []NotificationEvent

func (h MinHeapEvent) Len() int { return len(h) }

func (h MinHeapEvent) Less(i, j int) bool {
	return h[i].UnixTimestamp < h[j].UnixTimestamp
}

func (h MinHeapEvent) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *MinHeapEvent) Push(x interface{}) {
	*h = append(*h, x.(NotificationEvent))
}

func (h *MinHeapEvent) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *MinHeapEvent) Remove(eventID uint) {
	n := len(*h)
	idx := -1
	for i := 0; i < n; i++ {
		if (*h)[i].EventID == eventID {
			idx = i
			break
		}
	}
	if idx == -1 {
		// Event not found
		return
	}
	heap.Remove(h, idx)
}

type Worker struct {
	id       int
	minHeap  *MinHeapEvent
	mu       sync.Mutex
	newValue chan NotificationEvent
	stop     chan struct{}
	callback func(event *NotificationEvent)
}

func NewWorker(id int, callback func(event *NotificationEvent)) *Worker {
	h := &MinHeapEvent{}
	heap.Init(h)
	return &Worker{
		id:       id,
		minHeap:  h,
		newValue: make(chan NotificationEvent),
		stop:     make(chan struct{}),
		callback: callback,
	}
}

func (w *Worker) AddEvent(event NotificationEvent) {
	w.newValue <- event
}

func (w *Worker) Stop() {
	close(w.stop)
}

func (w *Worker) run() {
	count := 0
	logger.Info("scheduler started", "worker", w.id)
	for {
		w.mu.Lock()
		for w.minHeap.Len() == 0 {
			w.mu.Unlock()
			select {
			case newEvent := <-w.newValue:
				logger.Info("new event received", "event", newEvent, "worker", w.id)
				w.mu.Lock()
				heap.Push(w.minHeap, newEvent)
				w.mu.Unlock()
			case <-w.stop:
				fmt.Printf("Worker %d stopping\n", w.id)
				return
			}
			w.mu.Lock()
		}
		minEvent := (*w.minHeap)[0]
		w.mu.Unlock()
		now := time.Now().Unix()
		waitTime := time.Duration(minEvent.UnixTimestamp-now) * time.Second
		if waitTime > 0 {
			select {
			case <-time.After(waitTime):
				w.mu.Lock()
				if w.minHeap.Len() > 0 && (*w.minHeap)[0].UnixTimestamp == minEvent.UnixTimestamp {
					heap.Pop(w.minHeap)
					w.callback(&minEvent)
					count++
				}
				w.mu.Unlock()
			case newEvent := <-w.newValue:
				logger.Info("new event received", "event", newEvent, "worker", w.id)
				w.mu.Lock()
				heap.Push(w.minHeap, newEvent)
				w.mu.Unlock()
			case <-w.stop:
				return
			}
		} else {
			w.mu.Lock()
			if w.minHeap.Len() > 0 && (*w.minHeap)[0].UnixTimestamp == minEvent.UnixTimestamp {
				heap.Pop(w.minHeap)
				w.callback(&minEvent)
				count++
			}
			w.mu.Unlock()
		}
	}
}

type EventDispatcher struct {
	workers []*Worker
	wg      sync.WaitGroup

	mu sync.Mutex

	current int
}

func NewEventDispatcher(numWorkers int, callback func(event *NotificationEvent)) *EventDispatcher {
	workers := make([]*Worker, 0)
	for i := 0; i < numWorkers; i++ {
		workers = append(workers, NewWorker(i, callback))
	}
	return &EventDispatcher{
		workers: workers,
	}
}

func (d *EventDispatcher) Start() {
	for _, worker := range d.workers {
		d.wg.Add(1)
		go func(w *Worker) {
			defer d.wg.Done()
			w.run()
		}(worker)
	}
}

func (d *EventDispatcher) AddEvent(event *NotificationEvent) {
	logger.Info("adding event", "event", event)

	d.mu.Lock()
	worker := d.workers[d.current]
	d.current = (d.current + 1) % len(d.workers)
	d.mu.Unlock()
	worker.AddEvent(*event)

	logger.Info("event added", "event", event)
}

func (d *EventDispatcher) Stop() {
	for _, worker := range d.workers {
		worker.Stop()
	}
	d.wg.Wait()
}
