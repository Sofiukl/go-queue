package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var waitgroup sync.WaitGroup
var workQueue chan Work
var workerQueue chan chan Work

const (
	NWorkers = 4
	NTasks   = 10
)

func main() {
	fmt.Printf("Starting Worker..\nNo of Workers %d, No of Tasks %d\n", NWorkers, NTasks)
	workQueue = make(chan Work, 10)
	ds := Dispatcher{1, "work dispatcher"}
	ds.dispatchWork()
	allocateTasks()
	waitgroup.Wait()
	fmt.Printf("%d tasks are completed by %d Workers", NTasks, NWorkers)
}

type Work struct {
	id   int
	name string
}

type Worker struct {
	id            int
	workerChannel chan Work
	workerQueue   chan chan Work
}

func newWorker(id int, workerQueue chan chan Work) Worker {
	return Worker{id: id, workerChannel: make(chan Work), workerQueue: workerQueue}
}

func (worker *Worker) start() {

	go func() {
		for {
			worker.workerQueue <- worker.workerChannel
			select {
			case work := <-worker.workerChannel:

				fmt.Printf("worker %d: working on %s!\n", worker.id, work.name)
				time.Sleep(1000 * time.Millisecond)
				fmt.Printf("worker%d: %s completed\n", worker.id, work.name)
				waitgroup.Done()
			}
		}

	}()
}

type Dispatcher struct{
	id int
	name string
}
func (d *Dispatcher) dispatchWork() {
	workerQueue := make(chan chan Work, NWorkers)
	for i := 0; i < NWorkers; i++ {
		worker := newWorker(i+1, workerQueue)
		worker.start()
	}
	go func() {
		for {
			select {
			case work := <-workQueue:
				workerChannel := <-workerQueue
				workerChannel <- work

			}
		}

	}()
}

func allocateTasks() {
	for i := 0; i < NTasks; i++ {
		workName := "work" + strconv.Itoa(i+1)
		w1 := Work{id: i + 1, name: workName}
		waitgroup.Add(1)
		workQueue <- w1
	}
}
