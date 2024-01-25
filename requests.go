package kahva

type TorrentPriorityRequest struct {
	Priority int `json:"priority"`
}

type ThrottleRequest struct {
	Type      string `json:"type"`
	Kilobytes int    `json:"kilobytes"`
}
