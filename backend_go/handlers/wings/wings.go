package wings

import (
	"context" // Make sure context is imported

	"api.com/models"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterHandler registers the wings routes
func RegisterHandler(api huma.API, db_pool *pgxpool.Pool) {

	//Empty Route
	huma.Register(api, huma.Operation{ /*operation is equivalent of a route*/
		OperationID: "get-wings", /*Unique Id for OpenAPI docs*/
		Method:      "GET",
		Path:        "/api/wings", /*path prefix*/
	}, func(ctx context.Context, input *struct{}) (*models.WingsResponse, error) {
		/*ctx and struct are mandatory input parameters, they can be left empty*/
		/*output response is Wings.response defined in models like schema in FastAPI*/
		// 2. Instantiate and return your named response struct
		resp := &models.WingsResponse{}
		resp.Body.Message = "Wings endpoint working!"

		return resp, nil
	})
}
