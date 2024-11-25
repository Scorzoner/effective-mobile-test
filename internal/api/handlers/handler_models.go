package handlers

type BasicSongInfoJSON struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type additionalSongInfoJSON struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type UpdateRequestJSON struct {
	Id          int64  `json:"id"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type FilterRequest struct {
	GroupName             string
	SongName              string
	ReleaseDateLowerBound string
	ReleaseDateUpperBound string
	Text                  string
	Page                  int64
	PageSize              int64
}

type ListRowResult struct {
	Id          int32  `json:"id"`
	GroupName   string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}

type AddSongFailedExternalAPIResponse struct {
	Id          int32  `json:"id"`
	SongDetails string `json:"songDetails"`
}

type FilteredListResponse struct {
	FilteredRows any `json:"filteredRows"`
}
