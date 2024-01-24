package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/salimnassim/kahva"
)

func main() {
	transport := &http.Transport{}

	if os.Getenv("XMLRPC_URL") == "" {
		log.Fatal().Msgf("XMLRPC_URL environment variable is empty")
		return
	}

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
		log.Fatal().Err(err).Msgf("unable to create rtorrent client instance")
		return
	}
	defer rtorrent.Close()

	fs := http.FileServer(http.Dir("./www"))

	r := mux.NewRouter()
	r.Handle("/", fs)
	r.PathPrefix("/assets/").Handler(fs)

	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/view/{view}", kahva.ViewHandler(rtorrent))
	s.HandleFunc("/system", kahva.SystemHandler(rtorrent))
	s.HandleFunc("/load", kahva.LoadHandler(rtorrent)).Methods("POST")
	s.HandleFunc("/torrent/{hash}/{action}", kahva.TorrentHandler(rtorrent))
	s.Use(kahva.CORSMiddleware)

	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = "0.0.0.0:8080"
	}

	srv := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		Addr:              address,
		Handler:           r,
	}

	log.Info().Msgf("listen address: http://%s", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("cant serve")
		return
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
