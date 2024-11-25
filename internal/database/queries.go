package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Scorzoner/effective-mobile-test/internal/models"
)

type Queries struct {
	db       *sql.DB
	prepared map[string]*sql.Stmt // A map of prepared statements for use in database package functions
}

// Maps function names of database package to respective queries,
// used for preparing statements when Queries gets initialized
var queryMap = map[string]string{
	"AddSong": `
		INSERT INTO music_library (group_name, song_name)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
		RETURNING song_id`,
	"UpdateSongInfo": `
		UPDATE music_library
		SET release_date=$2, song_lyrics=$3, link=$4
		WHERE song_id=$1`,
	"isSongIdPresent": `
		SELECT EXISTS(
			SELECT 1 FROM music_library
			WHERE song_id=$1)`,
	"getSongId": `
		SELECT song_id FROM music_library
		WHERE group_name=$1 AND song_name=$2`,
	"DeleteSong": `
		DELETE FROM music_library
		WHERE song_id=$1`,
	"GetLyrics": `
		SELECT song_lyrics FROM music_library
		WHERE song_id=$1`,
	"GetFilteredList": `
		SELECT song_id,
			group_name,
			song_name,
			release_date,
			song_lyrics,
			link
		FROM music_library
		WHERE (group_name ILIKE '%' || $1 || '%' OR $1 IS NULL)
		AND (song_name ILIKE '%' || $2 || '%' OR $2 IS NULL)
		AND (release_date>=$3 OR $3 IS NULL)
		AND (release_date<=$4 OR $4 IS NULL)
		AND (song_lyrics ILIKE '%' || $5 || '%' OR $5 IS NULL)
		ORDER BY song_id ASC
		LIMIT $6 OFFSET $7`,
}

// Prepares statements from pre-written queries
func NewQueries(db *sql.DB) (*Queries, error) {
	var err error
	prepared := make(map[string]*sql.Stmt)
	for functionName, query := range queryMap {
		prepared[functionName], err = db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("failed to prepare sql statements at function %s: %s", functionName, err)
		}
	}

	return &Queries{db: db, prepared: prepared}, nil
}

// Returns [ErrSongAlreadyExists] if song already in the database
func (q *Queries) AddSong(song *models.BasicSongInfo) error {
	err := q.getSongId(song)
	if err != ErrSongNotFound {
		return ErrSongAlreadyExists
	}
	args := []any{song.GroupName, song.SongName}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = q.prepared["AddSong"].QueryRowContext(ctx, args...).Scan(&song.Id)
	return err
}

// Returns [ErrSongNotFound] if there's no matching song in the database
func (q *Queries) UpdateSongInfo(songId int64, info *models.AdditionalSongInfo) error {
	exists, err := q.isSongIdPresent(songId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSongNotFound
	}

	args := []any{songId, info.ReleaseDate, info.SongLyrics, info.Link}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = q.prepared["UpdateSongInfo"].ExecContext(ctx, args...)
	return err
}

// Returns whether a song with given id exists in the database
func (q *Queries) isSongIdPresent(songId int64) (bool, error) {
	if songId == 0 {
		return false, nil
	}

	args := []any{songId}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	exists := false
	err := q.prepared["isSongIdPresent"].QueryRowContext(ctx, args...).Scan(&exists)
	return exists, err
}

// Returns [ErrSongNotFound] if song is not present in the database,
// first checks id, then names,
// writes song_id into song.Id
func (q *Queries) getSongId(song *models.BasicSongInfo) error {
	exists, err := q.isSongIdPresent(song.Id)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	args := []any{song.GroupName, song.SongName}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = q.prepared["getSongId"].QueryRowContext(ctx, args...).Scan(&song.Id)
	if err == sql.ErrNoRows {
		return ErrSongNotFound
	}

	return nil
}

// Returns [ErrSongNotFound] if there's no song in the database
func (q *Queries) DeleteSong(songId int64) error {
	exists, err := q.isSongIdPresent(songId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSongNotFound
	}

	args := []any{songId}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err = q.prepared["DeleteSong"].ExecContext(ctx, args...)
	return err
}

// Writes lyrics into [info.SongLyrics].
// Returns [ErrSongNotFound] if there's no song in the database.
// Returns [ErrSongHasNoLyrics] if no lyrics were provided.
func (q *Queries) GetLyrics(songId int64, lyrics *string) error {
	exists, err := q.isSongIdPresent(songId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrSongNotFound
	}

	args := []any{songId}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var lyricsOrNull sql.NullString
	err = q.prepared["GetLyrics"].QueryRowContext(ctx, args...).Scan(&lyricsOrNull)
	if err != nil {
		return err
	}

	if !lyricsOrNull.Valid {
		return ErrSongHasNoLyrics
	}
	*lyrics = lyricsOrNull.String
	return nil
}

type ListFilter struct {
	GroupName             sql.NullString
	SongName              sql.NullString
	ReleaseDateLowerBound sql.NullTime
	ReleaseDateUpperBound sql.NullTime
	Lyrics                sql.NullString
	Limit                 int32
	Offset                int32
}

func (q *Queries) GetFilteredList(filter *ListFilter) ([]models.FullSongInfo, error) {
	args := []any{
		filter.GroupName,
		filter.SongName,
		filter.ReleaseDateLowerBound,
		filter.ReleaseDateUpperBound,
		filter.Lyrics,
		filter.Limit,
		filter.Offset,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	rows, err := q.prepared["GetFilteredList"].QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FullSongInfo

	for rows.Next() {
		var row models.FullSongInfo

		err := rows.Scan(
			&row.Id,
			&row.GroupName,
			&row.SongName,
			&row.ReleaseDate,
			&row.SongLyrics,
			&row.Link,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}
