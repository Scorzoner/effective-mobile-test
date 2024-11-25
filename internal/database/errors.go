package database

import (
	"errors"
)

var (
	ErrSongNotFound      = errors.New("no matching record in database")
	ErrSongAlreadyExists = errors.New("given song already exists in database")
	ErrSongHasNoLyrics   = errors.New("given song does not have any lyrics assigned")
)
