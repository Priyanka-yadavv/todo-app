package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app App

func CreateInitialTestTasks() ([]Task, int) {
	tasks := []Task{
		{ID: 1, Name: "Create project plan", Description: "Outline the goals and objectives for the project", DueDate: "2022-01-01"},
		{ID: 2, Name: "Gather requirements", Description: "Conduct research and gather input from stakeholders", DueDate: "2022-02-01"},
		{ID: 3, Name: "Design solution", Description: "Create a detailed design of the proposed solution", DueDate: "2022-03-01"},
	}
	return tasks, 3
}

func TestMain(m *testing.M) {
	tasks, id := CreateInitialTestTasks()
	app.Initialise(tasks, id)
	m.Run()
}

func TestGetAllTasks(t *testing.T) {
	allTasks, _ := getTasks()
	request, _ := http.NewRequest("GET", "/tasks", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var tasks []Task
	json.Unmarshal(response.Body.Bytes(), &tasks)

	if len(tasks) != len(allTasks) {
		t.Errorf("handler returned unexpected number of tasks: got %v want %v",
			len(tasks), len(allTasks))
	}
	for i := range tasks {
		if tasks[i].ID != allTasks[i].ID {
			t.Errorf("handler returned unexpected ID for task %d: got %v want %v", i+1, tasks[i].ID, allTasks[i].ID)
		}
		if tasks[i].Name != allTasks[i].Name {
			t.Errorf("handler returned unexpected Name for task %d: got %v want %v", i+1, tasks[i].Name, allTasks[i].Name)
		}
		if tasks[i].Description != allTasks[i].Description {
			t.Errorf("handler returned unexpected Description for task %d: got %v want %v", i+1, tasks[i].Description, allTasks[i].Description)
		}
		if tasks[i].DueDate != allTasks[i].DueDate {
			t.Errorf("handler returned unexpected DueDate for task %d: got %v want %v", i+1, tasks[i].DueDate, allTasks[i].DueDate)
		}
	}
}

func TestReadTask(t *testing.T) {
	// Test with valid task ID
	request, _ := http.NewRequest("GET", "/task/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var task Task
	json.Unmarshal(response.Body.Bytes(), &task)

	expectedTask := tasks[0]
	if task.ID != expectedTask.ID {
		t.Errorf("handler returned unexpected ID for task: got %v want %v", task.ID, expectedTask.ID)
	}
	if task.Name != expectedTask.Name {
		t.Errorf("handler returned unexpected Name for task: got %v want %v", task.Name, expectedTask.Name)
	}
	if task.Description != expectedTask.Description {
		t.Errorf("handler returned unexpected Description for task: got %v want %v", task.Description, expectedTask.Description)
	}
	if task.DueDate != expectedTask.DueDate {
		t.Errorf("handler returned unexpected DueDate for task: got %v want %v", task.DueDate, expectedTask.DueDate)
	}

	// Test with non-existent task ID
	request, _ = http.NewRequest("GET", "/task/10", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)

	var errorMessage map[string]string
	json.Unmarshal(response.Body.Bytes(), &errorMessage)

	if errorMessage["error"] != "task not found" {
		t.Errorf("handler returned unexpected error message: got %v want %v", errorMessage["error"], "task not found")
	}

	// Test with invalid task ID format
	request, _ = http.NewRequest("GET", "/task/abc", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusBadRequest, response.Code)

	json.Unmarshal(response.Body.Bytes(), &errorMessage)

	if errorMessage["error"] != "invalid task ID" {
		t.Errorf("handler returned unexpected error message: got %v want %v", errorMessage["error"], "invalid task ID")
	}
}

func TestCreateTask(t *testing.T) {
	// Test with valid task
	task := Task{Name: "Test task", Description: "This is a test task", DueDate: "2022-01-01"}
	payload, _ := json.Marshal(task)
	request, _ := http.NewRequest("POST", "/task", bytes.NewBuffer(payload))
	response := sendRequest(request)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var returnedTask Task
	json.Unmarshal(response.Body.Bytes(), &returnedTask)

	task.ID = currentID
	if returnedTask.ID != task.ID {
		t.Errorf("handler returned unexpected ID for task: got %v want %v", returnedTask.ID, task.ID)
	}
	if returnedTask.Name != task.Name {
		t.Errorf("handler returned unexpected Name for task: got %v want %v", returnedTask.Name, task.Name)
	}
	if returnedTask.Description != task.Description {
		t.Errorf("handler returned unexpected Description for task: got %v want %v", returnedTask.Description, task.Description)
	}
	if returnedTask.DueDate != task.DueDate {
		t.Errorf("handler returned unexpected DueDate for task: got %v want %v", returnedTask.DueDate, task.DueDate)
	}

	// Test with invalid task payload
	request, _ = http.NewRequest("POST", "/task", bytes.NewBuffer([]byte("invalid payload")))
	response = sendRequest(request)
	checkStatusCode(t, http.StatusBadRequest, response.Code)

	var errorResponse map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &errorResponse)

	if errorResponse["error"] != "Invalid request payload" {
		t.Errorf("handler returned unexpected error message: got %v want %v", errorResponse["error"], "Invalid request payload")
	}
}

func TestUpdateTask(t *testing.T) {
	// Test with valid task ID and payload
	task := tasks[1]
	task.Name = "Collect requirements"
	payload, _ := json.Marshal(task)
	request, _ := http.NewRequest("PUT", "/task/2", bytes.NewBuffer(payload))
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var returnedTask Task
	json.Unmarshal(response.Body.Bytes(), &returnedTask)

	if returnedTask.ID != task.ID {
		t.Errorf("handler returned unexpected ID for task: got %v want %v", returnedTask.ID, task.ID)
	}
	if returnedTask.Name != task.Name {
		t.Errorf("handler returned unexpected Name for task: got %v want %v", returnedTask.Name, task.Name)
	}

	// Test with invalid task ID
	request, _ = http.NewRequest("PUT", "/task/abc", bytes.NewBuffer(payload))
	response = sendRequest(request)
	checkStatusCode(t, http.StatusBadRequest, response.Code)

	var errorResponse map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &errorResponse)

	if errorResponse["error"] != "invalid task ID" {
		t.Errorf("handler returned unexpected error message: got %v want %v", errorResponse["error"], "invalid task ID")
	}

	// Test with invalid task payload
	request, _ = http.NewRequest("PUT", "/task/1", bytes.NewBuffer([]byte("invalid payload")))
	response = sendRequest(request)
	checkStatusCode(t, http.StatusBadRequest, response.Code)

	json.Unmarshal(response.Body.Bytes(), &errorResponse)

	if errorResponse["error"] != "Invalid request payload" {
		t.Errorf("handler returned unexpected error message: got %v want %v", errorResponse["error"], "Invalid request payload")
	}
}

func TestDeleteTask(t *testing.T) {
	// Test with valid task ID
	request, _ := http.NewRequest("DELETE", "/task/3", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)

	var responseBody map[string]string
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	if responseBody["result"] != "successful deletion" {
		t.Errorf("handler returned unexpected response: got %v want %v", responseBody["result"], "successful deletion")
	}

	// Test with non-existing task ID
	request, _ = http.NewRequest("DELETE", "/task/10", nil)
	response = sendRequest(request)
	checkStatusCode(t, http.StatusNotFound, response.Code)

	var errorResponse map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &errorResponse)

	if errorResponse["error"] != "task not found" {
		t.Errorf("handler returned unexpected error message: got %v want %v", errorResponse["error"], "task not found")
	}
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, request)
	return recorder
}
