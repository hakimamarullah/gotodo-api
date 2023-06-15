package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gotodo/controllers/todocontroller"
	"github.com/gotodo/middlewares"
	"github.com/gotodo/models"
	log "github.com/sirupsen/logrus"
)

var (
	PORT string
)

const (
	TODOS       = "/todos"
	TODOS_ID    = "/todos/{id}"
	BATCH_TODOS = "/batch/todos"
)

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	models.ConnectPostgreSQL()

	log.Info("GoTodo API is up and running!")
	router := mux.NewRouter()

	router.HandleFunc("/", todocontroller.Home)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middlewares.JWTMiddleware)

	api.HandleFunc(TODOS, todocontroller.AddTodoItem).Methods(http.MethodPost)
	api.HandleFunc(TODOS, todocontroller.GetAllTodos).Methods(http.MethodGet)
	api.HandleFunc(TODOS_ID, todocontroller.DeleteTodoById).Methods(http.MethodDelete)
	api.HandleFunc(TODOS_ID, todocontroller.UpdateTodoById).Methods(http.MethodPut)
	api.HandleFunc(BATCH_TODOS, todocontroller.BatchDeleteByIds).Methods(http.MethodDelete)

	fmt.Printf("Server up and running on port %s\n", PORT)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", PORT), router)

}
