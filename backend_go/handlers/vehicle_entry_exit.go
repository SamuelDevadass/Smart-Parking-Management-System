package handlers

import (
	"context"
	"fmt" // Added for missing case error checks
	"log"
	"net/http"
	"time"

	"api.com/models"
	"api.com/repositories"
	"api.com/services"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterVehicleEntryExitHandler(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-vehicle",
		Method:      http.MethodGet,
		Path:        "/api/vehicles/{license_plate}",
	}, func(ctx context.Context, input *models.GetVehicleInput) (*models.GetVehicleResponse, error) {
		vehicle, err := repositories.GetVehicle(ctx, input.LicensePlate)
		if err != nil {
			log.Println("Failed to fetch vehicle details \n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch vehicle details")
		}
		if vehicle == nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("No vehicle details found for license plate '%s'", input.LicensePlate))
		}
		resp := &models.GetVehicleResponse{}
		resp.Body.Message = fmt.Sprintf("Vehicle Details for license plate '%s'", input.LicensePlate)
		resp.Body.VehicleDetails = *vehicle
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "save-vehicle",
		Method:      http.MethodPut,
		Path:        "/api/vehicles",
	}, func(ctx context.Context, input *models.SaveVehicleInput) (*models.SaveVehicleResponse, error) {
		status, err := repositories.SaveVehicle(ctx, &input.VehicleDetails, input.LicensePlate)
		if err != nil || status != true {
			log.Println("Failed to save vehicle details \n", err)
			return nil, huma.Error500InternalServerError("Failed to save vehicle details")
		}
		resp := &models.SaveVehicleResponse{}
		resp.Body.Message = fmt.Sprintf("Vehicle Details for license plate '%s'", input.LicensePlate)
		resp.Body.Status = true
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "mark-entry",
		Method:      http.MethodPut,
		Path:        "/api/entries",
	}, func(ctx context.Context, input *models.MarkEntryInput) (*models.MarkEntryResonse, error) {
		now := time.Now()
		status, err := repositories.MarkEntry(ctx, input, now)
		if err != nil || status != true {
			log.Println("Failed to mark entry \n", err)
			return nil, huma.Error500InternalServerError("Failed to mark entry")
		}
		resp := &models.MarkEntryResonse{}
		resp.Body.Message = fmt.Sprintf("Vehicle Details for license plate '%s'", input.LicensePlate)
		resp.Body.Status = true
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-spot-details",
		Method:      http.MethodGet,
		Path:        "/api/vehicles/spot/{license_plate}",
	}, func(ctx context.Context, input *models.GetVehicleInput) (*models.GetSpotDetailsResponse, error) {
		ans, err := repositories.GetActiveSession(ctx, input)
		if err != nil {
			log.Println("Failed to mark entry \n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch session details")
		}
		if ans == nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("No session details found for license plate '%s'", input.LicensePlate))
		}
		resp := &models.GetSpotDetailsResponse{}
		resp.Body.Message = fmt.Sprintf("Session details for License Plate '%s'", input.LicensePlate)
		resp.Body.SpotNumber = ans.Body.SpotNumber
		resp.Body.Status = true
		return resp, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "mark-exit",
		Method:      http.MethodPost,
		Path:        "/api/exits",
	}, func(ctx context.Context, input *models.MarkExitInput) (*models.MarkExitResponse, error) {
		session, err := repositories.GetActiveSession_NoPath(ctx, input)
		if err != nil {
			log.Println("Failed to fetch active session \n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch session details")
		}
		if session == nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("No active session found for license plate '%s'", input.LicensePlate))
		}

		vehicleType, err := repositories.GetvehicleType(ctx, input.LicensePlate)
		if err != nil {
			log.Println("Failed to fetch vehicle type \n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch vehicle type")
		}
		if vehicleType == "" {
			return nil, huma.Error404NotFound(fmt.Sprintf("No vehicle type found for license plate '%s'", input.LicensePlate))
		}

		exitTime := time.Now()
		entryTime := session.Body.EntryTime // Adjust based on your session model structure

		// Go-specific time handling
		duration := exitTime.Sub(entryTime)
		amount := services.CalculateBillAmount(duration.Seconds(), vehicleType)

		status, err := repositories.RecordExit(ctx, input)
		if err != nil || !status {
			log.Println("Failed to record exit \n", err)
			return nil, huma.Error500InternalServerError("Failed to record exit")
		}

		// Populate the response matching your MarkExitResponse model
		resp := &models.MarkExitResponse{}
		resp.Body.Message = "Parking session ended successfully"
		resp.Body.EntryTime = entryTime
		resp.Body.ExitTime = exitTime
		resp.Body.Duration = duration.String() // Equivalent to Python's str(duration)
		resp.Body.Amount = float32(amount)     // Convert float64 to float32 for your model

		return resp, nil
	})
}
