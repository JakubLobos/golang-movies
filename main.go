package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// struct, like an object in js
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

// json:"firstname this is how we define json, encoding the json that will come from postman
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// slice of the type Movie
var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//import "w", which is the response we want to send and encode as json and encode movies, which is what we want to send in this slice
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// we pass the request to mux.Vars, because params are part of our request (query params)
	params := mux.Vars(r)

	// something like foreach in js
	for index, item := range movies {
		if item.ID == params["id"] {
			// we take movies[:index] and append everything movies[index+1:] except movies[:index]
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	// JSON decoding: We create a new JSON decoder that reads data from the r.Body (HTTP request body), and then decodes this data into a movie structure. We ignore the returned error (via _), which means that any errors during decoding are ignored.
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//	set json content type
	//	params
	//	loop over the movies range
	//	delete the movie with the i.d that you've sent
	//	add a new movie - the movie that we send in the body of postman

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)

			movies := append(movies[:index], movies[index+1:]...)

			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie one", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "21370", Title: "Movie two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	// If someone hits /movies then execute the getMovies function
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server on port :8000 \n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
