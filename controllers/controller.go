package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	model "github.com/aadiac/todo/models"
	utils "github.com/aadiac/todo/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://aadithmongodb:mongodb@cluster0.710qp.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const dbname = "DOdb"
const collectionName = "todo"

var SessionCtx model.User 

var Collection *mongo.Collection

func init() {
	fmt.Println("controllers.init() >>>")
	opt := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), opt)

	if err != nil {
		log.Fatal(err)

	}

	Collection = client.Database(dbname).Collection(collectionName)
	fmt.Println("controllers.init() <<<")
}

// ================HELPER METHODS==================
func insertLoginCreds(user model.UiRequset) {
	fmt.Println("controllers.insertLoginCreds() >>>")

	Inserted, err := Collection.InsertOne(context.Background(), user)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted Login Credential : ", Inserted.InsertedID)

	fmt.Println("controllers.insertLoginCreds() <<<")
}

func FindOne(email string) model.User {
	fmt.Println("controllers.FindOne() >>>")

	var fetchData model.User
	filter := bson.M{"email": email}

	if err := Collection.FindOne(context.Background(), filter).Decode(&fetchData); err != nil {
		fmt.Println("Err in FindOne() : ", err)
	}
	// fmt.Printf("Returning Fetch Data from Db  : %+v\n ", fetchData)

	fmt.Println("controllers.FindOne() <<<")

	return fetchData

}

func updateTask(updateReq model.EditTaskRequest, id primitive.ObjectID, isTaskDel bool) error {
	fmt.Println("controllers.updateTask() >>>")

	var update bson.M
	var operation string
	filter := bson.M{"task._id": id, "email": updateReq.Email}

	if !isTaskDel {
		update = bson.M{"$set": bson.M{"task.$.createdDate": time.Now().Format("2006-01-02"),
			"task.$.task":        updateReq.Task,
			"task.$.completedBy": updateReq.Completed,
		}}
		operation = "update"
	}else{
		update = bson.M{"$pull":bson.M{"task": bson.M{"_id":id}}}
		operation = "delete Task"
	}
	updateRes, err := Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if updateRes.ModifiedCount == 0 {
		return fmt.Errorf("task not found,no change made")
	}
	fmt.Println(operation,"task Result : ", updateRes.ModifiedCount)

	fmt.Println("controllers.updateTask() <<<")

	return nil
}

func appendTask(task []model.TaskList, email string) error {
	fmt.Println("controller.appendTask() >>>")

	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"task": task}}

	updateRes, err := Collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}
	if updateRes.ModifiedCount == 0 {
		return fmt.Errorf("task not appended(modified)")
	}

	fmt.Println("appendTask Result : ", updateRes.ModifiedCount)

	fmt.Println("controller.appendTask() <<<")
	return nil
}

//=================================================END OF HELPER METHODS

func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controllers.Register() >>>")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var myuser model.UiRequset
	json.NewDecoder(r.Body).Decode(&myuser)

	if myuser.Email == "" || myuser.Password == "" {
		fmt.Println("Email or Passowrd is empty")
		http.Error(w, "Email | Password is required", http.StatusBadRequest) //400

		return
	}

	//email is already registered or not
	var alreadyExist model.User
	filter := bson.M{"email": myuser.Email}
	if err := Collection.FindOne(context.Background(), filter).Decode(&alreadyExist); err == nil {
		fmt.Println("Email already registered")
		http.Error(w, "Email already Registerd", http.StatusConflict) //409
		return
	} else {
		fmt.Println("Regster  ,ERR : ", err)
	}
	//============================

	var mytask model.TaskList
	if myuser.Task == nil {
		mytask.TaskId = primitive.NewObjectIDFromTimestamp(time.Now())
		mytask.Created = time.Now().Format("2006-01-02")
		mytask.Task = "Task"
		mytask.Completed = "1"

		myuser.Task = append(myuser.Task, mytask)
	}
	insertLoginCreds(myuser)
	json.NewEncoder(w).Encode(myuser)

	fmt.Println("controllers.Register() <<<")
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controllers.Login() >>>")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)

		return
	}

	var inputData model.UiRequset
	json.NewDecoder(r.Body).Decode(&inputData)
	// fmt.Println("data :", inputData.Password, inputData.Email) //DEBUG

	filter := bson.M{"email": inputData.Email}
	var storedData model.User
	if err := Collection.FindOne(context.Background(), filter).Decode(&storedData); err == nil {

		if inputData.Password != storedData.Password {
			fmt.Println("Invalid password")
			http.Error(w, "Incorrect Password  , Try Again !", http.StatusUnauthorized) //401
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Login Successfull")
	} else {
		http.Error(w, "Incorrect email-id", http.StatusUnauthorized) //401
		return
	}
	fmt.Println("controllers.Login() <<<")
}

