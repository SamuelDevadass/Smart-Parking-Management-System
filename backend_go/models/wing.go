package models

// Define your output structure
type WingsResponse struct {
	Body struct {
		Message string `json:"message"`
	}
}
