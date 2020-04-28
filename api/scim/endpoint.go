package scim

import (
	"encoding/json"
	"fmt"

	"github.com/tcz001/databricks-sdk-go/client"
	"github.com/tcz001/databricks-sdk-go/models"
)

type Endpoint struct {
	Client *client.Client
}

func (c *Endpoint) ListServicePrincipal() (*models.ServicePrincipalsListResponse, error) {
	bytes, err := c.Client.Query("GET", "preview/scim/v2/ServicePrincipals", nil)
	resp := models.ServicePrincipalsListResponse{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) GetServicePrincipal(id string) (*models.ServicePrincipal, error) {
	if id == "" {
		return nil, fmt.Errorf("No Service Principal provided")
	}
	getSPUrl := fmt.Sprintf("preview/scim/v2/ServicePrincipals/%s", id)
	bytes, err := c.Client.Query("GET", getSPUrl, nil)
	resp := models.ServicePrincipal{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) CreateServicePrincipal(request *models.ServicePrincipalCreateRequest) (*models.ServicePrincipal, error) {
	bytes, err := c.Client.Query("POST", "preview/scim/v2/ServicePrincipals", request)
	resp := models.ServicePrincipal{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) UpdateServicePrincipal(updatedServicePrincipal *models.ServicePrincipal) (*models.ServicePrincipal, error) {
	if updatedServicePrincipal.Id == "" {
		return nil, fmt.Errorf("No Service Principal provided")
	}
	updateSPUrl := fmt.Sprintf("preview/scim/v2/ServicePrincipals/%s", updatedServicePrincipal.Id)
	bytes, err := c.Client.Query("PUT", updateSPUrl, updatedServicePrincipal)
	resp := models.ServicePrincipal{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) DeleteServicePrincipal(id string) error {
	if id == "" {
		return fmt.Errorf("No Service Principal provided")
	}
	deleteSPUrl := fmt.Sprintf("preview/scim/v2/ServicePrincipals/%s", id)
	resp, err := c.Client.Query("DELETE", deleteSPUrl, nil)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func (c *Endpoint) ListGroups() (*models.ListGroupRequestScim,error) {

	bytes, err := c.Client.Query("GET", "preview/scim/v2/Groups", nil)
	resp := models.ListGroupRequestScim{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Endpoint) CreateGroup(request *models.ScimGroup) (*models.ScimGroup, error) {
	bytes, err := c.Client.Query("POST", "preview/scim/v2/Groups", request)
	resp := models.ScimGroup{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Endpoint) GetGroup(id string) (*models.ScimGroup,error) {
	if id == "" {
		return nil, fmt.Errorf("No Group id provided")
	}
	getGroupUrl := fmt.Sprintf("preview/scim/v2/Groups/%s", id)
	bytes, err := c.Client.Query("GET", getGroupUrl, nil)
	resp := models.ScimGroup{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Endpoint) DeleteGroup(id string) (error) {
	if id == "" {
		return fmt.Errorf("No Group id provided")
	}
	deleteGroupUrl := fmt.Sprintf("preview/scim/v2/Groups/%s", id)
	resp, err := c.Client.Query("DELETE", deleteGroupUrl, nil)
	fmt.Println(resp)
	if err != nil {
		return err
	}
	return nil
}

func (c *Endpoint) UpdateGroup(id string,group models.ScimGroup) (*models.ScimGroup,error) {
	if id == "" {
		return nil,fmt.Errorf("No Group id provided")
	}
	updateGroupUrl := fmt.Sprintf("preview/scim/v2/Groups/%s", id)
	bytes, err := c.Client.Query("PUT", updateGroupUrl, group)
	resp := models.ScimGroup{}
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}