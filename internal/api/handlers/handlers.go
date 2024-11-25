package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Scorzoner/effective-mobile-test/internal/api/badresponses"
	"github.com/Scorzoner/effective-mobile-test/internal/api/jsonutil"
	"github.com/Scorzoner/effective-mobile-test/internal/config"
	"github.com/Scorzoner/effective-mobile-test/internal/database"
	"github.com/Scorzoner/effective-mobile-test/internal/logger"
	"github.com/Scorzoner/effective-mobile-test/internal/models"
)

type HandleQueries struct {
	connections *sql.DB
	q           *database.Queries
	cfg         config.Config
}

func NewHandlerQueries(connections *sql.DB, cfg config.Config) (*HandleQueries, error) {
	queries, err := database.NewQueries(connections)
	if err != nil {
		return nil, err
	}
	return &HandleQueries{connections: connections, q: queries, cfg: cfg}, nil
}

func (hq *HandleQueries) RequestLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)
		logger.Zap.Debug(
			"Method:", r.Method,
			"Duration:", time.Since(start),
			"URI:", r.RequestURI,
		)
	})
}

func (hq *HandleQueries) ResponseLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := logResponseWriter{
			ResponseWriter: w,
			status:         0,
			size:           0,
		}

		h.ServeHTTP(&lw, r)

		logger.Zap.Debug(
			"Method:", r.Method,
			"URI:", r.RequestURI,
			"Status:", lw.status,
			"Size:", lw.size,
		)
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *logResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

func (r *logResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.status = statusCode
}

func fetchSongDetails(bsi BasicSongInfoJSON, externalAPIURL string) (*additionalSongInfoJSON, error) {
	parsedURL, err := url.Parse(externalAPIURL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, fmt.Errorf("invalid or unsupported URL scheme: %s", externalAPIURL)
	}

	if parsedURL.Host == "" {
		return nil, fmt.Errorf("URL must have a valid host: %s", externalAPIURL)
	}

	fullURL := fmt.Sprintf("%s/info?group=%s&song=%s",
		parsedURL.String(), url.PathEscape(bsi.Group), url.PathEscape(bsi.Song))

	resp, errResp := http.Get(fullURL)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %s", resp.Status)
	}

	body, errBody := io.ReadAll(resp.Body)
	if errBody != nil {
		return nil, errBody
	}

	var asi additionalSongInfoJSON
	errJSON := json.Unmarshal(body, &asi)
	if errJSON != nil {
		return nil, errJSON
	}

	return &asi, nil
}

