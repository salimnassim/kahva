package kahva

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ViewResponse struct {
	Status   string    `json:"status"`
	Torrents []Torrent `json:"torrents"`
}

type SystemResponse struct {
	Status string `json:"status"`
	System System `json:"system"`
}

func respond(anything any, statusCode int, w http.ResponseWriter) {
	bytes, err := json.Marshal(anything)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(bytes)
}

func ViewHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		args := []interface{}{
			"", vars["view"],
			"d.hash=", "d.name=",
			"d.size_bytes=", "d.completed_bytes=", "d.up.rate=",
			"d.up.total=", "d.down.rate=", "d.down.total=",
			"d.message=", "d.is_active=", "d.is_open=",
			"d.is_hash_checking=", "d.peers_accounted=", "d.peers_complete=",
			"d.state=", "d.state_changed=", "d.state_counter=", "d.priority=",
			"d.custom1=", "d.custom2=", "d.custom3=",
			"d.custom4=", "d.custom5="}

		torrents, err := rt.DMulticall("main", args)
		if err != nil {
			log.Error().Err(err).Msgf("cant fetch view")
			respond(Response{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}

		respond(
			ViewResponse{
				Status:   "ok",
				Torrents: torrents,
			},
			http.StatusOK,
			w,
		)
	}
}

func SystemHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		args := []interface{}{
			[]interface{}{
				SystemCall{
					MethodName: "system.hostname",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "system.pid",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "system.time_seconds",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "system.api_version",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "system.client_version",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "system.library_version",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "throttle.global_down.total",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "throttle.global_up.total",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "throttle.global_down.rate",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "throttle.global_up.rate",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "throttle.global_down.max_rate",
					Params:     []string{""},
				},
				SystemCall{
					MethodName: "throttle.global_up.max_rate",
					Params:     []string{""},
				},
			},
		}

		result, err := rt.SystemMulticall(args)
		if err != nil {
			log.Error().Err(err).Msgf("cant fetch system")
			respond(Response{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusInternalServerError, w)
			return
		}
		respond(SystemResponse{
			Status: "ok",
			System: result,
		}, http.StatusOK, w)
	}
}
