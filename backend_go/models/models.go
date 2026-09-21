package models

import (
	"time"
)

// -----------------Spots:get-spot-availability-------------------
type GetSpotAvailabilityResponse struct {
	Body struct {
		Message               string `json:"message" doc:"Status Message"`
		TotalSpotsTwoWheeler  int    `json:"total_spots_two_wheeler" doc:"Count of Total Two Wheeler Spots"`
		FreeSpotsTwoWheeler   int    `json:"free_spots_two_wheeler" doc:"Count of Free Two Wheeler Spots"`
		TotalSpotsFourWheeler int    `json:"total_spots_four_wheeler" doc:"Count of Total Four Wheeler Spots"`
		FreeSpotsFourWheeler  int    `json:"free_spots_four_wheeler" doc:"Count of Free Four Wheeler Spots"`
	}
}

// -----------------Spots:get-available-spots-------------------
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

// -----------------VehicleEntryExit:get-vehicle-------------------
type GetVehicleInput struct {
	LicensePlate string `path:"license_plate" doc:"License Plate to fetch details of vehicle"`
}

type VehicleDetails struct {
	OwnerID int    `json:"owner_id" doc:"OwnerID for this license plate"`
	Model   string `json:"model" doc:"Model of Vehicle for this license plate"`
	Colour  string `json:"colour" doc:"Colour of vehicle for this license plate"`
	Type    string `json:"vehicle_type" doc:"Type of Vehicle for this license plate"`
	Phone   string `json:"phone" doc:"Owner's phone number for this license plate"`
	Name    string `json:"name" doc:"Owner's name for this license plate"`
}

type GetVehicleResponse struct {
	Body struct {
		Message        string         `json:"message" doc:"User & Vehicle details of this license plate"`
		VehicleDetails VehicleDetails `json:"vehicle_details"`
	}
}

// -----------------VehicleEntryExit:save-vehicle-------------------
type SaveVehicleInput struct {
	LicensePlate   string         `json:"license_plate" doc:"License Number"`
	VehicleDetails VehicleDetails `json:"vehicle_details"`
}
type SaveVehicleResponse struct {
	Body struct {
		Message string `json:"message" doc:"Save vehicle with details to database"`
		Status  bool   `json:"status" doc:"Status of saving details to database"`
	}
}

// -----------------VehicleEntryExit:mark-entry-------------------
type MarkEntryInput struct {
	LicensePlate string `json:"license_plate" doc:"License Number"`
	CentreID     int    `json:"centre_id" doc:"Centre id currently selected"`
	Wing         string `json:"wing" doc:"Wing currently selected"`
	Floor        string `json:"floor" doc:"Floor curently selected"`
	SpotNumber   string `json:"spot_number" doc:"Spot Number currectly selected"`
	FolderPath   string `json:"folder_path" doc:"Folder Path to save capture"`
}

type MarkEntryResonse struct {
	Body struct {
		Message string `json:"message" doc:"Message"`
		Status  bool   `json:"status" doc:"Status"`
	}
}

// -----------------VehicleEntryExit:get-spot-details-------------------
type GetActiveSessionResponse struct {
	Body struct {
		Message    string    `json:"message" doc:"Latest session details"`
		EntryTime  time.Time `json:"entry_time" doc:"Latest entry time"`
		CentreID   string    `json:"centre_id" doc:"Centre-id of latest session"`
		Wing       string    `json:"wing" doc:"Wing of latest session"`
		Floor      string    `json:"floor" doc:"Floor of latest session"`
		SpotNumber string    `json:"spot_number" doc:"Spot Number of latest session"`
	}
}

type GetSpotDetailsResponse struct {
	Body struct {
		Message    string `json:"message" doc:"Latest session details"`
		SpotNumber string `json:"spot_number" doc:"Spot Number of latest session"`
		Status     bool   `json:"status" doc:"Status"`
	}
}

// -----------------VehicleEntryExit:mark-exit-------------------

type MarkExitInput struct {
	LicensePlate string    `json:"license_plate" doc:"License Number"`
	EntryTime    time.Time `json:"entry_time" doc:"Entry Time"`
	ExitTime     time.Time `json:"exit_time" doc:"Exit Time"`
	Duration     string    `json:"duration" doc:"(string)ExitTime - EntryTime"`
	Amount       float32   `json:"amount" doc:"Final Bill Amount"`
	CentreID     int       `json:"centre_id" doc:"Centre id currently selected"`
	Wing         string    `json:"wing" doc:"Wing currently selected"`
	Floor        string    `json:"floor" doc:"Floor curently selected"`
	SpotNumber   string    `json:"spot_number" doc:"Spot Number currectly selected"`
}

type MarkExitResponse struct {
	Body struct {
		Message   string    `json:"message" doc:"Message"`
		EntryTime time.Time `json:"entry_time" doc:"Entry Time"`
		ExitTime  time.Time `json:"exit_time" doc:"Exit Time"`
		Duration  string    `json:"duration" doc:"(string)ExitTime - EntryTime"`
		Amount    float32   `json:"amount" doc:"Final Bill Amount"`
	}
}

// -----------------Bills:get-latest-bill-------------------
// uses LicensePlate as Path
type GetLatestBillResponse struct {
	Body struct {
		Message   string    `json:"message" doc:"Message"`
		EntryTime time.Time `json:"entry_time" doc:"Entry Time"`
		ExitTime  time.Time `json:"exit_time" doc:"Exit Time"`
		Duration  string    `json:"duration" doc:"(string)ExitTime - EntryTime"`
		Amount    float32   `json:"amount" doc:"Final Bill Amount"`
		OwnerName string    `json:"owner_name" doc:"Owner's name"`
	}
}
