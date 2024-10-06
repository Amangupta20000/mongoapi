package router

import (
	"github.com/Amangupta20000/mongoapi/controller"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/movies", controller.GetAllMovies).Methods("GET")
	router.HandleFunc("/api/movie/{id}", controller.FindOneMovie).Methods("GET")
	router.HandleFunc("/api/movie", controller.CreateMovies).Methods("POST")
	router.HandleFunc("/api/movie/{id}", controller.MarkAsWatched).Methods("PUT") //update
	router.HandleFunc("/api/movie/{id}", controller.DeleteOneMovie).Methods("DELETE")
	router.HandleFunc("/api/delete-all-movie", controller.DeleteAllMovies).Methods("DELETE")
	router.HandleFunc("/api/weather/{city}", controller.CheckWeather).Methods("GET")

	return router
}
