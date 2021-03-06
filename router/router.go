package router

import (
	"fam-photos-server/middleware"

	"github.com/gorilla/mux"
)

// Router defines the endpoints for the API.
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/first_name", middleware.GetFirstName).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/families", middleware.GetFamilies).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/image_data", middleware.GetImageDataByFamily).Methods("POST", "OPTIONS")

	return router
}
