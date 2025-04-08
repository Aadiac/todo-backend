package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Email    string     `json:"email" bson:"email"`
	Password string     `json:"password" bson:"password"`
	Task     []TaskList `json:"task" bson:"task"`
}
type TaskList struct {
	TaskId    primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Created   string `json:"createdDate" bson:"createdDate"`
	Task      string `json:"task" bson:"task"`
	Completed string    `json:"completedBy" bson:"completedBy"`
}

// for decoding data from user while login
type UiRequset struct {
	Email    string     `json:"email" bson:"email"`
	Password string     `json:"password" bson:"password"`
	Task     []TaskList `json:"task" bson:"task"`
}
//for decoding Edittask data 
type EditTaskRequest struct{
	TaskId    primitive.ObjectID  `json:"TaskId,omitempty"`
	Created   string 	`json:"createdDate" bson:"createdDate"`
	Task      string 	`json:"task" bson:"task"`
	Completed string    `json:"completedBy" bson:"completedBy"`
	Email    string     `json:"email" bson:"email"`

}
