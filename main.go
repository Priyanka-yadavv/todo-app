package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

var tasks []Task
var currentID int

type App struct {
	Router *mux.Router
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/tasks", app.getTasks).Methods("GET")
	app.Router.HandleFunc("/task/{id}", app.readTask).Methods("GET")
	app.Router.HandleFunc("/task", app.createTask).Methods("POST")
	app.Router.HandleFunc("/task/{id}", app.updateTask).Methods("PUT")
	app.Router.HandleFunc("/task/{id}", app.deleteTask).Methods("DELETE")
}

func (app *App) Initialise(initialTasks []Task, id int) {
	tasks = initialTasks
	currentID = id
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
}
func main() {
	app := App{}
	tasks, id := CreateInitialTasks()
	app.Initialise(tasks, id)
	app.Run("localhost:10000")
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	// your code goes here
}

func sendError(w http.ResponseWriter, statusCode int, err string) {
	// your code goes here
}

func (app *App) getTasks(writer http.ResponseWriter, request *http.Request) {
	// your code goes here
}

func (app *App) createTask(writer http.ResponseWriter, r *http.Request) {
	// your code goes here
}

func (app *App) readTask(writer http.ResponseWriter, request *http.Request) {
	// your code goes here
}

func (app *App) updateTask(writer http.ResponseWriter, request *http.Request) {
	// your code goes here
}

func (app *App) deleteTask(writer http.ResponseWriter, request *http.Request) {
	// your code goes here
}
