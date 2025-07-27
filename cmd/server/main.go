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

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
