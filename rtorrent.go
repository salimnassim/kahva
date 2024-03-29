package kahva

import (
	"encoding/base64"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/kolo/xmlrpc"
)

type Torrent struct {
	Hash           string `rt:"d.hash=" json:"hash"`
	Name           string `rt:"d.name=" json:"name"`
	SizeBytes      int64  `rt:"d.size_bytes=" json:"size_bytes"`
	CompletedBytes int64  `rt:"d.completed_bytes=" json:"completed_bytes"`
	UploadRate     int64  `rt:"d.up.rate=" json:"upload_rate"`
	UploadTotal    int64  `rt:"d.up.total=" json:"upload_total"`
	DownloadRate   int64  `rt:"d.down.rate=" json:"download_rate"`
	DownloadTotal  int64  `rt:"d.down.total=" json:"download_total"`
	Message        string `rt:"d.message=" json:"message"`
	BaseFilename   string `rt:"d.base_filename=" json:"base_filename"`
	BasePath       string `rt:"d.base_path=" json:"base_path"`
	IsActive       int64  `rt:"d.is_active=" json:"is_active"`
	IsOpen         int64  `rt:"d.is_open=" json:"is_open"`
	IsHashing      int64  `rt:"d.is_hash_checking=" json:"is_hashing"`
	Leechers       int64  `rt:"d.peers_accounted=" json:"leechers"`
	Seeders        int64  `rt:"d.peers_complete=" json:"seeders"`
	State          int64  `rt:"d.state=" json:"state"`
	StateChanged   int64  `rt:"d.state_changed=" json:"state_changed"`
	StateCounter   int64  `rt:"d.state_counter=" json:"state_counter"`
	Priority       int64  `rt:"d.priority=" json:"priority"`
	Custom1        string `rt:"d.custom1=" json:"custom1"`
	Custom2        string `rt:"d.custom2=" json:"custom2"`
	Custom3        string `rt:"d.custom3=" json:"custom3"`
	Custom4        string `rt:"d.custom4=" json:"custom4"`
	Custom5        string `rt:"d.custom5=" json:"custom5"`
}

type File struct {
	Path            string `rt:"f.path=" json:"path"`
	Size            int64  `rt:"f.size_bytes=" json:"size"`
	SizeChunks      int64  `rt:"f.size_chunks=" json:"size_chunks"`
	CompletedChunks int64  `rt:"f.completed_chunks=" json:"completed_chunks"`
	FrozenPath      string `rt:"f.frozen_path=" json:"frozen_path"`
	Priority        int64  `rt:"f.priority=" json:"priority"`
	IsCreated       int64  `rt:"f.is_created=" json:"is_created"`
	IsOpen          int64  `rt:"f.is_open=" json:"is_open"`
}

type Peer struct {
	ID               string `rt:"p.id=" json:"id"`
	Address          string `rt:"p.address=" json:"address"`
	Port             int64  `rt:"p.port=" json:"port"`
	Banned           int64  `rt:"p.banned=" json:"banned"`
	ClientVersion    string `rt:"p.client_version=" json:"client_version"`
	CompletedPercent int64  `rt:"p.completed_percent=" json:"completed_percent"`
	IsEncrypted      int64  `rt:"p.is_encrypted=" json:"is_encrypted"`
	IsIncoming       int64  `rt:"p.is_incoming=" json:"is_incoming"`
	IsObfuscated     int64  `rt:"p.is_obfuscated=" json:"is_obfuscated"`
	DownloadRate     int64  `rt:"p.peer_rate=" json:"down_rate"`
	DownloadTotal    int64  `rt:"p.peer_total=" json:"down_total"`
	UploadRate       int64  `rt:"p.up_rate=" json:"up_rate"`
	UploadTotal      int64  `rt:"p.up_total=" json:"up_total"`
}

