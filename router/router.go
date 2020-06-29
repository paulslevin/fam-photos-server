package router

import (
	"fam-photos-server/middleware"

	"github.com/gorilla/mux"
)

// Router is exported and used in application.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/first_name", middleware.GetFirstName).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/families", middleware.GetFamilies).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/image_urls", middleware.GetImageURLsByFamily).Methods("POST", "OPTIONS")

	return router
}
