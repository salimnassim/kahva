package kahva

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

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
			respond(ErrorResponse{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}

		respond(ViewResponse{
			Status:   "ok",
			Torrents: torrents,
		}, http.StatusOK, w)
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
			respond(ErrorResponse{
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

func ThrottleHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req ThrottleRequest
		err := decoder.Decode(&req)
		if err != nil {
			log.Error().Err(err).Msg("cant decode throttle request json")
			respond(ErrorResponse{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}

		if req.Type != "up" && req.Type != "down" {
			respond(ErrorResponse{
				Status:  "error",
				Message: "type must be up or down",
			}, http.StatusBadRequest, w)
			return
		}

		if req.Type == "up" {
			err := rt.GlobalThrottleUp(req.Kilobytes)
			if err != nil {
				log.Error().Err(err).Msg("cant set global up throttle")
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}
		}

		if req.Type == "down" {
			err := rt.GlobalThrottleDown(req.Kilobytes)
			if err != nil {
				log.Error().Err(err).Msg("cant set global down throttle")
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}
		}

		respond(Response{
			Status:  "ok",
			Message: "",
		}, http.StatusOK, w)
	}
}

func LoadHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)

		file, _, err := r.FormFile("file")
		if err != nil {
			log.Error().Err(err).Msg("cant read file form")
			respond(ErrorResponse{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}
		defer file.Close()

		buffer := bytes.NewBuffer(nil)
		_, err = io.Copy(buffer, file)
		if err != nil {
			log.Error().Err(err).Msg("cant copy file to buffer")
			respond(ErrorResponse{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusBadRequest, w)
			return
		}

		err = rt.LoadRawStart(buffer.Bytes())
		if err != nil {
			log.Error().Err(err).Msg("xmlrpc load raw start failed")
			respond(ErrorResponse{
				Status:  "error",
				Message: err.Error(),
			}, http.StatusInternalServerError, w)
			return
		}
		respond(Response{
			Status: "ok",
		}, http.StatusOK, w)
	}
}

func TorrentHandler(rt *Rtorrent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if vars["action"] == "stop" {
			err := rt.Stop(vars["hash"])
			if err != nil {
				log.Printf("error in action stop handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusInternalServerError, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "pause" {
			err := rt.Pause(vars["hash"])
			if err != nil {
				log.Printf("error in action stop handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusInternalServerError, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "resume" {
			err := rt.Resume(vars["hash"])
			if err != nil {
				log.Printf("error in action stop handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusInternalServerError, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "start" {
			err := rt.Start(vars["hash"])
			if err != nil {
				log.Printf("error in action start handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "hash" {
			err := rt.CheckHash(vars["hash"])
			if err != nil {
				log.Printf("error in action start handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}
			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "files" {
			args := []interface{}{vars["hash"], "",
				"f.path=", "f.size_bytes=", "f.size_chunks=",
				"f.completed_chunks=", "f.frozen_path=", "f.priority=",
				"f.is_created=", "f.is_open="}

			files, err := rt.FMulticall(args)
			if err != nil {
				log.Printf("error in action files handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(FilesResponse{
				Status: "ok",
				Files:  files,
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "peers" {
			args := []interface{}{vars["hash"], "",
				"p.id=", "p.address=", "p.port=",
				"p.banned=", "p.client_version=", "p.completed_percent=",
				"p.is_encrypted=", "p.is_incoming=", "p.is_obfuscated=",
				"p.peer_rate=", "p.peer_total=", "p.up_rate=", "p.up_total="}

			peers, err := rt.PMulticall(args)
			if err != nil {
				log.Printf("error in action peers handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(PeersResponse{
				Status: "ok",
				Peers:  peers,
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "trackers" {
			args := []interface{}{vars["hash"], "",
				"t.id=", "t.type=", "t.url=",
				"t.activity_time_last=", "t.activity_time_next=", "t.can_scrape=",
				"t.is_usable=", "t.is_enabled=", "t.failed_counter=",
				"t.failed_time_last=", "t.failed_time_next=", "t.is_busy=",
				"t.is_open=",
			}

			trackers, err := rt.TMulticall(args)
			if err != nil {
				log.Printf("error in action trackers handler: %s", err)
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(TrackersResponse{
				Status:   "ok",
				Trackers: trackers,
			}, http.StatusOK, w)
			return
		}

		if vars["action"] == "priority" {
			decoder := json.NewDecoder(r.Body)
			var req TorrentPriorityRequest
			err := decoder.Decode(&req)
			if err != nil {
				log.Error().Err(err).Msg("unable to decode priority request")
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusOK, w)
				return
			}

			err = rt.Priority(vars["hash"], req.Priority)
			if err != nil {
				log.Error().Err(err).Msg("unable to set torrent priority")
				respond(ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				}, http.StatusBadRequest, w)
				return
			}

			respond(Response{
				Status: "ok",
			}, http.StatusOK, w)
			return
		}

	}
}