type Tracker struct {
	ID               string `rt:"t.id=" json:"tracker_id"`
	ActivityTimeLast int64  `rt:"t.activity_time_last=" json:"activity_time_last"`
	ActivityTimeNext int64  `rt:"t.activity_time_next=" json:"activity_time_next"`
	CanScrape        int64  `rt:"t.can_scrape=" json:"can_scrape"`
	IsUsable         int64  `rt:"t.is_usable=" json:"t.is_usable"`
	IsEnabled        int64  `rt:"t.is_enabled=" json:"is_enabled"`
	FailedCounter    int64  `rt:"t.failed_counter=" json:"failed_counter"`
	FailedTimeLast   int64  `rt:"t.failed_time_last=" json:"failed_time_last"`
	FailedTimeNext   int64  `rt:"t.failed_time_next=" json:"failed_time_next"`
	LatestEvent      string `rt:"t.latest_event=" json:"latest_event"`
	IsBusy           int64  `rt:"t.is_busy=" json:"is_busy"`
	IsOpen           int64  `rt:"t.is_open=" json:"is_open"`
	Type             int64  `rt:"t.type=" json:"type"`
	URL              string `rt:"t.url=" json:"url"`
}

type System struct {
	APIVersion     string `rt:"system.api_version" json:"api_version"`
	ClientVersion  string `rt:"system.client_version" json:"client_version"`
	LibraryVersion string `rt:"system.library_version" json:"library_version"`

	Hostname string `rt:"system.hostname" json:"hostname"`
	PID      int64  `rt:"system.pid" json:"pid"`
	Time     int64  `rt:"system.time_seconds" json:"time_seconds"`

	ThrottleGlobalDownTotal   int64 `rt:"throttle.global_down.total" json:"throttle_global_down_total"`
	ThrottleGlobalUpTotal     int64 `rt:"throttle.global_up.total" json:"throttle_global_up_total"`
	ThrottleGlobalDownRate    int64 `rt:"throttle.global_down.rate" json:"throttle_global_down_rate"`
	ThrottleGlobalUpRate      int64 `rt:"throttle.global_up.rate" json:"throttle_global_up_rate"`
	ThrottleGlobalDownMaxRate int64 `rt:"throttle.global_down.max_rate" json:"throttle_global_down_max_rate"`
	ThrottleGlobalUpMaxRate   int64 `rt:"throttle.global_up.max_rate" json:"throttle_global_up_max_rate"`
}

type SystemCall struct {
	MethodName string      `xmlrpc:"methodName" json:"method_name"`
	Params     interface{} `xmlrpc:"params" json:"params"`
}

type Config struct {
	URL       string
	Transport http.RoundTripper
}

type Rtorrent struct {
	client *xmlrpc.Client
}

// Creates a new instance of Rtorrent client
func NewRtorrent(config Config) (*Rtorrent, error) {
	xmlrpcClient, err := xmlrpc.NewClient(config.URL, config.Transport)
	if err != nil {
		return nil, err
	}

	rtorrent := &Rtorrent{
		client: xmlrpcClient,
	}
	return rtorrent, nil
}

// Closes the underlying XMLRPC client.
func (rt *Rtorrent) Close() error {
	err := rt.client.Close()
	if err != nil {
		return err
	}
	return nil
}