// @Summary		Adds song into library
// @Tags			music-library
// @Description	Makes a request into externalAPIURL, if it fails, returns status 201,
// @Description saves basic song info and writes encountered errors into response body (field "songDetails")
// @Accept			json
// @Produce		json
// @Param			BasicSongInfoJSON	body		BasicSongInfoJSON	true	"group and song names"
// @Success		200					{object}	models.IdResponse
// @Success		201					{object}	AddSongFailedExternalAPIResponse
// @Failure		400					{object}	models.ErrorResponse
// @Failure		422					{object}	models.ErrorResponse
// @Failure		500					{object}	models.ErrorResponse
// @Router			/music-library/song [post]
func (hq *HandleQueries) AddSong(w http.ResponseWriter, r *http.Request) {
	var requestJSON BasicSongInfoJSON
	err := jsonutil.ReadJSON(w, r, &requestJSON)
	if err != nil {
		badresponses.BadRequestResponse(w, r, err)
		return
	}
	logger.Zap.Debug(fmt.Sprintf("addSong request json: %v", requestJSON))

	v := newValidator()
	validateBasicSongInfoJSON(v, &requestJSON, &hq.cfg)
	if !v.valid() {
		badresponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	bsi := models.BasicSongInfo{Id: 0, GroupName: requestJSON.Group, SongName: requestJSON.Song}

	err = hq.q.AddSong(&bsi)
	if err != nil {
		badresponses.BadRequestResponse(w, r, fmt.Sprintf("failed to add song: %s", err.Error()))
		return
	}

	sendOnlyBasicDataWasWrittenResponse := func(externalErr error) {
		result := map[string]any{
			"id":          bsi.Id,
			"songDetails": fmt.Sprintf("were not acquired due to: %s", externalErr.Error())}
		err = jsonutil.WriteJSON(w, http.StatusCreated, result, nil)
		if err != nil {
			badresponses.InternalServerErrorResponse(w, r, fmt.Errorf("failed writing response: %w", err))
			return
		}
	}

	externalResponseJSON, err := fetchSongDetails(requestJSON, hq.cfg.ExternalAPIURL)
	if err != nil {
		errorResult := fmt.Errorf("failed to fetch song details from external api: %w", err)
		logger.Zap.Debug(errorResult)
		sendOnlyBasicDataWasWrittenResponse(errorResult)
		return
	}

	validateAdditionalSongInfoJSON(v, externalResponseJSON, &hq.cfg)
	if !v.valid() {
		errorResult := fmt.Errorf("failed to validate song details from external api: %v", v.Errors)
		logger.Zap.Debug(errorResult)
		sendOnlyBasicDataWasWrittenResponse(errorResult)
		return
	}

	rd, _ := time.Parse("02.01.2006", externalResponseJSON.ReleaseDate)
	asi := models.AdditionalSongInfo{
		ReleaseDate: rd,
		SongLyrics:  externalResponseJSON.Text,
		Link:        externalResponseJSON.Link}

	err = hq.q.UpdateSongInfo(bsi.Id, &asi)
	if err != nil {
		errorResult := fmt.Errorf("failed to add additional info from external source: %w", err)
		logger.Zap.Debug(errorResult)
		sendOnlyBasicDataWasWrittenResponse(errorResult)
		return
	}

	result := map[string]any{"id": bsi.Id}
	err = jsonutil.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		errorResult := fmt.Errorf("failed writing response: %w", err)
		logger.Zap.Error(fmt.Errorf("failed writing response: %w", err))
		sendOnlyBasicDataWasWrittenResponse(errorResult)
		return
	}
}

// @Summary		Deletes song from library
// @Tags			music-library
// @Description	Returns provided id if deletion succeeded
// @Accept			plain
// @Produce		json
// @Param			id	query		int	true	"song id"
// @Success		200	{object}	models.IdResponse
// @Failure		400	{object}	models.ErrorResponse
// @Failure		422	{object}	models.ErrorResponse
// @Failure		500	{object}	models.ErrorResponse
// @Router			/music-library/song [delete]
func (hq *HandleQueries) DeleteSong(w http.ResponseWriter, r *http.Request) {
	stringId := r.URL.Query().Get("id")

	v := newValidator()
	songId := convertAndValidateStringToInt64(v, stringId, "id")
	if !v.valid() {
		badresponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err := hq.q.DeleteSong(songId)
	if err == database.ErrSongNotFound {
		badresponses.BadRequestResponse(w, r, fmt.Sprintf("failed to delete song: %s", err.Error()))
		return
	}
	if err != nil {
		badresponses.InternalServerErrorResponse(w, r, fmt.Errorf("failed to delete song: %w", err))
		return
	}

	result := map[string]any{"id": songId}
	err = jsonutil.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		badresponses.InternalServerErrorResponse(w, r, fmt.Errorf("failed writing response: %w", err))
		return
	}
}

