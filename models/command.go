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

type Athlete struct {
    ID      string `json:"id"`
    Nom     string `json:"nom"`
    Prenom  string `json:"prenom"`
    Age     string `json:"age,omitempty"`
}

type AthleteFile struct {
	AthleteID string `json:"athlete_id"`
	File 		string 	 `json:"file,omitempty"`
	TypeFile string  `json:"type_file,omitempty"`
	LinkFile 	string `json:"link_file,omitempty"`
}