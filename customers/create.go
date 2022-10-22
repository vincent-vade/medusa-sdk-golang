package customers

import (
	"encoding/json"
	"net/http"

	medusa "github.com/harshmngalam/medusa-sdk-golang"
	"github.com/harshmngalam/medusa-sdk-golang/request"
	"github.com/harshmngalam/medusa-sdk-golang/response"
	"github.com/harshmngalam/medusa-sdk-golang/schema"
	"github.com/harshmngalam/medusa-sdk-golang/utils"
)

type CreateCustomerData struct {
	Customer *schema.Customer `json:"customer"`
}

type CreateCustomerResponse struct {
	// Success response
	Data *CreateCustomerData

	// Error response
	Error *response.Error

	// Errors in case of multiple errors
	Errors *response.Errors
}

type CreateCustomer struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone,omitempty"`
}

func NewCreateCustomer() *CreateCustomer {
	return new(CreateCustomer)
}

func (c *CreateCustomer) SetFirstName(firstName string) *CreateCustomer {
	c.FirstName = firstName
	return c
}

func (c *CreateCustomer) SetLastName(lastName string) *CreateCustomer {
	c.LastName = lastName
	return c
}

func (c *CreateCustomer) SetEmail(email string) *CreateCustomer {
	c.Email = email
	return c
}

func (c *CreateCustomer) SetPassword(password string) *CreateCustomer {
	c.Password = password
	return c
}

func (c *CreateCustomer) SetPhone(phone string) *CreateCustomer {
	c.Phone = phone
	return c
}

func (c *CreateCustomer) Create(config *medusa.Config) (*CreateCustomerResponse, error) {

	path := "/store/customers"
	resp, err := request.NewRequest().SetMethod(http.MethodPost).SetPath(path).SetData(c).Send(config)
	if err != nil {
		return nil, err
	}
	body, err := utils.ParseResponseBody(resp)
	if err != nil {
		return nil, err
	}

	respBody := new(CreateCustomerResponse)
	switch resp.StatusCode {
	case http.StatusOK:
		respData := new(CreateCustomerData)
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
