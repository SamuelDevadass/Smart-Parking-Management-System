package handlers

import (
	"context"
	"fmt" // Added for missing case error checks
	"log"
	"net/http"

	"api.com/models"
	"api.com/repositories"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterBillsHandler(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-latest-bill",
		Method:      http.MethodGet,
		Path:        "/api/bills/{license_plate}/latest",
	}, func(ctx context.Context, input *models.GetLicensePlate) (*models.GetLatestBillResponse, error) {
		ans, err := repositories.GetLatestBill(ctx, input)
		if err != nil {
			log.Println("Failed to fetch vehicle details \n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch bill details")
		}
		if ans == nil {
			return nil, huma.Error404NotFound(fmt.Sprintf("No bill details found for license plate '%s'", input.LicensePlate))
		}
		ans.Body.Message = fmt.Sprintf("Bill Details for license plate '%s'", input.LicensePlate)
		return ans, nil
	})
}
