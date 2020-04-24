package scim

import (
	"encoding/json"
	"github.com/tcz001/databricks-sdk-go/client"
	"github.com/tcz001/databricks-sdk-go/models"
)

type Endpoint struct {
	Client *client.Client
}

func (c *Endpoint) List() (*models.ServicePrincipalsListResponse, error) {
	bytes, err := c.Client.Query("GET", "preview/scim/v2/ServicePrincipals", nil)
	resp := models.ServicePrincipalsListResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
