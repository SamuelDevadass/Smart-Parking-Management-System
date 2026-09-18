package models

type GetSpotAvailabilityResponse struct {
	Body struct {
		Message               string `json:"message" doc:"Status Message"`
		TotalSpotsTwoWheeler  int    `json:"total_spots_two_wheeler" doc:"Count of Total Two Wheeler Spots"`
		FreeSpotsTwoWheeler   int    `json:"free_spots_two_wheeler" doc:"Count of Free Two Wheeler Spots"`
		TotalSpotsFourWheeler int    `json:"total_spots_four_wheeler" doc:"Count of Total Four Wheeler Spots"`
		FreeSpotsFourWheeler  int    `json:"free_spots_four_wheeler" doc:"Count of Free Four Wheeler Spots"`
	}
}

type GetAvailableSpotsInput struct {
	Wing      string `path:"wing" doc:"Wing name"`
	Centre_id int    `query:"centre_id" doc:"Centre-id of wing"`
	Size      string `query:"size" doc:"Size of Vehicle-Two Wheeler or Four Wheeler"`
}

type GetAvailableSpotsResponse struct {
	Body struct {
		Message            string              `json:"message" doc:"List of all available spots in format floor: <>, spot_number: <>, size: <>"`
		AvailableSpotsList []map[string]string `json:"available_spots_list" doc:"return list of all available/free spots in current wing"`
	}
}
