package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	nedusa "github.com/vincent-vade/medusa-sdk-golang"
	"github.com/vincent-vade/medusa-sdk-golang/request"
	"github.com/vincent-vade/medusa-sdk-golang/response"
	"github.com/vincent-vade/medusa-sdk-golang/utils"
)

type ExistsData struct {
	// Whether email exists or not.
	Exists bool `json:"exists"`
}

type ExistsResponse struct {
	// Success response
	Data *ExistsData

	// Error response
	Error *response.Error

	// Errors in case of multiple errors
	Errors *response.Errors
}

// verify customer emaail address
func Exists(email string, config *nedusa.Config) (*ExistsResponse, error) {
	path := fmt.Sprintf("/store/auth/%v", email)
	resp, err := request.NewRequest().SetMethod(http.MethodGet).SetPath(path).Send(config)
	if err != nil {
		return nil, err
	}
	body, err := utils.ParseResponseBody(resp)
	if err != nil {
		return nil, err
	}

	respBody := new(ExistsResponse)
	switch resp.StatusCode {
	case http.StatusOK:
		respData := new(ExistsData)
		if json.Unmarshal(body, respData); err != nil {
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
