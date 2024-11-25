package main

import (
	_ "github.com/Scorzoner/effective-mobile-test/docs"
	"github.com/Scorzoner/effective-mobile-test/internal/server"
)

//	@title		Music Library API
//	@version	1.0

// @host		localhost:8080
// @BasePath	/
func main() {
	server.Run()
}
