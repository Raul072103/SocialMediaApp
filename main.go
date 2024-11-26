package main

import "net/http"

type api struct {
	addr string
}

func (s *api) GetIndexPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index page"))
}

func (s *api) GetUsersPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("users page"))
}

func main() {
	api := &api{addr: ":8080"}

	// Initialize the ServeMux
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("/", api.GetIndexPage)
	mux.HandleFunc("/users", api.GetUsersPage)

	err := srv.ListenAndServe()

	if err != nil {
	}
}
