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
