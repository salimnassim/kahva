package kahva

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrorResponse struct {
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

type FilesResponse struct {
	Status string `json:"status"`
	Files  []File `json:"files"`
}

type PeersResponse struct {
	Status string `json:"status"`
	Peers  []Peer `json:"peers"`
}

type TrackersResponse struct {
	Status   string    `json:"status"`
	Trackers []Tracker `json:"trackers"`
}
