package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Task struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Is_completed bool   `json:"is_completed"`
}

type BulkTask struct {
	BulkTask []Task `json:"tasks"`
}

var id rune = '0'
var tasks = allTask{}

type allTask []Task

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createTask(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)

	id = rune(int(id) + 1)
	if err != nil {
		fmt.Println("Kindly enter valid data in order to:", err)
		return
	}
	newTask := Task{}
	json.Unmarshal(reqBody, &newTask)

	newTask.Id = string(id)

	tasks = append(tasks, newTask)
	fmt.Fprintf(w, "{id: %v } ", newTask.Id)
	w.WriteHeader(http.StatusCreated)
}

func bulkTaskCreation(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	newData := BulkTask{}
	json.Unmarshal([]byte(reqBody), &newData)
	tasks = append(tasks, newData.BulkTask...)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newData)
}

func getAllTask(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]
	var updatedTask Task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter valid data in order to update")
	}
	json.Unmarshal(reqBody, &updatedTask)
	for i, singleTask := range tasks {
		if singleTask.Id == taskId {
			singleTask.Title = updatedTask.Title
			singleTask.Is_completed = updatedTask.Is_completed
			tasks[i] = singleTask
			json.NewEncoder(w).Encode(singleTask)
			return
		}
	}
	json.NewEncoder(w).Encode(http.StatusNotFound)

}

// func bulkTaskDeletion(w http.ResponseWriter, r *http.Request) {
// 	reqBody, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		fmt.Println("Error reading JSON data:", err)
// 		return
// 	}
// 	newData := BulkTask{}
// 	json.Unmarshal([]byte(reqBody), &newData)
// 	for i, singleTask := range tasks {
// 		for uniqueTaskId := range newData.BulkTask {

// 			if singleTask.Id == newData.BulkTask[uniqueTaskId].Id {
// 				println(singleTask.Id, newData.BulkTask[uniqueTaskId].Id)
// 				print(i, j)
// 				// tasks = append(tasks[:i], tasks[i+1:]...)
// 			}
// 		}

// 	}
// }

func deleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]

	for i, singleTask := range tasks {
		if singleTask.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			json.NewEncoder(w).Encode(singleTask)
			return
		}

	}
	json.NewEncoder(w).Encode(http.StatusNotFound)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home)
	router.HandleFunc("/v1/tasks", createTask).Methods("POST")
	router.HandleFunc("/v2/tasks", bulkTaskCreation).Methods("POST")
	router.HandleFunc("/v1/tasks", getAllTask).Methods("GET")
	router.HandleFunc("/v1/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/v1/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/v2/tasks", bulkTaskCreation).Methods("POST")
	router.HandleFunc("/v1/tasks", bulkTaskDeletion).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

}