// @Summary		Fetches lyrics divided into verses
// @Tags			music-library
// @Description	If page/pageSize combination results in an empty page, you will get status code 400 :)
// @Accept			plain
// @Produce		json
// @Param			id			query		int	true	"song id"
// @Param			page		query		int	true	"page number"
// @Param			pageSize	query		int	true	"number of verses per page"
// @Success		200			{object}	models.VersesResponse
// @Failure		400			{object}	models.ErrorResponse
// @Failure		422			{object}	models.ErrorResponse
// @Failure		500			{object}	models.ErrorResponse
// @Router			/music-library/lyrics [get]
func (hq *HandleQueries) GetSongLyrics(w http.ResponseWriter, r *http.Request) {
	stringId := r.URL.Query().Get("id")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")

	v := newValidator()
	songId := convertAndValidateStringToInt64(v, stringId, "id")
	pageAsInt := convertAndValidateStringToInt64(v, page, "page")
	pageSizeAsInt := convertAndValidateStringToInt64(v, pageSize, "pageSize")
	if !v.valid() {
		badresponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	var lyrics string
	err := hq.q.GetLyrics(songId, &lyrics)
	logger.Zap.Debug(fmt.Sprintf("lyrics fetched: %s", lyrics))
	if err == database.ErrSongHasNoLyrics {
		badresponses.BadRequestResponse(w, r, fmt.Sprintf("failed to fetch song lyrics: %s", err.Error()))
		return
	}
	if err == database.ErrSongNotFound {
		badresponses.BadRequestResponse(w, r, fmt.Sprintf("failed to fetch song lyrics: %s", err.Error()))
		return
	}
	if err != nil {
		badresponses.InternalServerErrorResponse(w, r, fmt.Errorf("failed to get song lyrics: %w", err))
		return
	}

	verses := strings.Split(lyrics, "\n\n")

	//overflow detection
	x := pageAsInt * pageSizeAsInt
	if pageAsInt != 0 && x/pageAsInt != pageSizeAsInt || x > 1<<31 {
		badresponses.BadRequestResponse(w, r,
			fmt.Sprintf("provided page (%s) and pageSize (%s) are too big, "+
				"consider providing numbers that dont overflow int", page, pageSize))
		return
	}

	lowerBound := int((pageAsInt - 1) * pageSizeAsInt)
	if len(verses) <= lowerBound {
		badresponses.BadRequestResponse(w, r,
			fmt.Sprintf("provided page (%s) and pageSize(%s) combination "+
				"yeilds an empty page, song only has %d verses", page, pageSize, len(verses)))
		return
	}

	var upperBound int
	if int(pageAsInt*pageSizeAsInt) <= len(verses) {
		upperBound = int(pageAsInt * pageSizeAsInt)
	} else {
		upperBound = len(verses)
	}

	result := map[string]any{
		"verses": verses[lowerBound:upperBound]}
	err = jsonutil.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		badresponses.InternalServerErrorResponse(w, r, fmt.Errorf("failed writing response: %w", err))
		return
	}
}

