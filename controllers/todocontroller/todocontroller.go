package todocontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gotodo/helpers"
	"github.com/gotodo/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	APPLICATION_JSON = "application/json"
	TODOS            = "/todos"
	TODOS_ID         = "/todos/{id}"
	BATCH_TODOS      = "/batch/todos"
)

func Home(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	helpers.ResponseJSON(w, models.ResponseBody{Data: map[string]bool{"alive": true}})
}

func AddTodoItem(w http.ResponseWriter, r *http.Request) {
	desc := struct {
		Description string `json:"description"`
	}{}

	id := helpers.GetUserIdFromContext(r.Context())

	err := json.NewDecoder(r.Body).Decode(&desc)
	if err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "invalid data", Code: http.StatusBadRequest})
		return
	}

	todo := &models.TodoItem{Description: desc.Description, UserID: id}
	models.DB.Create(&todo)
	models.DB.Last(&todo)
	helpers.ResponseJSON(w, models.ResponseBody{Message: "Todo added!", Data: todo, Code: http.StatusCreated})
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.TodoItem
	userId := helpers.GetUserIdFromContext(r.Context())
	result := models.DB.Where(&models.TodoItem{UserID: userId}).Find(&todos)
	helpers.ResponseJSON(w, models.ResponseBody{Message: "Data retrieved", Count: int(result.RowsAffected), Data: todos})

}

func DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result := models.DB.Where("user_id = ? AND id = ?", helpers.GetUserIdFromContext(r.Context()), params["id"]).Delete(&models.TodoItem{})
	if result.RowsAffected == 0 {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Todo Item Not Found for the given id", Code: http.StatusNotFound})
		return
	}
	helpers.ResponseJSON(w, models.ResponseBody{Message: "Todo deleted", Data: params["id"], Count: int(result.RowsAffected)})
}

func UpdateTodoById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var todo models.TodoItem
	userId := helpers.GetUserIdFromContext(r.Context())
	result := models.DB.Where("user_id = ? AND id = ?", userId, params["id"]).First(&todo)
	if result.Error != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Todo not found", Code: http.StatusNotFound})
		return
	}
	var updateTodo models.TodoItem
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&updateTodo)

	if err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Invalid data", Code: http.StatusBadRequest})
		return
	}

	result = models.DB.Model(&todo).Updates(updateTodo)
	helpers.ResponseJSON(w, models.ResponseBody{Message: fmt.Sprintf("Item %d updated", todo.ID), Count: int(result.RowsAffected)})

}

func BatchDeleteByIds(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Ids []int `json:"ids"`
	}{}
	userId := helpers.GetUserIdFromContext(r.Context())
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Invalid data", Code: http.StatusBadRequest})
		return
	}

	result := models.DB.Where("user_id = ? AND id IN (?)", userId, data.Ids).Delete(&models.TodoItem{})
	if result.RowsAffected < int64(len(data.Ids)) {
		helpers.ResponseJSON(w, models.ResponseBody{Message: "Some data could not be found on the database. So, not deleted", Count: int(result.RowsAffected), Code: http.StatusOK})
		return
	}
	helpers.ResponseJSON(w, models.ResponseBody{Message: "Items Deleted", Count: int(result.RowsAffected), Code: http.StatusOK})
}