func deleteOneTask(taskId string) error {
	fmt.Println("controllers.deleteOneTask >>>")
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}

	delCount, _ := Collection.DeleteOne(context.Background(), filter)
	if delCount.DeletedCount != 0 {
		return fmt.Errorf("error in deleteing task")
	}
	fmt.Println("controllers.deleteOneTask <<<")
	return nil

}

func FetchTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controllers.FetchTasks() >>>")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	params := mux.Vars(r)

	var storedData model.User
	// fmt.Println("email : ", params["email"])
	if storedData = FindOne(params["email"]); storedData.Task == nil {
		fmt.Printf("stredData : %+v\n", storedData)
		http.Error(w, "No task are available for you", 404)
		return
	}

	
	SessionCtx = storedData
	fmt.Printf("Session Context : %+v\n",SessionCtx)

	json.NewEncoder(w).Encode(&storedData.Task)
	fmt.Println("controllers.FetchTasks() <<<")
}

func EditTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controllers.EditTask() >>>")

	utils.PreflightCheck(w, r)
	//  utils.PrintRawData(r)

	var inputData model.EditTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		fmt.Println("Err in decoding data: ", err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Type of Completeed BY : %T", inputData.Created)
	storedData := FindOne(inputData.Email) //real data from db

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["taskid"])
	utils.ErrHandle(err)

	fmt.Println("inputdata tskid : ", inputData.TaskId, " url id : ", id)

	//handling if user press save if there is nothing to save
	for index, task := range storedData.Task {
		if task.TaskId == id {
			if storedData.Task[index].Completed == inputData.Completed && storedData.Task[index].Created == time.Now().Format("2006-01-02") && storedData.Task[index].Task == inputData.Task {
				json.NewEncoder(w).Encode(storedData)
			}
		}
	}
	fmt.Printf("after creating new updated data : %+v\n", storedData)

	// var updatedData []model.TaskList
	if err := updateTask(inputData, id, false); err != nil {
		http.Error(w, "Error In updating Db", http.StatusInternalServerError)
		fmt.Println("Err in Updating DB ", err)
		return
	}
	result := FindOne(inputData.Email)
	fmt.Printf("Task Updated Succesfull ! : %+v\n", result)
	json.NewEncoder(w).Encode(result)

	fmt.Println("controllers.FetchTasks() <<<")

}

func AddNewTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controllers.AddNewTask()>>>")

	utils.PreflightCheck(w, r)

	// utils.PrintRawData(r)
	var userInputData model.EditTaskRequest // same structure type of data are we used here

	if err := json.NewDecoder(r.Body).Decode(&userInputData); err != nil {
		fmt.Println("Error in decoding AddnewTask ", err)
		http.Error(w, "Invalid Body", http.StatusBadRequest)
		return
	}

	storedData := FindOne(userInputData.Email)

	var newTask model.TaskList
	newTask.TaskId = primitive.NewObjectIDFromTimestamp(time.Now())
	newTask.Created = userInputData.Created
	newTask.Task = userInputData.Task
	newTask.Completed = userInputData.Completed

	storedData.Task = append(storedData.Task, newTask)

	if err := appendTask(storedData.Task, userInputData.Email); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := FindOne(userInputData.Email)
	fmt.Printf("Updated SuccesFully : %+v\n", result)

	json.NewEncoder(w).Encode(result)
	// w.WriteHeader(http.StatusOK)

	fmt.Println("controllers.AddNewTask()<<<")
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("controllers.DeleteTask()>>>")
	utils.PreflightCheck(w, r)

	 utils.PrintRawData(r)
	var request model.EditTaskRequest
	request.Email = SessionCtx.Email
	// if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
	// 	fmt.Println("Error in decoding DeleteTask Req ", err)
	// 	http.Error(w, "Invalid Body", http.StatusBadRequest)
	// 	return
	// }
	// fmt.Println("Request email: ", request.Email)
	//

	
	param := mux.Vars(r)
	id,err := primitive.ObjectIDFromHex(param["taskid"])
	if err != nil{
		fmt.Println("Error in decoding taskid ", err)
		http.Error(w, "decoding Taskid Failed", http.StatusBadRequest)
		return
	}

	if err := updateTask(request,id,true); err != nil {
		fmt.Println("Err in Deleting Task:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Println("Task Deleted Succesfully ")
	w.WriteHeader(http.StatusNoContent)

	fmt.Println("controllers.DeleteTask()<<<")

}
