package products

import (
	"encoding/json"
	"fmt"
	"net/http"

	medusa "github.com/vincent-vade/medusa-sdk-golang"
	"github.com/vincent-vade/medusa-sdk-golang/request"
	"github.com/vincent-vade/medusa-sdk-golang/response"
	"github.com/vincent-vade/medusa-sdk-golang/schema"
	"github.com/vincent-vade/medusa-sdk-golang/utils"
)

type RetrieveProductData struct {
	Product []*schema.Product `json:"product"`
}

type RetrieveProductResponse struct {
	// Success response
	Data *RetrieveProductData

	// Error response
	Error *response.Error

	// Errors in case of multiple errors
	Errors *response.Errors
}

// Retrieves a Product.
func Retrieve(id string, config *medusa.Config) (*RetrieveProductResponse, error) {
	path := fmt.Sprintf("/store/products/%v", id)
	resp, err := request.NewRequest().SetMethod(http.MethodGet).SetPath(path).Send(config)
	if err != nil {
		return nil, err
	}
	body, err := utils.ParseResponseBody(resp)
	if err != nil {
		return nil, err
	}
	respBody := new(RetrieveProductResponse)
	switch resp.StatusCode {
	case http.StatusOK:
		respData := new(RetrieveProductData)
		if err := json.Unmarshal(body, respData); err != nil {
			return nil, err
		}
		respBody.Data = respData

	case http.StatusUnauthorized:
		respErr := utils.UnauthorizeError()
		respBody.Error = respErr

	case http.StatusBadRequest:
		respErrors, err := utils.ParseErrors(body)
		if err != nil {
			return nil, err
		}
		if len(respErrors.Errors) == 0 {
			respError, err := utils.ParseError(body)
			if err != nil {
				return nil, err
			}
			respBody.Error = respError
		} else {
			respBody.Errors = respErrors
		}

	default:
		respErr, err := utils.ParseError(body)
		if err != nil {
			return nil, err
		}
		respBody.Error = respErr
	}

	return respBody, nil
}
