package models

// Define your output structure
//highest level has to be a struct
type GetWingsResponse struct {
	Body struct {
		Message string   `json:"message" doc:"Status Message"`
		Wings   []string `json:"wings" doc:"List of distinct parking wings"`
	}
}
type WingInput struct {
	Wing string `query:"wing" doc:"Wing name to get its centre id"`
}

type WingPathInput struct {
	Wing string `path:"wing" doc:"Wing name"`
}
type GetCentreForWingsResponse struct {
	Body struct {
		Message string `json:"message" doc:"Status Message"`
		Centre  int    `json:"centre_id" doc:"Centre-id for current wing"`
	}
}
