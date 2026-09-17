package wings

import (
	"context"
	"fmt" // Added for missing case error checks
	"log"
	"net/http"

	"api.com/models"
	"api.com/repositories"
	"github.com/danielgtaylor/huma/v2"
)

// RegisterHandler registers the wings routes
func RegisterHandler(api huma.API) {

	// Route 1: Get All Wings
	huma.Register(api, huma.Operation{
		OperationID: "get-wings",
		Method:      http.MethodGet,
		Path:        "/api/wings",
	}, func(ctx context.Context, input *struct{}) (*models.GetWingsResponse, error) {
		wings, err := repositories.ListWings(ctx)
		if err != nil {
			log.Println("Failed to fetch data\n", err)
			return nil, huma.Error500InternalServerError("Failed to look up wings")
		}

		resp := &models.GetWingsResponse{}
		resp.Body.Message = "List of all Wings registered"
		resp.Body.Wings = wings

		return resp, nil
	}) // <-- Cleanly closed Route 1

	// Route 2: Get Centre ID For Specific Wing
	huma.Register(api, huma.Operation{
		OperationID: "get-centre-id-for-wing",
		Method:      http.MethodGet,
		Path:        "/api/wings/{wing}/centre",
	}, func(ctx context.Context, input *models.WingInput) (*models.GetCentreForWingsResponse, error) { // FIX 1: Removed invalid '{}' from type parameter

		// Access the bound structural field string value
		centre_id, err := repositories.GetCentreForWings(ctx, input.Wing)
		if err != nil {
			log.Println("Failed to fetch data\n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch centre_id")
		}

		// Validation check: Return a 404 error if the string is empty
		if centre_id == "" {
			return nil, huma.Error404NotFound(fmt.Sprintf("No centre found for wing '%s'", input.Wing))
		}

		resp := &models.GetCentreForWingsResponse{}
		resp.Body.Message = "Centre_id for current Wing"
		resp.Body.Centre = centre_id

		return resp, nil
	}) // <-- FIX 3: Placed closing brackets directly after the handler block
}
