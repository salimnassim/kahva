package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/salimnassim/kahva"
)

func main() {
	transport := &http.Transport{}

	// enable basic auth if env is set
	if os.Getenv("XMLRPC_USERNAME") != "" && os.Getenv("XMLRPC_PASSWORD") != "" {
		transport.RegisterProtocol("https",
			newBasicAuthTransport(
				os.Getenv("XMLRPC_USERNAME"),
				os.Getenv("XMLRPC_PASSWORD"),
			),
		)
	}

	rtorrent, err := kahva.NewRtorrent(
		kahva.Config{
			URL:       os.Getenv("XMLRPC_URL"),
			Transport: transport,
		},
	)

	if err != nil {
		log.Fatalf("unable to create rtorrent client instance: %v", err)
		return
	}

	defer rtorrent.Close()

	r := mux.NewRouter()
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/view/{view}", kahva.ViewHandler(rtorrent))
	s.HandleFunc("/system", kahva.SystemHandler(rtorrent))
	// s.HandleFunc("/load", LoadHandler(rtorrent)).Methods("POST")
	// s.HandleFunc("/methods", MethodsHandler(rtorrent))
	// s.HandleFunc("/torrent/{hash}/{action}", TorrentHandler(rtorrent))
	s.Use(kahva.CORSMiddleware)

	srv := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		Addr:              os.Getenv("SERVER_ADDRESS"),
		Handler:           r,
	}

	log.Printf("listen address: http://%s", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failure: %s", err)
	}

}

type basicAuthTransport struct {
	Username string
	Password string
}

func newBasicAuthTransport(username, password string) *basicAuthTransport {
	return &basicAuthTransport{
		Username: username,
		Password: password,
	}
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultTransport.RoundTrip(req)
}
