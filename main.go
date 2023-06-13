package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	PORT string
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprint(w, "Hello GoTodo!")
}

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", Home)
	fmt.Printf("Server up and running on port %s", PORT)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", PORT), router)

}
