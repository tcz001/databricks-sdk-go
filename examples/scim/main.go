package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tcz001/databricks-sdk-go/api/scim"
	"github.com/tcz001/databricks-sdk-go/client"
	"github.com/tcz001/databricks-sdk-go/models"
	"io/ioutil"
)

func main() {
	flag.Parse() // required to suppress warnings from glog

	secrets := loadSecrets()

	cl, err := client.NewClient(client.Options{
		Domain:                              &secrets.Domain,
		Token:                               &secrets.Token,
		XDatabricksAzureWorkspaceResourceId: &secrets.DBWorkspaceResourceId,
		XDatabricksAzureSPManagementToken:   &secrets.SPMgmtToken,
	})
	if err != nil {
		panic(err)
	}

	endpoint := scim.Endpoint{
		Client: cl,
	}

	//printServicePrincipals(listServicePrincipals(endpoint))
	//printUserGroups(listUserGroups(endpoint))
	//printCreatedGroups(createUserGroups(endpoint))
	//printGetGroup(getUserGroup(endpoint,"2812227734411558"))
	deleteUserGroup(endpoint,"4593788545285908")
}

func printGetGroup(group *models.ScimGroup) {
	fmt.Println("GroupId:",group.Id)
	fmt.Println("GroupDisplayName:",group.DisplayName)
	for index, element := range group.Members {
		fmt.Println(index, "displayName: =>", element.Display)
		fmt.Println(index, "id: =>", element.Value)
	}
}

func getServicePrincipal(endpoint scim.Endpoint, id string) *models.ServicePrincipal {
	fmt.Println("Getting Service Principals for id:", id)
	resp, err := endpoint.GetServicePrincipal(id)
	fmt.Println("Response:\n", resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func deleteServicePrincipal(endpoint scim.Endpoint, id string) {
	fmt.Println("Deleting Service Principals for id:", id)
	err := endpoint.DeleteServicePrincipal(id)
	if err != nil {
		panic(err)
	}
}

func createServicePrincipal(endpoint scim.Endpoint, appId string, displayName string) *models.ServicePrincipal {
	fmt.Println("Attaching Service Principals for id:", appId)
	req := models.ServicePrincipalCreateRequest{}
	req.ApplicationId = appId
	req.DisplayName = displayName
	resp, err := endpoint.CreateServicePrincipal(&req)
	if err != nil {
		panic(err)
	}
	return resp
}

func updateServicePrincipal(endpoint scim.Endpoint, id string, updatedSP models.ServicePrincipal) *models.ServicePrincipal {
	fmt.Println("Updating Service Principals for SP id:", id)
	resp, err := endpoint.UpdateServicePrincipal(&updatedSP)
	if err != nil {
		panic(err)
	}
	return resp
}

func printServicePrincipals(principals []models.ServicePrincipal) {
	fmt.Println(principals)
}

func printUserGroups(groups []models.ScimGroup) {
	for index, element := range groups {
		fmt.Println(index, "displayName: =>", element.DisplayName)
		fmt.Println(index, "id: =>", element.Id)
		members := element.Members
		for _,element2 := range members {
			fmt.Println("value id: =>", element2.Value)
			fmt.Println("displayName with Member: =>", element2.Display)
		}

	}
	fmt.Println(groups)
}

func printCreatedGroups(displayName string, members []models.ScimMember) {
	fmt.Println("group displayName :",displayName)
	for index, element := range members {
		fmt.Println(index, "Member id: =>", element.Value)
		fmt.Println(index, "Member displayNAme: =>", element.Display)
		fmt.Println(index, "Member ref: =>", element.Ref)
	}
}

func listServicePrincipals(endpoint scim.Endpoint) []models.ServicePrincipal {
	fmt.Println("Listing Service Principals")
	resp, err := endpoint.ListServicePrincipal()
	if err != nil {
		panic(err)
	}
	return resp.Resources
}

type secrets struct {
	Domain                string `json:"domain"`
	Token                 string `json:"token"`
	ClusterName           string `json:"cluster_name"`
	SPToken               string `json:"sp_token"`
	SPMgmtToken           string `json:"sp_mgmt_token"`
	DBWorkspaceResourceId string `json:"workspace_resource_id"`
}

func loadSecrets() *secrets {
	content, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		panic(err)
	}

	var sc secrets
	err = json.Unmarshal(content, &sc)
	if err != nil {
		panic(err)
	}

	return &sc
}

func listUserGroups(endpoint scim.Endpoint) []models.ScimGroup {
	fmt.Println("Listing User Groups")
	resp, err := endpoint.ListUserGroups()
	if err != nil {
		panic(err)
	}
	return resp.Resources
}

func createUserGroups(endpoint scim.Endpoint) (string,[]models.ScimMember) {
	fmt.Println("Creating User Group")
	member := models.ScimMember{
		Display: "",
		Value:   "5648206897659689",
		Ref:     "",
	}
	group := models.ScimGroup{
		Entitlements: nil,
		DisplayName:  "blah6",
		Members:      []models.ScimMember{member},
		Groups:       nil,
		Id:           "",
	}
	resp, err := endpoint.CreateUserGroup(&group)
	if err != nil {
		panic(err)
	}
	return resp.DisplayName,resp.Members
}

func getUserGroup(endpoint scim.Endpoint,id string) *models.ScimGroup {
	fmt.Println("Getting User Group with id:",id)
	resp, err := endpoint.GetUserGroup(id)
	if err != nil {
		panic(err)
	}
	return resp
}

func deleteUserGroup(endpoint scim.Endpoint,id string) {
	fmt.Println("Deleting User Group with id:",id)
	err := endpoint.DeleteUserGroup(id)
	if err != nil {
		panic(err)
	}
}