package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

//go:embed app
var embededFiles embed.FS

func main() {
	reverseProxyURL, ok := os.LookupEnv("REVERSE_PROXY_URL")
	if !ok || reverseProxyURL == "" {
		golog.Fatalf("Web: environment variable not declared: reverseProxyURL")
	}
	webPort, ok := os.LookupEnv("WEB_PORT")
	if !ok || webPort == "" {
		golog.Fatalf("Web: environment variable not declared: webPort")
	}
	e := echo.New()
	useOs := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useOs))
	e.GET("/", echo.WrapHandler(assetHandler))
}

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")

		return http.FS(os.DirFS("app"))
	}

	log.Print("using embed mode")

	fsys, err := fs.Sub(embededFiles, "app")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
