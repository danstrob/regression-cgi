package main

import (
	"html/template"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Data struct {
	X, Y             []float64 // Slices with X and Y values as data for the scatter plot.
	Intercept, Slope float64   // User guess for intercept and slope of the line graph.
	rssOLS, rssGuess float64   // Residual sum of squares of OLS model and user guess.
	ImagePath        string    // Path to plot image (PNG format).
}

func handler(w http.ResponseWriter, r *http.Request) {
	wg.Add(1)
	go filepath.Walk("images", RemoveOldFiles)
	d := &Data{
		X: []float64{5, 3, 6, 3, 5, 2, 0, 6, 8, 10},
		Y: []float64{3, 5, 3, 7, 0, 8, 6, 0, 0, 0}}

	// Convert HTML input from string to floats.
	keys := []string{"intercept", "slope"}
	floatMap := inputToFloat(r, keys)
	d.Intercept, d.Slope = floatMap[keys[0]], floatMap[keys[1]]

	// Run Regression, draw plot and serve HTML from template.
	d.rssOLS, d.rssGuess = Regression(d)
	d.ImagePath = DrawPlot(d)

	t, err := template.ParseFiles("templates/input.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, d)
	wg.Wait()
}

func main() {
	if err := cgi.Serve(http.HandlerFunc(handler)); err != nil {
		log.Fatal(err)
	}
}

// inputToFloat takes in a pointer to an http.Request and
// a slice of strings with keys to search for in the query.
// It returns a map of the keys (string) to converted floats.
func inputToFloat(r *http.Request, keys []string) map[string]float64 {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	keyMap := make(map[string]float64)
	for _, k := range keys {
		str := r.FormValue(k)
		if str != "" {
			keyMap[k], err = strconv.ParseFloat(str, 64)

			switch {
			case err != nil:
				log.Fatal(err)
			case keyMap[k] > 50:
				log.Fatalf("Invalid user input: Value for %s too large.", k)
			case keyMap[k] < -50:
				log.Fatalf("Invalid user input: Value for %s too small.", k)
			}
		}
	}
	return keyMap
}

// RemoveOldFiles is a WalkFunc which removes all files in root dir older than 30 seconds.
func RemoveOldFiles(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Fatal(err)
	}
	if !info.IsDir() {
		if time.Since(info.ModTime()).Seconds() > 30 {
			err := os.Remove(path)
			if err != nil {
				log.Print(err)
			}
		}
	}
	return err
}
