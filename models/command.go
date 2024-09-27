package models

import "time"

type ServiceModels struct {
	Id string `json:"id"`
	Content string `json:"content"`
	Action string `json:"action"`
	PreviousContent string `json:"previous_content"`
	CreatedDate time.Time `json:"created_date"`
	LastModified time.Time `json:"last_modified_date"`
	Title string `json:"title"`
	TypeAction string `json:"type_action"` 
	TypeActionName string `json:"type_action_name"` 
}