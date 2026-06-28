package scrapper

type LikeEntry struct {
	ID       int64  `json:"id,string"`
	Title    string `json:"title"`
	URL      string `json:"webpage_url"`
	FilePath string
	ProxyURL string
}

type metadata struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Uploader   string `json:"uploader"`
	Thumbnails []struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"thumbnails"`
}

type Track struct {
	ID       int64
	Title    string
	Artist   string
	CoverURL string
	URL      string
}
