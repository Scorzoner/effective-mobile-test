package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Scorzoner/effective-mobile-test/internal/config"
)

type validator struct {
	Errors map[string]string
}

func newValidator() *validator {
	return &validator{Errors: make(map[string]string)}
}

func (v *validator) valid() bool {
	return len(v.Errors) == 0
}

func (v *validator) addError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *validator) check(ok bool, key, message string) {
	if !ok {
		v.addError(key, message)
	}
}

func convertAndValidateStringToInt64(v *validator, numberAsStr string, name string) int64 {
	v.check(len(numberAsStr) > 0, name, "should be provided")

	numberAsInt, err := strconv.ParseInt(numberAsStr, 10, 64)
	v.check(err == nil, name, "should be an integer")
	v.check(numberAsInt > 0, name, "should be positive")
	v.check(numberAsInt < 1<<32, name, fmt.Sprintf("should be less than %d", 1<<32))
	return numberAsInt
}

func convertAndValidateStringToDate(v *validator, dateAsStr string, name string) time.Time {
	date, err := time.Parse("02.01.2006", dateAsStr)
	v.check(err == nil, name,
		fmt.Sprintf("expected DD.MM.YYYY format, date provided: %v", dateAsStr))
	if err == nil {
		v.check(date.Before(time.Now()), name,
			fmt.Sprintf("expected to be in the past, date provided: %v", dateAsStr))
	}
	return date
}

/*func validateNumberAsString(v *validator, numberAsStr string, name string) {
	v.check(len(numberAsStr) > 0, name, "should be provided")

	numberAsInt, err := strconv.ParseInt(numberAsStr, 10, 64)
	v.check(err == nil, name, "should be an integer")
	v.check(numberAsInt > 0, name, "should be positive")
}*/

func validateBasicSongInfoJSON(v *validator, j *BasicSongInfoJSON, cfg *config.Config) {
	v.check(len(j.Group) > 0, "group", "should be provided")
	v.check(len(j.Group) <= cfg.MaxGroupNameLen, "group",
		fmt.Sprintf("should be no more than %v characters long, current length %v",
			cfg.MaxGroupNameLen, len(j.Group)))

	v.check(len(j.Song) > 0, "song", "should be provided")
	v.check(len(j.Song) <= cfg.MaxSongNameLen, "song",
		fmt.Sprintf("should be no more than %v characters long, current length %v",
			cfg.MaxSongNameLen, len(j.Song)))
}

func validateAdditionalSongInfoJSON(v *validator, j *additionalSongInfoJSON, cfg *config.Config) {
	releaseDate, err := time.Parse("02.01.2006", j.ReleaseDate)
	v.check(err == nil, "releaseDate",
		fmt.Sprintf("expected DD.MM.YYYY format, date provided: %v", j.ReleaseDate))
	if err == nil {
		v.check(releaseDate.Before(time.Now()), "releaseDate",
			fmt.Sprintf("expected to be in the past, date provided: %v", j.ReleaseDate))
	}

	v.check(len(j.Text) > 0, "text", "should be provided")
	v.check(len(j.Text) <= cfg.MaxSongLyricsLen, "text",
		fmt.Sprintf("should be no more than %v characters long, current length %v",
			cfg.MaxSongLyricsLen, len(j.Text)))

	v.check(len(j.Link) > 0, "link", "should be provided")
	v.check(len(j.Link) <= cfg.MaxSongLinkLen, "link",
		fmt.Sprintf("should be no more than %v characters long, current length %v",
			cfg.MaxSongLinkLen, len(j.Link)))
}

func validateUpdateRequestJSON(v *validator, j *UpdateRequestJSON, cfg *config.Config) {
	v.check(j.Id > 0, "id", "should be positive")

	asij := additionalSongInfoJSON{ReleaseDate: j.ReleaseDate, Text: j.Text, Link: j.Link}
	validateAdditionalSongInfoJSON(v, &asij, cfg)
}
