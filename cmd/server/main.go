package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"money-map/internal/api"
	"money-map/internal/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// the path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("could not load .env file")
	}
	database.Connect()
	database.CreateTables()

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/income", api.GetIncome).Methods("GET")
	apiRouter.HandleFunc("/expenses", api.GetExpenses).Methods("GET")
	apiRouter.HandleFunc("/income", api.AddIncome).Methods("POST")
	apiRouter.HandleFunc("/expenses", api.AddExpense).Methods("POST")
	apiRouter.HandleFunc("/income/{id}", api.DeleteIncome).Methods("DELETE")
	apiRouter.HandleFunc("/expenses/{id}", api.DeleteExpense).Methods("DELETE")
	apiRouter.HandleFunc("/funds", api.GetFunds).Methods("GET")
	apiRouter.HandleFunc("/funds", api.UpdateFunds).Methods("PUT")
	apiRouter.HandleFunc("/delta", api.GetDelta).Methods("GET")
	apiRouter.HandleFunc("/upload", api.UploadCSV).Methods("POST")

	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}