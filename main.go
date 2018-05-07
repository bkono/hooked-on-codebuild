package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bkono/way"
)

var (
	port = flag.String("port", "8080", "port to listen on")
)

func handleBitbucket(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.UserAgent(), "Bitbucket-Webhooks") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if key := r.Header.Get("X-Event-Key"); !(key == "repo:push") {
		log.Println("Bitbucket web handled called for something other than repo:push", key)
		return
	}

	var push RepoPushHook
	err := json.NewDecoder(r.Body).Decode(&push)
	if err != nil {
		log.Println("err decoding json body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	owner := push.Repository.Owner.Username
	repo := push.Repository.Name
	log.Println("processing", owner, repo)

	prjNames := make(map[string]int)
	for _, c := range push.Push.Changes {
		log.Println("new", c.New.Name)
		prjNames[fmt.Sprintf("%s-%s-%s", owner, repo, c.New.Name)] += 1
	}

	log.Println("prjNames is ready", prjNames)

	// here is where I go fetch the projects and see which ones I can run
}

func main() {
	router := way.NewRouter()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "This is not the page you are looking for")
	})

	router.HandleFunc("POST", "/bitbucket", handleBitbucket)

	// log.Fatalln(gateway.ListenAndServe(":8080", router))
	log.Fatalln(http.ListenAndServe(":"+*port, router))
}

type RepoPushHook struct {
	Actor struct {
		Type        string `json:"type"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
	} `json:"actor"`
	Repository struct {
		Type  string `json:"type"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		UUID     string `json:"uuid"`
		FullName string `json:"full_name"`
		Name     string `json:"name"`
		Website  string `json:"website"`
		Owner    struct {
			Type        string `json:"type"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			UUID        string `json:"uuid"`
			Links       struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
		} `json:"owner"`
		Scm       string `json:"scm"`
		IsPrivate bool   `json:"is_private"`
	} `json:"repository"`
	Push struct {
		Changes []struct {
			New struct {
				Type   string `json:"type"`
				Name   string `json:"name"`
				Target struct {
					Type    string    `json:"type"`
					Hash    string    `json:"hash"`
					Message string    `json:"message"`
					Date    time.Time `json:"date"`
					Parents []struct {
						Type  string `json:"type"`
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
					} `json:"parents"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
				} `json:"target"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Commits struct {
						Href string `json:"href"`
					} `json:"commits"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
				} `json:"links"`
			} `json:"new"`
			Old struct {
				Type   string `json:"type"`
				Name   string `json:"name"`
				Target struct {
					Type    string    `json:"type"`
					Hash    string    `json:"hash"`
					Message string    `json:"message"`
					Date    time.Time `json:"date"`
					Parents []struct {
						Type  string `json:"type"`
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
					} `json:"parents"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
				} `json:"target"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Commits struct {
						Href string `json:"href"`
					} `json:"commits"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
				} `json:"links"`
			} `json:"old"`
			Links struct {
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Diff struct {
					Href string `json:"href"`
				} `json:"diff"`
				Commits struct {
					Href string `json:"href"`
				} `json:"commits"`
			} `json:"links"`
			Created bool `json:"created"`
			Forced  bool `json:"forced"`
			Closed  bool `json:"closed"`
			Commits []struct {
				Hash    string `json:"hash"`
				Type    string `json:"type"`
				Message string `json:"message"`
				Links   struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
				} `json:"links"`
			} `json:"commits"`
			Truncated bool `json:"truncated"`
		} `json:"changes"`
	} `json:"push"`
}
