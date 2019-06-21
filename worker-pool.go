package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var waitgroup sync.WaitGroup

// It's a work queue. All the new works are enqueued to this queue.
var workQueue chan Work

// workerQueue - Its a buffered channel of WorkChannel. Its capacity should be equal to no of workers.
// It is shared between different workers. It keeps track of available worker's workChannel.
// When new work came in workQueue channel, dispatcher dequeue one work from workQueue
// and enqueue to available worker's workChannel
// Worker enqueue work from its workChannel and process it.
var workerQueue chan chan Work

// NoWorkers - This denotes the no of workers
// NoTasks - This denotes the total no of tasks
const (
	NoWorkers = 4
	NoTasks   = 10
)

func main() {
	fmt.Printf("Starting Worker..\nNo of Workers %d, No of Tasks %d\n", NoWorkers, NoTasks)
	workQueue = make(chan Work, 10)
	ds := Dispatcher{1, "work dispatcher"}
	ds.dispatchWork()
	allocateTasks()
	waitgroup.Wait()
	fmt.Printf("%d tasks are completed by %d Workers", NoTasks, NoWorkers)
}

// Work struct
type Work struct {
	id   int
	name string
}

// Worker - This is contract for Worker
// Worker always listen to its workerChannel. when a new work is enqueqed to the
// a worker channel, then worker dequeue it and strats t process. After work is done then
// again workChannel will be available for accepting further task.
// workerChannel - Its a unbuffered channel.
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

// Dispatcher - This is the contract for task dispatcher
// Dispatcher dequeues from work queue and enqueue to the workChannel of available
// worker, then worker process takes the work and process it.
type Dispatcher struct {
	id   int
	name string
}

func (d *Dispatcher) dispatchWork() {
	workerQueue := make(chan chan Work, NoWorkers)
	for i := 0; i < NoWorkers; i++ {
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
	for i := 0; i < NoTasks; i++ {
		workName := "work" + strconv.Itoa(i+1)
		w1 := Work{id: i + 1, name: workName}
		waitgroup.Add(1)
		workQueue <- w1
	}
}
