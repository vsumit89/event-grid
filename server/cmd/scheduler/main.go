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
	minHeap  *MinHeap
	mu       sync.Mutex
	newValue chan int64
	stop     chan struct{}
}

func NewWorker() *Worker {
	h := &MinHeap{}
	heap.Init(h)
	return &Worker{
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

					fmt.Printf("Time reached: %v\n", minTime)
				}
				w.mu.Unlock()
			case newTime := <-w.newValue:
				w.mu.Lock()
				heap.Push(w.minHeap, newTime)
				w.mu.Unlock()
			case <-w.stop:
				return
			}
		} else {
			w.mu.Lock()
			if w.minHeap.Len() > 0 && (*w.minHeap)[0] == minTime {
				heap.Pop(w.minHeap)
				fmt.Printf("Time reached: %v\n", minTime)
			}
			w.mu.Unlock()
		}
	}
}

func main() {
	worker := NewWorker()

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

	go worker.run()

	// Example usage
	worker.AddTime(time.Now().Add(12 * time.Second).Unix())
	worker.AddTime(time.Now().Add(5 * time.Second).Unix())
	worker.AddTime(time.Now().Add(3 * time.Second).Unix())

	time.Sleep(20 * time.Second)
	worker.Stop()
}
