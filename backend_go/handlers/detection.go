package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"api.com/models"
	"api.com/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

// const pythonBaseURL = "http://localhost:8001"

func RegisterDetectionHandler(r chi.Router, api huma.API, python_url string) {

	//-------------------Start-Detection-----------------------------

	/*huma.Register(api, huma.Operation{
		OperationID: "start-detection",
		Method:      http.MethodPost,
		Path:        "/api/detection/start",
	}, func(ctx context.Context, input *struct{}) (*models.StartDetectionResponse, error) {
		// Read the JPEG bytes sent by React.
		body, err := io.ReadAll(ctx.Value(huma.ContextKey).(*http.Request).Body)
		if err != nil {
			return nil, huma.Error400BadRequest("Failed to read image body")
		}

		// Forward the same JPEG bytes to FastAPI.
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			python_url+"/api/detection/start",
			bytes.NewReader(body),
		)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to create OCR request")
		}

		req.Header.Set("Content-Type", "image/jpeg")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, huma.Error500InternalServerError(
				fmt.Sprintf("Failed to reach OCR service: %v", err),
			)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, huma.Error500InternalServerError(
				fmt.Sprintf("OCR service returned HTTP %d", resp.StatusCode),
			)
		}

		var pyResp struct {
			Ok bool `json:"ok"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&pyResp); err != nil {
			return nil, huma.Error500InternalServerError("Invalid response from OCR service")
		}

		response := &models.StartDetectionResponse{}
		response.Body.Ok = pyResp.Ok
		return response, nil
	})*/
	r.Post("/api/detection/start", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		resp, err := http.NewRequestWithContext(
			req.Context(),
			http.MethodPost,
			python_url+"/api/detection/start",
			req.Body,
		)
		if err != nil {
			http.Error(w, "Failed to create OCR request", http.StatusInternalServerError)
			return
		}

		resp.Header.Set("Content-Type", "image/jpeg")

		upstream, err := http.DefaultClient.Do(resp)
		if err != nil {
			http.Error(w, "Failed to reach OCR service", http.StatusBadGateway)
			return
		}
		defer upstream.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(upstream.StatusCode)
		_, _ = io.Copy(w, upstream.Body)
	})

	//-------------------Detection-Status-----------------------------
	huma.Register(api, huma.Operation{
		OperationID: "detection-status",
		Method:      http.MethodGet,
		Path:        "/api/detection/status",
	}, func(ctx context.Context, input *struct{}) (
		*models.DetectionStatusResponse, error) {
		status, err := services.GetDetectionStatus(python_url)
		if err != nil {
			log.Println("Unable to reach OCR service:", err)
			return nil, huma.Error500InternalServerError(
				fmt.Sprintf("Failed to reach OCR service: %v", err),
			)
		}
		response := &models.DetectionStatusResponse{}
		response.Body.Status = status.Status
		response.Body.LicensePlate = status.LicensePlate
		response.Body.FolderPath = status.FolderPath
		return response, nil
	})

	//-------------------Stop-Detection-----------------------------
	huma.Register(api, huma.Operation{
		OperationID: "stop-detection",
		Method:      http.MethodPost,
		Path:        "/api/detection/stop",
	}, func(ctx context.Context, input *struct{}) (
		*models.StopDetectionResponse, error) {
		resp, err := http.Post(python_url+"/api/detection/stop", "application/json", nil)
		if err != nil {
			log.Println("Unable to start OCR Engine\n", err)
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
