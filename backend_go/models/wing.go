package models

// Define your output structure
//highest level has to be a struct
type WingsResponse struct {
	Body struct {
		Message string   `json:"message" doc:"Status Message"`
		Wings   []string `json:"wings" doc:"List of distinct parking wings"`
	}
}
