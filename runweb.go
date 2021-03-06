package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed htdocs proverbs templates
var embFS embed.FS

// runWEB registers handlers and starts the web server.
func runWEB() {
	var err error
	rtr := mux.NewRouter()
	rtr.Use(prometheusMiddleware)
	s := &http.Server{
		Handler:           rtr,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	htd, err := fs.Sub(embFS, "htdocs")
	if err != nil {
		log.Println("Error returning embedded directory htdocs:", err.Error())
	}

	prv, err := fs.Sub(embFS, "proverbs")
	if err != nil {
		log.Println("Error returning embedded directory proverbs:", err.Error())
	}

	rtr.HandleFunc("/mortgage", getWebiHtml)
	rtr.HandleFunc("/mortgagejson", getWebiJson)
	rtr.HandleFunc("/json1", getWebiJson)
	rtr.HandleFunc("/amortize", getAmorHtml)
	rtr.HandleFunc("/amortizejson", getAmorJson)
	rtr.HandleFunc("/json2", getAmorJson)
	rtr.HandleFunc("/test", test)

	rtr.Path("/prometheus").Handler(promhttp.Handler())
	rtr.Path("/webinfo").Handler(http.RedirectHandler("/mortgage", http.StatusMovedPermanently))
	rtr.Path("/webinfojson").Handler(http.RedirectHandler("/mortgagejson", http.StatusMovedPermanently))
	rtr.Path("/cake").Handler(http.RedirectHandler("/cake.html", http.StatusMovedPermanently))

	rtr.PathPrefix("/proverbs/").Handler(http.StripPrefix("/proverbs/", http.FileServer(http.FS(prv))))
	rtr.PathPrefix("/").Handler(http.FileServer(http.FS(htd)))

	err = s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

/*
    Elected to stick with basic startup directly above.
	Startup using a goroutine and signal catching is below.

		go func() {
			err := s.ListenAndServe()
			if err != nil {
				log.Fatal(err)
				return
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		sig := <-sigs
		secs := 3 * time.Second
		log.Println("Terminating in", secs, ": Received signal", sig)
		time.Sleep(secs)
		s.Shutdown(context.TODO())
*/
