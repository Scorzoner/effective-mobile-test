package models

import (
	"database/sql"
	"time"
)

type BasicSongInfo struct {
	Id        int64
	GroupName string
	SongName  string
}

type AdditionalSongInfo struct {
	ReleaseDate time.Time
	SongLyrics  string
	Link        string
}

type FullSongInfo struct {
	Id          int64          `json:"id"`
	GroupName   string         `json:"group"`
	SongName    string         `json:"song"`
	ReleaseDate sql.NullTime   `json:"releaseDate"`
	SongLyrics  sql.NullString `json:"text"`
	Link        sql.NullString `json:"link"`
}

type ErrorResponse struct {
	Error any `json:"errors"`
}

type IdResponse struct {
	Id int32 `json:"id"`
}

type VersesResponse struct {
	Verses any `json:"verses"`
}
