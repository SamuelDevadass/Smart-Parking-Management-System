package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"api.com/models"
	"github.com/danielgtaylor/huma/v2"
)

const pythonBaseURL = "http://localhost:8001"

func RegisterDetectionHandler(api huma.API) {
	// 1. POST /api/detection/start
	huma.Register(api, huma.Operation{
		OperationID: "start-detection",
		Method:      http.MethodPost,
		Path:        "/api/detection/start",
	}, func(ctx context.Context, input *struct{}) (*models.StartDetectionResponse, error) {
		resp, err := http.Post(pythonBaseURL+"/api/detection/start", "application/json", nil)
		if err != nil {
			log.Println("Failed to start deteection\n", err)
			return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to reach OCR service: %v", err))
		}
		defer resp.Body.Close()

		var pyResp struct {
			Ok bool `json:"ok"`
		}
		json.NewDecoder(resp.Body).Decode(&pyResp)

		response := &models.StartDetectionResponse{}
		response.Body.Ok = pyResp.Ok
		return response, nil
	})

	// 2. GET /api/detection/status
	huma.Register(api, huma.Operation{
		OperationID: "detection-status",
		Method:      http.MethodGet,
		Path:        "/api/detection/status",
	}, func(ctx context.Context, input *struct{}) (*models.DetectionStatusResponse, error) {
		resp, err := http.Get(pythonBaseURL + "/api/detection/status")
		if err != nil {
			log.Println("Unable to start OCR Engine\n", err)
			return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to reach OCR service: %v", err))
		}
		defer resp.Body.Close()

		var statusData struct {
			Status       string `json:"status"`
			LicensePlate string `json:"license_plate"`
			FolderPath   string `json:"folder_path"`
		}
		json.NewDecoder(resp.Body).Decode(&statusData)

		response := &models.DetectionStatusResponse{}
		response.Body.Status = statusData.Status
		response.Body.LicensePlate = statusData.LicensePlate
		response.Body.FolderPath = statusData.FolderPath
		return response, nil
	})

	// 3. POST /api/detection/stop
	huma.Register(api, huma.Operation{
		OperationID: "stop-detection",
		Method:      http.MethodPost,
		Path:        "/api/detection/stop",
	}, func(ctx context.Context, input *struct{}) (*models.StopDetectionResponse, error) {
		resp, err := http.Post(pythonBaseURL+"/api/detection/stop", "application/json", nil)
		if err != nil {
			return nil, huma.Error500InternalServerError(fmt.Sprintf("Failed to reach OCR service: %v", err))
		}
		defer resp.Body.Close()

		var pyResp struct {
			Ok bool `json:"ok"`
		}
		json.NewDecoder(resp.Body).Decode(&pyResp)

		response := &models.StopDetectionResponse{}
		response.Body.Ok = pyResp.Ok
		return response, nil
	})
}
