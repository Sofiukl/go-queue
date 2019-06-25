package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// GithubAPIBaseURL - Github API base url
var (
	GithubAPIBaseURL = "https://api.github.com/users/"
)

// GithubRepo - This is Github Repo
// ID - Id of github repo
// Name - Name of github repo
type GithubRepo struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

var sem chan int
var waitgroup sync.WaitGroup

func main() {
	users := []string{"sofiukl", "adityakeyal", "rajibdas008", "libra"}
	sem = make(chan int, 3) // semaphore
	waitgroup.Add(len(users))
	for i := 0; i < len(users); i++ {
		go func(user string, i int) {
			sem <- i
			fmt.Printf("Fetching repos for user [%s] \n", user)
			time.Sleep(1000 * time.Millisecond)
			url := GithubAPIBaseURL + user + "/repos"
			resp, err := http.Get(url)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			var repos []GithubRepo
			json.Unmarshal(body, &repos)
			for i, repo := range repos {
				fmt.Printf("%d: %s/%v\n", i, user, repo.Name)
			}
			fmt.Printf("Fetching repos for user [%s] is completed\n", user)
			<-sem
			waitgroup.Done()
		}(users[i], i)

	}
	waitgroup.Wait()
}
