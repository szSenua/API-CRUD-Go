package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//las estructuras son como objetos en java

type Movie struct {
	ID    string `json:"id"`
	Isbn  string `json:"isbn"`
	Title string `json:"title"`
	//Asociación de película con un director
	Director *Director `json:"director"`
	//El asterisco es un puntero, lo cual quiere decir que si creo una estructura llamada Director
	//se asociará a la estructura de la película

}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

//Función obtener películas

func getMovies(w http.ResponseWriter, r *http.Request) {
	//Establecemos el contenido como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	//Establecemos el contenido como JSON
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	//recorremos la lista de peliculas
	for index, item := range movies {
		//si encuentra una con el mismo id que le pasamos
		if item.ID == params["id"] {
			//la borra con la función append
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

//Función obtener película por id

func getMovie(w http.ResponseWriter, r *http.Request) {
	//Establecemos el contenido como un objeto JSON
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	//Recorre la lista de películas
	for _, item := range movies {
		//Si encuentra una con el id proporcionado
		if item.ID == params["id"] {
			//lo devuelve como un objeto json
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

//Función crear película

func createMovie(w http.ResponseWriter, r *http.Request) {
	//Establecemos el contenido como un objeto JSON
	w.Header().Set("Content-Type", "application/json")
	//Creamos una variable de tipo pelicula
	var movie Movie
	//decodificamos los datos entrantes y los transformamos en un objeto json
	_ = json.NewDecoder(r.Body).Decode(&movie)
	//generamos un identificador único para cada película y los parseamos a string
	//que será el ID
	movie.ID = strconv.Itoa(rand.Intn(1000000000))
	//la película recién creada se agrega al final de la lista de películas existentes
	//(almacenada en una variable global llamada "movies") utilizando la función "append".
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

//Función actualizar película

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//Establecemos el contenido como un objeto JSON
	w.Header().Set("Content-Type", "application/json")
	//lo que vamos a hacer básicamente es borrar la pelicula que coincida con el id y añadir
	//la nueva modificada
	params := mux.Vars(r)

	//Recorremos la lista comparando el id proporcionado
	for index, item := range movies {
		if item.ID == params["id"] {
			//borramos la pelicula que coincida usando la función append
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			//se decodifica el objeto JSON enviado en la solicitud HTTP entrante
			//y se crea una nueva variable "movie" de tipo "Movie" utilizando la función "json.NewDecoder"
			_ = json.NewDecoder(r.Body).Decode(&movie)
			//El identificador de la nueva película se establece en el identificador de la película original
			//(obtenido de los parámetros de la URL)
			movie.ID = params["id"]
			//Agregamos la nueva película al último lugar de la lista
			movies = append(movies, movie)
			//y se devuelve como respuesta JSON utilizando la función "json.NewEncoder" y la respuesta HTTP.
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {

	//Definimos las rutas
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "454657", Title: "Movie Two", Director: &Director{Firstname: "Steven", Lastname: "Smith"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))

}