// Lists available XMLRPC methods
func (rt *Rtorrent) ListMethods() ([]string, error) {
	var result []string
	err := rt.client.Call("system.listMethods", nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Load and start a torrent
func (rt *Rtorrent) LoadRawStart(file []byte) error {
	base64 := base64.StdEncoding.EncodeToString(file)

	err := rt.client.Call("load.raw_start_verbose", []interface{}{"", xmlrpc.Base64(base64)}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Stop torrent with the specified hash
func (rt *Rtorrent) Stop(hash string) error {
	err := rt.client.Call("d.stop", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

// Start torrent with the specified hash
func (rt *Rtorrent) Start(hash string) error {
	err := rt.client.Call("d.start", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

// Pause torrent with the specified hash
func (rt *Rtorrent) Pause(hash string) error {
	err := rt.client.Call("d.pause", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

// Pause torrent with the specified hash
func (rt *Rtorrent) Resume(hash string) error {
	err := rt.client.Call("d.resume", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

// Pause torrent with the specified hash
func (rt *Rtorrent) CheckHash(hash string) error {
	err := rt.client.Call("d.check_hash", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

// Erase torrent with the specified hash
func (rt *Rtorrent) Erase(hash string) error {
	err := rt.client.Call("d.erase", hash, nil)
	if err != nil {
		return err
	}
	return nil
}

// Set torrent priority
func (rt *Rtorrent) Priority(hash string, priority int) error {
	if priority < 0 || priority > 3 {
		return errors.New("invalid priority")
	}

	err := rt.client.Call("d.priority.set", []interface{}{hash, priority}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Set global down throttle.
func (rt *Rtorrent) GlobalThrottleDown(kilobytes int) error {
	kb := strconv.Itoa(kilobytes)
	err := rt.client.Call("throttle.global_down.max_rate.set_kb", []interface{}{"", kb}, nil)
	if err != nil {
		return err
	}
	return nil
}

// Set global up throttle.
func (rt *Rtorrent) GlobalThrottleUp(kilobytes int) error {
	kb := strconv.Itoa(kilobytes)
	err := rt.client.Call("throttle.global_up.max_rate.set_kb", []interface{}{"", kb}, nil)
	if err != nil {
		return err
	}
	return nil
}

// View multicall.
func (rt *Rtorrent) DMulticall(target string, args interface{}) ([]Torrent, error) {
	var result interface{}
	err := rt.client.Call("d.multicall2", args, &result)
	if err != nil {
		return nil, err
	}

	torrents := multicallTags[Torrent](result, args)
	return torrents, nil
}

// File multicall.
func (rt *Rtorrent) FMulticall(args interface{}) ([]File, error) {
	var result interface{}
	err := rt.client.Call("f.multicall", args, &result)
	if err != nil {
		return nil, err
	}

	files := multicallTags[File](result, args)
	return files, nil
}

// Peer multicall.
func (rt *Rtorrent) PMulticall(args interface{}) ([]Peer, error) {
	var result interface{}
	err := rt.client.Call("p.multicall", args, &result)
	if err != nil {
		return nil, err
	}

	peers := multicallTags[Peer](result, args)
	return peers, nil
}

// Torrent multicall.
func (rt *Rtorrent) TMulticall(args interface{}) ([]Tracker, error) {
	var result interface{}
	err := rt.client.Call("t.multicall", args, &result)
	if err != nil {
		return nil, err
	}

	trackers := multicallTags[Tracker](result, args)
	return trackers, nil
}

// System multicall.
func (rt *Rtorrent) SystemMulticall(args interface{}) (System, error) {
	var result interface{}
	err := rt.client.Call("system.multicall", args, &result)
	if err != nil {
		return System{}, err
	}

	system := systemTags(result, args)
	return system, nil
}

// Maps XMLRPC result to a struct using fields from args with reflection
func multicallTags[T File | Torrent | Peer | Tracker](result interface{}, args interface{}) []T {
	items := make([]T, 0)
	for _, outer := range result.([]interface{}) {
		item := new(T)
		for idx := 2; idx < len(args.([]interface{})); idx++ {
			ref := outer.([]interface{})[idx-2]
			fname := args.([]interface{})[idx].(string)
			vo := reflect.ValueOf(item)
			el := vo.Elem()
			for i := 0; i < el.NumField(); i++ {
				field := el.Type().Field(i)
				if fname == field.Tag.Get("rt") {
					if ref == nil {
						continue
					}
					el.Field(i).Set(reflect.ValueOf(ref))
				}
			}

		}
		items = append(items, *item)
	}
	return items
}

func systemTags(result interface{}, args interface{}) System {
	system := &System{}
	r := result.([]interface{})
	a := args.([]interface{})
	for idx := 0; idx < len(r); idx++ {
		ref := r[idx].([]interface{})[0]
		fname := a[0].([]interface{})[idx].(SystemCall).MethodName
		vo := reflect.ValueOf(system)
		el := vo.Elem()
		for i := 0; i < el.NumField(); i++ {
			field := el.Type().Field(i)
			if fname == field.Tag.Get("rt") {
				if ref == nil {
					continue
				}
				el.Field(i).Set(reflect.ValueOf(ref))
			}
		}
	}
	return *system
}
