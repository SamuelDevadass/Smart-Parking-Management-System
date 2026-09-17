package wings

import (
	"context"
	"log"
	"net/http"

	"api.com/models"
	"api.com/repositories"
	"github.com/danielgtaylor/huma/v2"
)

// RegisterHandler registers the wings routes
func RegisterHandler(api huma.API) {

	//Empty Route
	huma.Register(api, huma.Operation{ /*operation is equivalent of a route*/
		OperationID: "get-wings", /*Unique Id for OpenAPI docs*/
		Method:      http.MethodGet,
		Path:        "/api/wings", /*path prefix*/
	}, func(ctx context.Context, input *struct{}) (*models.WingsResponse, error) {
		/*ctx and struct are mandatory input parameters, they can be left empty*/
		/*output response is Wings.response defined in models like schema in FastAPI*/

		//Call the function
		wings, err := repositories.ListWings(ctx)
		if err != nil {
			log.Println("Failed to fetch data\n", err)
			return nil, huma.Error500InternalServerError("Failed to look up wings")
		}
		//Instantiate and return your named response struct
		resp := &models.WingsResponse{}
		resp.Body.Message = "Wings endpoint working!"
		resp.Body.Wings = wings

		return resp, nil
	})
}
