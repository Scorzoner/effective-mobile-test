package router

import (
	"github.com/Scorzoner/effective-mobile-test/internal/api/badresponses"
	"github.com/Scorzoner/effective-mobile-test/internal/api/handlers"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(hq *handlers.HandleQueries) *chi.Mux {
	router := chi.NewRouter()

	router.MethodNotAllowed(badresponses.MethodNotAllowedResponse)
	router.NotFound(badresponses.NotFoundResponse)

	router.Group(func(r chi.Router) {
		r.Use(hq.RequestLogging)

		r.Post("/music-library/song", hq.AddSong)
		r.Put("/music-library/song", hq.UpdateSongInfo)
		r.Delete("/music-library/song", hq.DeleteSong)
	})

	router.Group(func(r chi.Router) {
		r.Use(hq.ResponseLogging)

		r.Get("/music-library/lyrics", hq.GetSongLyrics)
		r.Get("/music-library/list", hq.GetFilteredList)
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	return router
}
