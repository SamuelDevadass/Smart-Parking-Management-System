package handlers

import (
	"context"
	"fmt" // Added for missing case error checks
	"log"
	"net/http"
	"time"

	"api.com/models"
	"api.com/repositories"
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
}