// @Summary		Fetches song data in pages
// @Tags			music-library
// @Description	page and pageSize are required, every other field is a filter, if it's empty, it is treated as absence of filter
// @Accept			plain
// @Produce		json
// @Param			group				query		string	false	"group name"
// @Param			song				query		string	false	"song name"
// @Param			releaseDateLower	query		string	false	"dates before this will not show up"
// @Param			releaseDateUpper	query		string	false	"dates after this will not show up"
// @Param			text				query		string	false	"lyrics"
// @Param			page				query		int	true	"page number"
// @Param			pageSize			query		int	true	"number of songs displayed per page"
// @Success		200					{object}	FilteredListResponse
// @Failure		400					{object}	models.ErrorResponse
// @Failure		422					{object}	models.ErrorResponse
// @Failure		500					{object}	models.ErrorResponse
// @Router			/music-library/list [get]
func (hq *HandleQueries) GetFilteredList(w http.ResponseWriter, r *http.Request) {
	var filter FilterRequest
	rq := r.URL.Query()
	filter.GroupName = rq.Get("group")
	filter.SongName = rq.Get("song")
	filter.ReleaseDateLowerBound = rq.Get("releaseDateLower")
	filter.ReleaseDateUpperBound = rq.Get("releaseDateUpper")
	filter.Text = rq.Get("text")

	v := newValidator()
	filter.Page = convertAndValidateStringToInt64(v, rq.Get("page"), "page")
	filter.PageSize = convertAndValidateStringToInt64(v, rq.Get("pageSize"), "pageSize")

	var dbFilter database.ListFilter

	if filter.GroupName == "" {
		dbFilter.GroupName.Valid = false
	} else {
		dbFilter.GroupName.Valid = true
		dbFilter.GroupName.String = filter.GroupName
	}

	if filter.SongName == "" {
		dbFilter.SongName.Valid = false
	} else {
		dbFilter.SongName.Valid = true
		dbFilter.SongName.String = filter.SongName
	}

	if filter.ReleaseDateLowerBound == "" {
		dbFilter.ReleaseDateLowerBound.Valid = false
	} else {
		dbFilter.ReleaseDateLowerBound.Time = convertAndValidateStringToDate(
			v, filter.ReleaseDateLowerBound, "releaseDateLower")
		dbFilter.ReleaseDateLowerBound.Valid = true
	}

	if filter.ReleaseDateUpperBound == "" {
		dbFilter.ReleaseDateUpperBound.Valid = false
	} else {
		dbFilter.ReleaseDateUpperBound.Time = convertAndValidateStringToDate(
			v, filter.ReleaseDateUpperBound, "releaseDateUpper")
		dbFilter.ReleaseDateUpperBound.Valid = true
	}

	if filter.Text == "" {
		dbFilter.Lyrics.Valid = false
	} else {
		dbFilter.Lyrics.Valid = true
		dbFilter.Lyrics.String = filter.Text
	}

	dbFilter.Limit = int32(filter.PageSize)
	dbFilter.Offset = int32((filter.Page - 1) * filter.PageSize)

	if !v.valid() {
		badresponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	resultNullable, err := hq.q.GetFilteredList(&dbFilter)
	if err != nil {
		badresponses.InternalServerErrorResponse(
			w, r, fmt.Errorf("failed to get filtered list: %w", err))
		return
	}

	var result []ListRowResult
	for _, row := range resultNullable {
		var res ListRowResult
		res.Id = int32(row.Id)
		res.GroupName = row.GroupName
		res.SongName = row.SongName
		if row.ReleaseDate.Valid {
			res.ReleaseDate = time.Time.Format(row.ReleaseDate.Time, "02.01.2006")
		}
		if row.SongLyrics.Valid {
			res.Text = row.SongLyrics.String
		}
		if row.Link.Valid {
			res.Link = row.Link.String
		}

		result = append(result, res)
	}

	resultMap := map[string]any{"filteredRows": result}
	err = jsonutil.WriteJSON(w, http.StatusOK, resultMap, nil)
	if err != nil {
		badresponses.InternalServerErrorResponse(w, r, fmt.Errorf("failed writing response: %w", err))
		return
	}
}

// @Summary		Updates song info
// @Tags			music-library
// @Description	You need to provide id and 3 other fields, on success returns provided id
// @Accept			json
// @Produce		json
// @Param			UpdateRequestJSON	body		UpdateRequestJSON	true	"id and additional info"
// @Success		200					{object}	models.IdResponse
// @Failure		400					{object}	models.ErrorResponse
// @Failure		422					{object}	models.ErrorResponse
// @Failure		500					{object}	models.ErrorResponse
// @Router			/music-library/song [put]
func (hq *HandleQueries) UpdateSongInfo(w http.ResponseWriter, r *http.Request) {
	var requestJSON UpdateRequestJSON
	err := jsonutil.ReadJSON(w, r, &requestJSON)
	if err != nil {
		badresponses.BadRequestResponse(w, r, fmt.Sprintf("failed to update song info: %s", err.Error()))
		return
	}

	v := newValidator()
	validateUpdateRequestJSON(v, &requestJSON, &hq.cfg)
	if !v.valid() {
		badresponses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	rd, _ := time.Parse("02.01.2006", requestJSON.ReleaseDate)
	asi := models.AdditionalSongInfo{
		ReleaseDate: rd,
		SongLyrics:  requestJSON.Text,
		Link:        requestJSON.Link}

	err = hq.q.UpdateSongInfo(requestJSON.Id, &asi)
	if err == database.ErrSongNotFound {
		badresponses.BadRequestResponse(w, r, fmt.Sprintf("failed to update song info: %s", err.Error()))
		return
	}
	if err != nil {
		badresponses.InternalServerErrorResponse(w, r, fmt.Sprintf("failed to update song info: %s", err.Error()))
		return
	}

	result := map[string]any{"id": requestJSON.Id}
	err = jsonutil.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		badresponses.InternalServerErrorResponse(w, r, fmt.Errorf("failed writing response: %w", err))
		return
	}
}
