package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

// MinHeap represents a min-heap of integers
type MinHeap []int64

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int64))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Worker handles incoming times and waits until the minimum time is reached
type Worker struct {
	id       int
	minHeap  *MinHeap
	mu       sync.Mutex
	newValue chan int64
	stop     chan struct{}
}

func NewWorker(id int) *Worker {
	h := &MinHeap{}
	heap.Init(h)
	return &Worker{
		id:       id,
		minHeap:  h,
		newValue: make(chan int64),
		stop:     make(chan struct{}),
	}
}

func (w *Worker) AddTime(unixTime int64) {
	w.newValue <- unixTime
}

func (w *Worker) Stop() {
	close(w.stop)
}

func (w *Worker) run() {
	count := 0

	for {
		w.mu.Lock()
		for w.minHeap.Len() == 0 {
			w.mu.Unlock()
			select {
			case newTime := <-w.newValue:
				w.mu.Lock()
				heap.Push(w.minHeap, newTime)
				w.mu.Unlock()
			case <-w.stop:
				fmt.Printf("Worker %d stopping\n", w.id)
				return
			}
			w.mu.Lock()
		}
		minTime := (*w.minHeap)[0]
		w.mu.Unlock()

		now := time.Now().Unix()
		waitTime := time.Duration(minTime-now) * time.Second
		if waitTime > 0 {
			select {
			case <-time.After(waitTime):
				w.mu.Lock()
				if w.minHeap.Len() > 0 && (*w.minHeap)[0] == minTime {
					heap.Pop(w.minHeap)
					count++
					fmt.Printf("Worker %d: Time reached: %v\n", w.id, minTime)
				}
				w.mu.Unlock()
			case newTime := <-w.newValue:
				w.mu.Lock()
				heap.Push(w.minHeap, newTime)
				w.mu.Unlock()
			case <-w.stop:
				fmt.Printf("Worker %d stopping\n", w.id)
				return
			}
		} else {
			w.mu.Lock()
			if w.minHeap.Len() > 0 && (*w.minHeap)[0] == minTime {
				heap.Pop(w.minHeap)
				fmt.Printf("Worker %d: Time reached: %v\n", w.id, minTime)
			}
			w.mu.Unlock()
		}
	}
}

type Dispatcher struct {
	workers []*Worker
	current int
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func NewDispatcher(numWorkers int) *Dispatcher {
	workers := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = NewWorker(i)
	}
	return &Dispatcher{
		workers: workers,
	}
}

func (d *Dispatcher) Start() {
	for _, worker := range d.workers {
		d.wg.Add(1)
		go func(w *Worker) {
			defer d.wg.Done()
			w.run()
		}(worker)
	}
}

func (d *Dispatcher) AddTime(unixTime int64) {
	d.mu.Lock()
	worker := d.workers[d.current]
	d.current = (d.current + 1) % len(d.workers)
	d.mu.Unlock()
	worker.AddTime(unixTime)
}

func (d *Dispatcher) Stop() {
	for _, worker := range d.workers {
		worker.Stop()
	}
	d.wg.Wait() // Wait for all workers to finish
}

func main() {
	numWorkers := 4 // Number of workers

	ticker := time.NewTicker(time.Second)

	go func() {
		count := 0
		for {
			<-ticker.C
			count = count + 1
			fmt.Println("this is count", count)
			// worker.AddTime(time.Now().Unix())
		}
	}()

	dispatcher := NewDispatcher(numWorkers)

	dispatcher.Start()

	// Example usage
	for i := 0; i < 100; i++ {
		dispatcher.AddTime(time.Now().Add(time.Duration(5) * time.Second).Unix())
	}

	time.Sleep(1 * time.Minute)
	dispatcher.Stop()
}
