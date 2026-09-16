package wings

import (
	"context" // Make sure context is imported

	"api.com/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterHandler registers the wings routes
func RegisterHandler(api huma.API, db_pool *pgxpool.Pool) {
	huma.Register(api, huma.Operation{
		OperationID: "get-wings",
		Method:      "GET",
		Path:        "/wings",
	}, func(ctx context.Context, input *struct{}) (*models.WingsResponse, error) {

		// 2. Instantiate and return your named response struct
		resp := &models.WingsResponse{}
		resp.Body.Message = "Wings endpoint working!"

		return resp, nil
	})
}
