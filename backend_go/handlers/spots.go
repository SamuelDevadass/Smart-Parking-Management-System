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

func RegisterHandler(api huma.API) {

	huma.Register(api, huma.Operation{
		OperationID: "get-spot-availability",
		Method:      http.MethodGet,
		Path:        "/api/{wing}/spots/availability",
	}, func(ctx context.Context, input *models.WingInput) (*models.GetSpotAvailabilityResponse, error) {
		spots, err := repositories.GetSpotAvailability(ctx, input.Wing)
		if err != nil {
			log.Println("Failed to fetch spots\n", err)
			return nil, huma.Error500InternalServerError("Failed to fetch spots")
		}

		if len(spots) == 0 {
			return nil, huma.Error404NotFound(fmt.Sprintf("No spots found for wing '%s'", input.Wing))
		}
		resp := &models.GetSpotAvailabilityResponse{}
		resp.Body.Message = fmt.Sprintf("Total and Free spots for wing '%s'", input.Wing)
		resp.Body.TotalSpotsTwoWheeler = spots["total_spots_two_wheeler"]
		resp.Body.FreeSpotsTwoWheeler = spots["free_spots_two_wheeler"]
		resp.Body.TotalSpotsFourWheeler = spots["total_spots_four_wheeler"]
		resp.Body.FreeSpotsFourWheeler = spots["free_spots_four_wheeler"]

		return resp, nil
	})
}
