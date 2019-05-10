package restapi

import (
	"fmt"
	"net/http"

	"github.com/sofiukl/go-queue/worker"
)

// WorkQueue - This is the work queues
var WorkQueue = make(chan worker.Work, 100)

// ReceiveWork - This method receives the task requests
func ReceiveWork(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}
	work := worker.Work{Name: name}
	WorkQueue <- work
	fmt.Println("Work request queued")

	w.WriteHeader(http.StatusCreated)
	return
}
