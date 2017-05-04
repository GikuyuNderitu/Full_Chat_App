package main

import (
	"net/http"
)

func main() {

	// Handle Serving Static Files
	assetFiles := http.FileServer(http.Dir("assets"))
	distFiles := http.FileServer(http.Dir("dist"))
	bowerFiles := http.FileServer(http.Dir("bower_components"))
	srcFiles := http.FileServer(http.Dir("src"))

	http.HandleFunc("/", HomeHandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", assetFiles))
	http.Handle("/dist/", http.StripPrefix("/dist/", distFiles))
	http.Handle("/bower_components/", http.StripPrefix("/bower_components/", bowerFiles))
	http.Handle("/src/", http.StripPrefix("/src/", srcFiles))
	http.ListenAndServe(":8080", nil)
}

// HomeHandler serves the index.html file
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
