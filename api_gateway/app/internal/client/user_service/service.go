package user_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Cadet-Blue/backend-go/api_gateway/internal/apperror"
	"github.com/Cadet-Blue/backend-go/api_gateway/pkg/logging"
	"github.com/Cadet-Blue/backend-go/api_gateway/pkg/rest"
	"github.com/fatih/structs"
	"github.com/go-playground/validator/v10"
)

var _ UserService = &client{}

type client struct {
	base     rest.BaseClient
	Resource string
}

func NewService(baseURL string, resource string, logger logging.Logger) UserService {
	c := client{
		Resource: resource,
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
	return &c
}

type UserService interface {
	Create(ctx context.Context, dto CreateUserDTO) error
}

func (c *client) Create(ctx context.Context, dto CreateUserDTO) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("convert dto to map")
	structs.DefaultTagName = "json"
	data := structs.Map(dto)

	c.base.Logger.Debug("marshal map to bytes")
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal dto")
	}

	validate := validator.New()

	if err := validate.Struct(dto); err != nil {
		errs := []string{}
		validationErrors := err.(validator.ValidationErrors)
		for _, validationError := range validationErrors {
			errs = append(errs, validationError.Error())
		}

		return apperror.ValidationError(errs)
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %w", err)
	}

	if response.IsOk {
		c.base.Logger.Debug("parse location header")
		userURL, err := response.Location()
		if err != nil {
			return fmt.Errorf("failed to get Location header")
		}
		c.base.Logger.Tracef("Location: %s", userURL.String())
		return nil
	}
	return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
}
