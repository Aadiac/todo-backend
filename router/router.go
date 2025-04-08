package router

import (
	controller "github.com/aadiac/todo/controllers"
	"github.com/gorilla/mux"
)

func Router() (*mux.Router){
	r := mux.NewRouter()
	
	r.HandleFunc("/register",controller.Register).Methods("POST")
	r.HandleFunc("/login",controller.Login).Methods("POST")
	r.HandleFunc("/tasks/{email}",controller.FetchTasks).Methods("GET")//display tasks 
	r.HandleFunc("/updatetask/{taskid}",controller.EditTask).Methods("PUT")//update particular task
	r.HandleFunc("/addnewtask",controller.AddNewTask).Methods("POST")
	
	return r
}
