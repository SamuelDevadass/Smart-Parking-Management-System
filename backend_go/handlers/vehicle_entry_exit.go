package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"api.com/models"
	"api.com/repositories"
	"api.com/services"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterVehicleEntryExitHandler(api huma.API) {

	//-------------------Get-Vehicle------------------------------
	huma.Register(api, huma.Operation{
		OperationID: "get-vehicle",
		Method:      http.MethodGet,
		Path:        "/api/vehicles/{license_plate}",
	}, func(ctx context.Context, input *models.GetLicensePlate) (
		*models.GetVehicleResponse, error) {
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

	//-------------------Save-Vehicle-----------------------------
	huma.Register(api, huma.Operation{
		OperationID: "save-vehicle",
		Method:      http.MethodPost,
		Path:        "/api/vehicles",
	}, func(ctx context.Context, input *models.SaveVehicleInput) (
		*models.SaveVehicleResponse, error) {
		vehicle := &models.VehicleDetails{
			OwnerID: input.Body.OwnerID, Model: input.Body.Model,
			Colour: input.Body.Colour, Type: input.Body.Type,
			Phone: input.Body.Phone, Name: input.Body.Name,
		}

		status, err := repositories.SaveVehicle(ctx, vehicle, input.Body.LicensePlate)
		if err != nil || status != true {
			log.Println("Failed to save vehicle details \n", err)
			return nil, huma.Error500InternalServerError("Failed to save vehicle details")
		}
		resp := &models.SaveVehicleResponse{}
		resp.Body.Message = fmt.Sprintf("Vehicle Details for license plate '%s'", input.Body.LicensePlate)
		resp.Body.Status = true
		return resp, nil
	})

	//-------------------Mark-Entry-----------------------------
	huma.Register(api, huma.Operation{
		OperationID: "mark-entry",
		Method:      http.MethodPost,
		Path:        "/api/entries",
	}, func(ctx context.Context, input *models.MarkEntryInput) (
		*models.MarkEntryResonse, error) {
		now := time.Now()
		status, err := services.GetDetectionStatus()
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to fetch detection status")
		}
		input.Body.FolderPath = services.CreateFolderPath(input.Body.LicensePlate, status.FolderPath)
		status_py, err_py := repositories.MarkEntry(ctx, input, now)
		if err_py != nil || status_py != true {
			log.Println("Failed to mark entry \n", err_py)
			return nil, huma.Error500InternalServerError("Failed to mark entry")
		}
		resp := &models.MarkEntryResonse{}
		resp.Body.Message = fmt.Sprintf("Vehicle Details for license plate '%s'", input.Body.LicensePlate)
		resp.Body.Status = true
		return resp, nil
	})

	//-------------------Get-Spot-Details-----------------------------
	huma.Register(api, huma.Operation{
		OperationID: "get-spot-details",
		Method:      http.MethodGet,
		Path:        "/api/vehicles/spot/{license_plate}",
	}, func(ctx context.Context, input *models.GetLicensePlate) (
		*models.GetSpotDetailsResponse, error) {
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

	//-------------------Mark-Exit-----------------------------
	huma.Register(api, huma.Operation{
		OperationID: "mark-exit",
		Method:      http.MethodPut,
		Path:        "/api/exits",
	}, func(ctx context.Context, input *models.MarkExitInput) (
		*models.MarkExitResponse, error) {
		session, err := repositories.GetActiveSession_NoPath(ctx, input)
		if err != nil {
			log.Println("Failed to fetch active session \n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch session details")
		}
		if session == nil {
			return nil, huma.Error404NotFound(
				fmt.Sprintf("No active session found for license plate '%s'", input.Body.LicensePlate),
			)
		}
		vehicleType, err := repositories.GetvehicleType(ctx, input.Body.LicensePlate)
		if err != nil {
			log.Println("Failed to fetch vehicle type \n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch vehicle type")
		}
		if vehicleType == "" {
			return nil, huma.Error404NotFound(
				fmt.Sprintf("No vehicle type found for license plate '%s'", input.Body.LicensePlate),
			)
		}
		exitTime := time.Now()
		entryTime := session.Body.EntryTime
		duration := exitTime.Sub(entryTime)
		amount := services.CalculateBillAmount(duration.Seconds(), vehicleType)
		details := &models.RecordExitDetails{
			EntryTime: entryTime, ExitTime: exitTime,
			Duration: duration.String(), Amount: float32(amount),
			CentreID: session.Body.CentreID, Wing: session.Body.Wing,
			Floor: session.Body.Floor, SpotNumber: session.Body.SpotNumber,
		}
		status, err := repositories.RecordExit(ctx, details)
		if err != nil || !status {
			log.Println("Failed to record exit \n", err)
			return nil, huma.Error500InternalServerError("Failed to record exit")
		}
		resp := &models.MarkExitResponse{}
		resp.Body.Message = "Parking session ended successfully"
		resp.Body.EntryTime = entryTime
		resp.Body.ExitTime = exitTime
		resp.Body.Duration = duration.String()
		resp.Body.Amount = float32(amount)
		return resp, nil
	})
}
