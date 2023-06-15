package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"encoding/json"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	PORT string
)

const (
	APPLICATION_JSON = "application/json"
	TODOS            = "/todos"
	TODOS_ID         = "/todos/{id}"
	BATCH_TODOS      = "/batch/todos"
)

/*DB CONNECTION*/
var dsn string = "host=localhost user=go password=go dbname=gotodo port=5432 sslmode=disable TimeZone=Asia/Jakarta"
var db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})

/* MODEL */
type TodoItem struct {
	gorm.Model
	Description string
	Completed   bool `gorm:"default=false"`
}

type ResponseBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
	Code    int         `json:"code"`
}

/* Handler */
func Home(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func AddTodoItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	desc := struct {
		Description string `json:"description"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&desc)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseBody{Message: "Invalid data", Code: http.StatusBadRequest})
		return
	}

	todo := &TodoItem{Description: desc.Description}
	db.Create(&todo)
	db.Last(&todo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ResponseBody{Message: "Todo added!", Data: todo, Code: http.StatusCreated})
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	var todos []TodoItem
	result := db.Find(&todos)
	json.NewEncoder(w).Encode(ResponseBody{Message: "Data retrieved", Count: int(result.RowsAffected), Data: todos, Code: http.StatusOK})

}

func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	params := mux.Vars(r)
	result := db.Delete(&TodoItem{}, params["id"])
	json.NewEncoder(w).Encode(ResponseBody{Message: "Todo deleted", Data: params["id"], Count: int(result.RowsAffected)})
}

func UpdateTodoItemById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)
	params := mux.Vars(r)
	var todo TodoItem
	result := db.First(&todo, params["id"])
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseBody{Message: "Todo not found", Code: http.StatusNotFound})
		return
	}
	var updateTodo TodoItem
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&updateTodo)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseBody{Message: "Invalid data", Code: http.StatusBadRequest})
		return
	}

	result = db.Model(&todo).Updates(updateTodo)

	json.NewEncoder(w).Encode(ResponseBody{Message: fmt.Sprintf("Item %d updated", todo.ID), Count: int(result.RowsAffected), Code: http.StatusOK})

}

func BatchDeleteByIds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", APPLICATION_JSON)

	data := struct {
		Ids []int `json:"ids"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseBody{Message: "Invalid data", Code: http.StatusBadRequest})
		return
	}

	result := db.Delete(&TodoItem{}, data.Ids)
	json.NewEncoder(w).Encode(ResponseBody{Message: "Items Deleted", Count: int(result.RowsAffected), Code: http.StatusOK})
}

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	db.Migrator().CreateTable(&TodoItem{})

	log.Info("GoTodo API is up and running!")
	router := mux.NewRouter()

	router.HandleFunc("/", Home)
	router.HandleFunc(TODOS, AddTodoItem).Methods(http.MethodPost)
	router.HandleFunc(TODOS, GetTodos).Methods(http.MethodGet)
	router.HandleFunc(TODOS_ID, DeleteTodoById).Methods(http.MethodDelete)
	router.HandleFunc(TODOS_ID, UpdateTodoItemById).Methods(http.MethodPut)
	router.HandleFunc(BATCH_TODOS, BatchDeleteByIds).Methods(http.MethodDelete)

	fmt.Printf("Server up and running on port %s", PORT)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", PORT), router)

}
