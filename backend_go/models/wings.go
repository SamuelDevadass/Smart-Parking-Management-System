package models

// Define your output structure
//highest level has to be a struct
// --------------------Wings:get-wings-----------------------
type GetWingsResponse struct {
	Body struct {
		Message string   `json:"message" doc:"Status Message"`
		Wings   []string `json:"wings" doc:"List of distinct parking wings"`
	}
}

// --------------------Wings:get-centre-id-for-wing-----------------------
type WingPathInput struct {
	Wing string `path:"wing" doc:"Wing name"`
}
type GetCentreForWingsResponse struct {
	Body struct {
		Message string `json:"message" doc:"Status Message"`
		Centre  int    `json:"centre_id" doc:"Centre-id for current wing"`
	}
}
